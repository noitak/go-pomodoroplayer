package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

const playcmd = "/usr/bin/afplay" // mac only!!

type Pomodoro struct {
	WorkMin   time.Duration
	RestMin   time.Duration
	WorkSongs []string
	RestSongs []string
}

func (p *Pomodoro) Start() {
	for pomodoroCount, workSong, restSong := 0, 0, 0; ; pomodoroCount++ {
		fmt.Printf("[%d] Pomodoro Start [%v]\n", pomodoroCount+1, p.WorkMin)

		workSong = play(p.WorkMin, p.WorkSongs, workSong)

		fmt.Printf("[%d] Pomodoro End\n", pomodoroCount+1)

		restSong = play(p.RestMin, p.RestSongs, restSong)
	}
}

func timer(t time.Duration, ch chan string) {
	time.Sleep(t)
	ch <- "timeout"
}

func play(t time.Duration, songs []string, songPos int) int {
	ch := make(chan string)

	go timer(t, ch)

	var cmd *exec.Cmd
	for {
		if songPos >= len(songs) {
			songPos = 0
		}
		cmd = exec.Command(playcmd, songs[songPos])

		go func() {
			cmd.Start()
			cmd.Wait()
			ch <- "playend"
		}()

		s := <-ch
		if s == "timeout" {
			cmd.Process.Kill()
			break
		} else if s == "playend" {
			songPos++
		}
	}
	return songPos
}

func main() {
	var configfile = flag.String("c", "", "music_config.json")
	flag.Parse()
	if len(*configfile) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	raw, err := ioutil.ReadFile(*configfile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var pomodoro Pomodoro

	json.Unmarshal(raw, &pomodoro)

	if pomodoro.WorkMin == 0 || pomodoro.RestMin == 0 {
		fmt.Println("Set WorkMin and RestMin > 1 min")
		os.Exit(1)
	}
	if len(pomodoro.WorkSongs) == 0 || len(pomodoro.RestSongs) == 0 {
		fmt.Printf("No songs in %s\n", *configfile)
		os.Exit(1)
	}
	pomodoro.WorkMin *= time.Minute
	pomodoro.RestMin *= time.Minute

	pomodoro.Start()
}
