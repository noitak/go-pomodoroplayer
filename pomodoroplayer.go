package main

import (
	"fmt"
	"os/exec"
	"time"
)

const playcmd = "/usr/bin/afplay" // mac only!!

type Pomodoro struct {
	worksongs []string
	restsongs []string
	worktime  time.Duration
	resttime  time.Duration
}

func NewPomodoro(
	worksongs []string,
	restsongs []string,
	worktime time.Duration,
	resttime time.Duration) *Pomodoro {
	return &Pomodoro{worksongs, restsongs, worktime, resttime}
}

func (p *Pomodoro) Start() {
	working(p.worktime, p.worksongs)
	working(p.resttime, p.restsongs)
}

func play(cmd *exec.Cmd, playend chan struct{}) {
	cmd.Run()
	playend <- struct{}{}
}

func working(t time.Duration, worksongs []string) {
	var cmd *exec.Cmd
	go func() {
		playend := make(chan struct{})
		for _, song := range worksongs {
			cmd = exec.Command(playcmd, song)

			fmt.Printf("Start: %s\n", song)
			go play(cmd, playend)
			<-playend
		}
	}()

	time.Sleep(t)
	cmd.Process.Kill()

}

func main() {
	pomodoro := NewPomodoro(
		[]string{
			"./Japanese_School_Bell02-12.mp3",
			"./Japanese_School_Bell02-12.mp3",
		},
		[]string{
			"./gup01.mp4",
		},
		11*time.Second,
		5*time.Second)
	pomodoro.Start()
}
