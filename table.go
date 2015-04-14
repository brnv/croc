package main

import (
	"bytes"
	"os"
	"strconv"
	"text/template"

	"github.com/seletskiy/tplutil"
)

const (
	tableTpl = `
	[hand {{.Hero.Hand}}, 
	position {{.Hero.Position}}, 
	board {{.Board}}, 
	pot {{.Pot}}]`
)

type Table struct {
	Window Window
	Image  Image
	Hero
	Board
	Limpers []Limper
	Pot     int
	Button  string
	Errors  []string
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

func (table Table) ClickFold() {
	table.Window.Click(440, 520)
}

func (table Table) ClickCheck() {
	table.Window.Click(560, 520)
}

func (table Table) ClickRaise() {
	table.Window.Click(620, 440)
	table.Window.Click(720, 520)
}

func (table Table) ClickSteal() {
	table.Window.Click(720, 520)
}

func (table Table) ClickThreeBet() {
	table.Window.Click(560, 440)
	table.Window.Click(720, 520)
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

func (table *Table) HeroPositionRecognize() {
	buttonNum, _ := strconv.Atoi(table.Button)

	table.Hero.Position = len(positions) + 1 - buttonNum
}

func (table Table) HeroMoveIsPending() bool {
	maxButton := ImageSnippet{
		61, 23, 719, 432,
	}

	_, err := recognize(
		table.Image.Crop(maxButton),
		"/tmp/croc/button_max",
		0.05,
	)

	if err != nil {
		return false
	}

	return true
}

func (table Table) FastFoldButtonIsVisible() bool {
	fastFoldButton := ImageSnippet{
		128, 25, 382, 490,
	}

	_, err := recognize(
		table.Image.Crop(fastFoldButton),
		"/tmp/croc/button_fast_fold_minimized",
		0.05,
	)

	if err != nil {
		return false
	}

	return true
}

func (table *Table) Check() bool {
	if _, err := os.Stat(table.Image.Path); os.IsNotExist(err) {
		table.Errors = append(
			table.Errors, "no such file or directory: "+table.Image.Path,
		)
	}

	if len(table.Errors) != 0 {
		return false
	}

	return true
}
