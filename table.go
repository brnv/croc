package main

import (
	"bytes"
	"os"
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
	Limpers        []Limper
	Pot            int
	ButtonPosition int
	Errors         []string
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
	table.Window.Click(400, 505)
}

func (table Table) ClickCheck() {
	table.Window.Click(540, 505)
}

func (table Table) ClickRaise() {
	table.Window.Click(680, 440)
	table.Window.Click(680, 505)
}

func (table Table) ClickSteal() {
	table.Window.Click(680, 440)
	table.Window.Click(680, 505)
}

func (table Table) ClickThreeBet() {
	table.Window.Click(680, 440)
	table.Window.Click(680, 505)
}

func (table Table) ClickFastFoldToAnyBet() {
	table.Window.Click(12, 396)
}

//@TODO:all in 760 440

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
	table.Hero.Position = len(positions) + 1 - table.ButtonPosition
}

func (table Table) HeroMoveIsPending() bool {
	maxButton := ImageSnippet{
		15, 23, 765, 493,
	}

	_, err := recognize(
		table.Image.Crop(maxButton),
		"/tmp/croc/raise_button_top_right_corner",
		0.05,
	)

	if err != nil {
		return false
	}

	return true
}

func (table Table) FoldButtonIsVisible() bool {
	fastFoldButton := ImageSnippet{
		15, 23, 382, 490,
	}

	_, err := recognize(
		table.Image.Crop(fastFoldButton),
		"/tmp/croc/fold_button_top_left_corner",
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

func (table Table) FastFoldToAnyBetIsUnchecked() bool {
	fastFoldCheckbox := ImageSnippet{
		65, 18, 5, 386,
	}

	_, err := recognize(
		table.Image.Crop(fastFoldCheckbox),
		"/tmp/croc/fast_fold_checkbox",
		0.05,
	)

	if err != nil {
		return false
	}

	return true
}
