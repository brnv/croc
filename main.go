package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/docopt/docopt-go"
	"github.com/op/go-logging"
	"github.com/thearkit/runcmd"
)

var (
	log       = logging.MustGetLogger("croc")
	cmdRunner *runcmd.Local
)

var (
	importCmd           = "/bin/import png:%s"
	compareCmd          = "/bin/compare -quiet -metric RMSE %s %s NULL:"
	convertCmd          = "/bin/convert -crop %dx%d+%d+%d %s %s"
	reCompareErrorLevel = regexp.MustCompile("\\((.*)\\)$")
)

const usage = `
	Usage:
	croc [<filepath>] [--call=CALL] [--blinds=BLINDS] [--ante=ANTE]
`

type Table struct {
	Hero
	Blinds string
	Ante   string
	Pot    string
	Board  string
}

type Image struct {
	Path string
}

type Hero struct {
	Chips string
	Call  string
	Hand  string
}

type ImageSnippet struct {
	Width   int
	Height  int
	OffsetX int
	OffsetY int
}

func (image Image) Crop(snippet ImageSnippet) string {
	croppedPath, _ := getTmpFilename()

	command, _ := cmdRunner.Command(fmt.Sprintf(
		convertCmd,
		snippet.Width, snippet.Height, snippet.OffsetX, snippet.OffsetY,
		image.Path, croppedPath),
	)
	_, _ = command.Run()

	return croppedPath
}

func getImageSnippets(
	width int,
	height int,
	offsetY int,
	offsets []int,
) []ImageSnippet {
	imageSnippets := make([]ImageSnippet, len(offsets))

	for index, offset := range offsets {
		imageSnippets[index] = ImageSnippet{
			Width:   width,
			Height:  height,
			OffsetX: offset,
			OffsetY: offsetY,
		}
	}

	return imageSnippets
}

func main() {
	var err error

	cmdRunner, err = runcmd.NewLocalRunner()
	if err != nil {
		log.Fatal(err)
	}

	image := Image{}

	args, _ := docopt.Parse(usage, nil, true, "croc", false)
	if args["<filepath>"] != nil {
		image.Path = args["<filepath>"].(string)
	} else {
		image.Path, err = makeScreenshot()
		if err != nil {
			log.Fatal(err)
		}
	}

	table := Table{
		Hero: Hero{},
	}

	wg := &sync.WaitGroup{}
	wg.Add(4)

	go func() {
		table.Hero.Hand = image.HandRecognize()
		wg.Done()
	}()

	go func() {
		table.Pot = image.PotRecognize()
		wg.Done()
	}()

	go func() {
		table.Hero.Chips = image.HeroChipsRecognize()
		wg.Done()
	}()

	go func() {
		table.Board = image.BoardRecognize()
		wg.Done()
	}()

	if args["--call"] == nil {
		wg.Add(1)
		go func() {
			table.Hero.Call = image.CallRecognize()
			wg.Done()
		}()
	} else {
		table.Hero.Call = args["--call"].(string)
	}

	wg.Wait()

	log.Notice("Input: %v", image.Path)
	log.Notice("Pot: %v", table.Pot)
	log.Notice("Board: %v", table.Board)
	log.Notice("Hero hand: %v", table.Hero.Hand)
	log.Notice("Hero chips: %v", table.Hero.Chips)
	log.Notice("Hero call: %v", table.Hero.Call)

	if args["--blinds"] != nil {
		table.Blinds = args["--blinds"].(string)
		log.Notice("Blinds: %v", table.Blinds)
	}

	if args["--ante"] != nil {
		table.Ante = args["--ante"].(string)
		log.Notice("Ante: %v", table.Ante)
	}
}

func recognize(
	input string,
	samplesFilepathPattern string,
	compareThreshold float64,
) (string, error) {
	samples, _ := filepath.Glob(samplesFilepathPattern)

	for _, sample := range samples {
		command, _ := cmdRunner.Command(fmt.Sprintf(compareCmd, sample, input))

		_, err := command.Run()
		if err == nil {
			return path.Base(sample), nil
		}

		errorLevel, _ := strconv.ParseFloat(
			reCompareErrorLevel.FindStringSubmatch(err.Error())[1],
			32,
		)

		if errorLevel < compareThreshold {
			return path.Base(sample), nil
		}
	}

	return "", errors.New(
		fmt.Sprintf("Recognition failed! Input file: %s", input),
	)
}

func makeScreenshot() (string, error) {
	screenshot, err := getTmpFilename()
	if err != nil {
		return "", err
	}

	command, err := cmdRunner.Command(fmt.Sprintf(importCmd, screenshot))

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
