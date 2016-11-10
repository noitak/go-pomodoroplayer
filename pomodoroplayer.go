package main

import (
	"os/exec"
	"fmt"
	"time"
)

const playcmd = "/usr/bin/afplay" // mac only!!

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
			<- playend
		}
	}()

	time.Sleep(t)
	cmd.Process.Kill()
	
}

func main() {
	worksongs := []string{
		"./Japanese_School_Bell02-12.mp3",
		"./Japanese_School_Bell02-12.mp3",
	}
	restsongs := []string{
		"./gup01.mp4",
	}

	working(11 * time.Second, worksongs)
	working(5 * time.Second, restsongs)
}
