package main

import (
	"bytes"
	"html/template"
	"os"
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

const tableTpl = `
	Hero hand: {{.Hero.Hand}}{{"\n"}}
	Hero position: {{.Hero.Position}}{{"\n"}}
	Hero chips: {{.Hero.Chips}}{{"\n"}}
	Pot size: {{.Pot}}{{"\n"}}
	Board: {{.Board}}{{"\n"}}
`

const usage = `
	Usage:
	croc [<filepath>]`

func main() {
	runtime.GOMAXPROCS(4)

	var err error

	logging.SetLevel(logging.NOTICE, "")

	cmdRunner, err = runcmd.NewLocalRunner()
	if err != nil {
		log.Fatal(err)
	}

	image := Image{}

	window := Window{}

	args, _ := docopt.Parse(usage, nil, true, "croc", false)

	if args["<filepath>"] != nil {
		image.Path = args["<filepath>"].(string)
	} else {
		window, err = getWindow()
		if err != nil {
			log.Fatal(err)
		}

		image.Path, err = getWindowScreenshot(window.Id)
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
		table.Limpers = image.LimpersRecognize()
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

	strategy := MSSStrategy{
		Strategy: Strategy{
			Table: table,
		},
	}

	strategy.Run()
}

type Table struct {
	Hero
	Board
	Pot     int
	Limpers []Limper
	Button  string
}

type Image struct {
	Path string
}

type Hero struct {
	Chips    string
	Hand     Hand
	Position int
}

type ImageSnippet struct {
	Width   int
	Height  int
	OffsetX int
	OffsetY int
}

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

func (table Table) GetHeroPosition() int {
	buttonNum, _ := strconv.Atoi(table.Button)
	return len(positions) + 1 - buttonNum
}
