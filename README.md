# go-pomodoro player

Music Player for [Pomodoro Technique](http://pomodorotechnique.com/).

## Requirements
* afplay command (Mac only)

## How to use

	$ go run pomodoro.go music.json

music.json is configuration 

	{
		"workmin": 25,
		"restmin": 5,
		"worksongs": [
			"/path/to/working/music01.mp3",
			"/path/to/working/music02.m4a"
		],
		"restsongs": [
			"/path/to/rest/music01.mp3",
			"/path/to/rest/music02.m4a",
		]
	}
