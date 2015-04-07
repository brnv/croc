package main

import (
	"errors"
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

func getWindow() (Window, error) {
	window := Window{}

	command, _ := cmdRunner.Command(windowInfoCmd)
	output, err := command.Run()
	if err != nil {
		return window, err
	}

	matches := reWindowId.FindStringSubmatch(output[4])
	if len(matches) != 0 {
		window.Id = matches[1]
	} else {
		return window, errors.New("No window id found")
	}

	matches = reWindowX.FindStringSubmatch(output[6])
	if len(matches) != 0 {
		window.X, _ = strconv.Atoi(matches[1])
	}

	matches = reWindowY.FindStringSubmatch(output[7])
	if len(matches) != 0 {
		window.Y, _ = strconv.Atoi(matches[1])
	}

	return window, nil
}

func getWindowCoordinates(id string) (int, int) {
	command, _ := cmdRunner.Command(fmt.Sprintf(windowInfoByIdCmd, id))
	output, err := command.Run()
	if err != nil {
		return 0, 0
	}

	x, y := 0, 0

	matches := reWindowX.FindStringSubmatch(output[2])
	if len(matches) != 0 {
		x, _ = strconv.Atoi(matches[1])
	}

	matches = reWindowY.FindStringSubmatch(output[3])
	if len(matches) != 0 {
		y, _ = strconv.Atoi(matches[1])
	}

	return x, y
}

func getWindowScreenshot(windowId string) (string, error) {
	screenshot, err := getTmpFilename()
	if err != nil {
		return "", err
	}

	command, err := cmdRunner.Command(fmt.Sprintf(
		importCmd, windowId, screenshot,
	))

	if err != nil {
		return "", err
	}

	_, err = command.Run()
	if err != nil {
		return "", err
	}

	return screenshot, nil
}

func getTmpFilename() (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), "croc")
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}
