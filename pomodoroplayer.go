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
	WorkSongs []string
	WorkMin   time.Duration
	RestSongs []string
	RestMin   time.Duration
}

func NewPomodoro(
	worksongs []string,
	workmin time.Duration,
	restsongs []string,
	restmin time.Duration) *Pomodoro {
	return &Pomodoro{worksongs, workmin, restsongs, restmin}
}

func (p *Pomodoro) Start() {
	fmt.Printf("Pomodoro Start [%v]\n", p.WorkMin)
	working(p.WorkMin, p.WorkSongs)
	working(p.RestMin, p.RestSongs)
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

		fmt.Printf("start: %s\n", song)
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

	pomodoro.WorkMin *= time.Minute
	pomodoro.RestMin *= time.Minute

	if pomodoro.WorkMin == 0 || pomodoro.RestMin == 0 {
		fmt.Println("set WorkMin and RestMin > 1 min")
		os.Exit(0)
	}
	pomodoro.Start()
}
