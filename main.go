package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
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

const tableTpl = `
	Hero hand: {{.Hero.Hand}}{{"\n"}}
	Hero position: {{.Hero.Position}}{{"\n"}}
	Hero chips: {{.Hero.Chips}}{{"\n"}}
	Pot size: {{.Pot}}{{"\n"}}
	Board: {{.Board}}{{"\n"}}
`

const usage = `
	Usage:
	croc [<filepath>] [--wid=<window_id>] [-v]`

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

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
	} else if args["--wid"] != nil {
		window.Id = args["--wid"].(string)
		window.X, window.Y = getWindowCoordinates(window.Id)
	} else {
		window, err = getWindow()
		if err != nil {
			log.Fatal(err)
		}
	}

	if image.Path == "" {
		image.Path, err = getWindowScreenshot(window.Id)
		if err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(image.Path); os.IsNotExist(err) {
		log.Error("no such file or directory: " + image.Path)
		os.Exit(1)
	}

	if !image.CheckIfHeroTurn() {
		fmt.Print(".")
		return
	}

	fmt.Println("")

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

	decision := strategy.Run()

	//@TODO: remember decision for current window id
	//@TODO: refactor all click logic

	if window.Id != "" {
		mouseX, mouseY := rememberMousePosition()

		switch decision {
		case "CHECK":
			clickOnCheckButton(window)
		case "FOLD":
			clickOnFoldButton(window)

		case "RAISE/FOLD":
			clickOnRaiseButton(window)
		case "RAISE/ALL-IN":
			clickOnRaiseButton(window)

		case "STEAL/FOLD":
			clickOnStealButton(window)
		case "STEAL/ALL-IN":
			clickOnStealButton(window)

		case "3-BET/FOLD if raiser >= EP":
			clickOnThreeBetButton(window)
		case "3-BET/ALL-IN if raiser >= EP":
			clickOnThreeBetButton(window)

		case "3-BET/FOLD if raiser >= MP":
			clickOnThreeBetButton(window)
		case "3-BET/ALL-IN if raiser >= MP":
			clickOnThreeBetButton(window)

		case "3-BET/FOLD if raiser >= LATER":
			clickOnThreeBetButton(window)
		case "3-BET/ALL-IN if raiser >= LATER":
			clickOnThreeBetButton(window)

		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= EP":
			clickOnThreeBetButton(window)
		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= MP":
			clickOnThreeBetButton(window)
		case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= LATER":
			clickOnThreeBetButton(window)

		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= EP":
			clickOnThreeBetButton(window)
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= MP":
			clickOnThreeBetButton(window)
		case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= LATER":
			clickOnThreeBetButton(window)
		}

		restoreMousePosition(mouseX, mouseY)
	}

	if args["-v"].(bool) != false {
		fmt.Print(table)
	}

	fmt.Println(decision)
}

func clickOnCheckButton(window Window) {
	click(window.X+560, window.Y+520)
}

func clickOnFoldButton(window Window) {
	click(window.X+440, window.Y+520)
}

func clickOnRaiseButton(window Window) {
	click(window.X+620, window.Y+440)
	click(window.X+720, window.Y+520)
}

func clickOnStealButton(window Window) {
	click(window.X+720, window.Y+520)
}

func clickOnThreeBetButton(window Window) {
	click(window.X+560, window.Y+440)
	click(window.X+720, window.Y+520)
}

var (
	reMouseX = regexp.MustCompile("x:(\\d+)\\s")
	reMouseY = regexp.MustCompile("y:(\\d+)\\s")
)

func rememberMousePosition() (string, string) {
	command, _ := cmdRunner.Command(
		fmt.Sprintf("/bin/xdotool getmouselocation"),
	)
	output, _ := command.Run()
	mouseX := reMouseX.FindStringSubmatch(output[0])
	mouseY := reMouseY.FindStringSubmatch(output[0])

	return mouseX[1], mouseY[1]
}

func restoreMousePosition(x string, y string) {
	command, _ := cmdRunner.Command(
		fmt.Sprintf("/bin/xdotool mousemove %s %s", x, y),
	)
	command.Run()
}

func click(x int, y int) {
	command, _ := cmdRunner.Command(
		fmt.Sprintf("/bin/xdotool mousemove %d %d click 1", x, y),
	)
	command.Run()
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

func (image Image) CheckIfHeroTurn() bool {
	maxButton := ImageSnippet{
		61, 23, 719, 432,
	}

	_, err := recognize(
		image.Crop(maxButton),
		"/tmp/croc/button_max",
		0.05,
	)

	if err != nil {
		return false
	}

	return true
}
