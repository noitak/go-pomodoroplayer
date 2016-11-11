package main

import (
	"encoding/json"
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
	fmt.Printf("Pomodoro Start [%v]\n", p.WorkMin)

	play(p.WorkMin, p.WorkSongs)
	play(p.RestMin, p.RestSongs)

	fmt.Println("Pomodoro End")
}

func timer(t time.Duration, ch chan string) {
	time.Sleep(t)
	ch <- "timeout"
}

func play(t time.Duration, songs []string) {
	ch := make(chan string)

	go timer(t, ch)

	var cmd *exec.Cmd
	for _, song := range songs {
		cmd = exec.Command(playcmd, song)

		fmt.Printf("start: %s\n", song)
		go func() {
			cmd.Start()
			cmd.Wait()
			ch <- "playend"
		}()

		s := <-ch
		if s == "timeout" {
			cmd.Process.Kill()
			break
		}
	}
}

func main() {
	var musiclist string

	if len(os.Args) > 1 {
		musiclist = os.Args[1]
	} else {
		fmt.Println("Usage: pmodoroplayer musiclist.json")
		os.Exit(1)
	}
	raw, err := ioutil.ReadFile(musiclist)
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
		fmt.Printf("No songs in %s\n", musiclist)
		os.Exit(1)
	}
	pomodoro.WorkMin *= time.Minute
	pomodoro.RestMin *= time.Minute

	pomodoro.Start()
}
