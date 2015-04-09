package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

var (
	reWindowId = regexp.MustCompile("Window id: (0x[a-z0-9]+)\\s+.*")
	reWindowX  = regexp.MustCompile("Absolute upper-left X:\\s+(\\d+)")
	reWindowY  = regexp.MustCompile("Absolute upper-left Y:\\s+(\\d+)")

	windowInfoCmd     = "/bin/xwininfo"
	windowInfoByIdCmd = "/bin/xwininfo -id %s"
	importCmd         = "/bin/import -window %s png:%s"
)

type Window struct {
	Id string
	X  int
	Y  int
}

func (window *Window) InitId() {
	command, _ := cmdRunner.Command(windowInfoCmd)

	output, _ := command.Run()

	matches := reWindowId.FindStringSubmatch(output[4])
	if len(matches) != 0 {
		window.Id = matches[1]
	}
}

func (window *Window) InitCoordinates() {
	command, _ := cmdRunner.Command(fmt.Sprintf(
		windowInfoByIdCmd, window.Id),
	)

	output, err := command.Run()
	if err != nil {
		return
	}

	matches := reWindowX.FindStringSubmatch(output[2])
	if len(matches) != 0 {
		window.X, _ = strconv.Atoi(matches[1])
	}

	matches = reWindowY.FindStringSubmatch(output[3])
	if len(matches) != 0 {
		window.Y, _ = strconv.Atoi(matches[1])
	}
}

func (window Window) Screenshot() string {
	screenshot, err := getTmpFilename()
	if err != nil {
		return ""
	}

	command, err := cmdRunner.Command(fmt.Sprintf(
		importCmd, window.Id, screenshot,
	))

	if err != nil {
		return ""
	}

	_, err = command.Run()
	if err != nil {
		return ""
	}

	return screenshot
}

func (window Window) Click(offsetX int, offsetY int) {
	window.InitCoordinates()

	command, _ := cmdRunner.Command(
		fmt.Sprintf(
			"/bin/xdotool mousemove %d %d click 1",
			window.X+offsetX,
			window.Y+offsetY,
		),
	)

	command.Run()
}

func getTmpFilename() (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "croc")
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
