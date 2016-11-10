package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
	"io/ioutil"
)

const playcmd = "/usr/bin/afplay" // mac only!!

type Pomodoro struct {
	WorkSongs []string
	WorkTime  time.Duration
	RestSongs []string
	RestTime  time.Duration
}

func NewPomodoro(
	worksongs []string,
	worktime time.Duration,
	restsongs []string,
	resttime time.Duration) *Pomodoro {
	return &Pomodoro{worksongs, worktime, restsongs, resttime}
}

func (p *Pomodoro) Start() {
	working(p.WorkTime, p.WorkSongs)
	working(p.RestTime, p.RestSongs)
	fmt.Println("Pomodoro End")
}

func play(cmd *exec.Cmd, ch chan string) {
	cmd.Start()
	cmd.Wait()
	ch <- "playend"
}

func timer(t time.Duration, ch chan string) {
	time.Sleep(t)
	ch <- "timeout"
}

func working(t time.Duration, worksongs []string) {
	playend := make(chan string)

	go timer(t, playend)

	var cmd *exec.Cmd
		for _, song := range worksongs {
			cmd = exec.Command(playcmd, song)

			fmt.Printf("Start: %s\n", song)
			go play(cmd, playend)

			s := <-playend
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

	pomodoro.WorkTime = 4*time.Minute
	pomodoro.RestTime = 2*time.Second

	pomodoro.Start()
}
