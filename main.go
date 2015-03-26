package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/docopt/docopt-go"
	"github.com/op/go-logging"
	"github.com/seletskiy/tplutil"
	"github.com/thearkit/runcmd"
)

var (
	log       = logging.MustGetLogger("croc")
	cmdRunner *runcmd.Local
)

var (
	importCmd  = "/bin/import png:%s"
	compareCmd = "/bin/compare " +
		"-dissimilarity-threshold 1 -quiet -metric RMSE %s %s NULL:"
	convertCmd          = "/bin/convert -crop %dx%d+%d+%d %s %s"
	reCompareErrorLevel = regexp.MustCompile("\\((.*)\\).*$")
)

const usage = `
	Usage:
	croc [<filepath>] [--call=CALL] [--blinds=BLINDS] [--ante=ANTE]
`

type Table struct {
	Hero
	Blinds    string
	Ante      string
	Pot       int
	Board     Board
	Opponents []Opponent
	Button    string
}

type Image struct {
	Path string
}

type Hero struct {
	Chips    string
	Call     string
	Hand     Hand
	Position int
}

type ImageSnippet struct {
	Width   int
	Height  int
	OffsetX int
	OffsetY int
}

const tableTpl = `
	Hero hand: {{.Hero.Hand}}{{"\n"}}
	Hero position: {{.Hero.Position}}{{"\n"}}
	Hero chips: {{.Hero.Chips}}{{"\n"}}
	Pot size: {{.Pot}}{{"\n"}}
	Board: {{.Board}}{{"\n"}}
`

func (table Table) String() string {
	myTpl := template.Must(
		template.New("table").Parse(tplutil.Strip(
			tableTpl,
		)),
	)

	buf := bytes.NewBuffer([]byte{})

	myTpl.Execute(buf, table)

	return buf.String()
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
	runtime.GOMAXPROCS(4)

	var err error

	logging.SetLevel(logging.NOTICE, "")

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

	if _, err := os.Stat(image.Path); os.IsNotExist(err) {
		log.Error("no such file or directory: " + image.Path)
		os.Exit(1)
	}

	table := Table{
		Hero: Hero{},
	}

	wg := &sync.WaitGroup{}
	wg.Add(6)

	go func() {
		table.Hero.Hand = image.HandRecognize()
		wg.Done()
	}()

	go func() {
		table.Pot, _ = strconv.Atoi(strings.TrimLeft(image.PotRecognize(), "0"))
		wg.Done()
	}()

	go func() {
		table.Opponents = image.OpponentsRecognize()
		wg.Done()
	}()

	go func() {
		table.Hero.Chips = strings.TrimLeft(image.HeroChipsRecognize(), "0")
		wg.Done()
	}()

	go func() {
		table.Board = image.BoardRecognize()
		wg.Done()
	}()

	go func() {
		table.Button = image.ButtonRecognize()
		table.Hero.Position = table.GetHeroPosition()
		wg.Done()
	}()

	wg.Wait()

	if args["--call"] != nil {
		table.Hero.Call = args["--call"].(string)
	}

	if args["--blinds"] != nil {
		table.Blinds = args["--blinds"].(string)
		fmt.Printf("Blinds: %v\n", table.Blinds)
	}

	if args["--ante"] != nil {
		table.Ante = args["--ante"].(string)
		fmt.Printf("Ante: %v\n", table.Ante)
	}

	strategy := MSSStrategy{
		Strategy: Strategy{
			Table: table,
		},
	}

	strategy.Run()
}

func (table Table) GetHeroPosition() int {
	buttonNum, _ := strconv.Atoi(table.Button)
	return len(positions) + 1 - buttonNum
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

		compareErrorLevel := reCompareErrorLevel.FindStringSubmatch(err.Error())

		if len(compareErrorLevel) > 0 {
			errorLevel, _ := strconv.ParseFloat(
				compareErrorLevel[1],
				32,
			)

			if errorLevel < compareThreshold {
				return path.Base(sample), nil
			}

		}
	}

	return "", errors.New(fmt.Sprintf("%s failed!", input))

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
