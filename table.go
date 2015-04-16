package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/seletskiy/tplutil"
)

const (
	tableTpl = `
	[hand {{.Hero.Hand}}, 
	position {{.Hero.Position}}, 
	board {{.Board}}, 
	pot {{.Pot}}]`

	potSaneLimitForThreeBet = 18
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

func (table Table) Fold() {
	table.Window.Click(400, 505)
}

func (table Table) Check() {
	table.Window.Click(540, 505)
}

func (table Table) Raise() {
	table.Window.Click(680, 440)
	table.Window.Click(680, 505)
}

func (table Table) Steal() {
	table.Window.Click(680, 440)
	table.Window.Click(680, 505)
}

func (table Table) ThreeBet() {
	table.Window.Click(680, 440)
	table.Window.Click(680, 505)
}

func (table Table) AllIn() {
	table.Window.Click(760, 440)
	table.Window.Click(680, 505)
}

func (table Table) Bet() {
	table.Window.Click(630, 440)
	table.Window.Click(680, 505)
}

func (table Table) FastFoldToAnyBet() {
	table.Window.Click(12, 396)
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
	table.Hero.Position = len(positions) + 1 - table.ButtonPosition
}

func (table Table) HeroMoveInProgress() bool {
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

func (table *Table) Validate() bool {
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

func (table Table) FastFoldToAnyBetIsChecked() bool {
	fastFoldCheckbox := ImageSnippet{
		65, 18, 5, 386,
	}

	_, err := recognize(
		table.Image.Crop(fastFoldCheckbox),
		"/tmp/croc/fast_fold_checkbox",
		0.05,
	)

	if err != nil {
		return true
	}

	return false
}

func (table Table) PerformAutomatedActions(decision string) {
	switch decision {
	case "CHECK":
		table.Check()
	case "FOLD":
		table.Fold()

	case "RAISE/FOLD":
		table.RaiseFold()
	case "STEAL/FOLD":
		table.StealFold()

	case "3-BET/FOLD if raiser >= EP":
		table.ThreeBetFold()
	case "3-BET/FOLD if raiser >= MP":
		table.ThreeBetFold()
	case "3-BET/FOLD if raiser >= LATER":
		table.ThreeBetFold()
	case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= EP":
		table.ThreeBetFold()
	case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= MP":
		table.ThreeBetFold()
	case "RESTEAL/FOLD\n3-BET/FOLD if raiser >= LATER":
		table.ThreeBetFold()

	case "RESTEAL/ALL-IN\n3-BET/FOLD if raiser >= EP":
		table.ThreeBetFold()
	case "RESTEAL/ALL-IN\n3-BET/FOLD if raiser >= MP":
		table.ThreeBetFold()
	case "RESTEAL/ALL-IN\n3-BET/FOLD if raiser >= LATER":
		table.ThreeBetFold()

	case "RAISE/ALL-IN":
		table.RaiseAllIn()
	case "STEAL/ALL-IN":
		table.StealAllIn()
	case "3-BET/ALL-IN":
		table.ThreeBetAllIn()
	case "3-BET/ALL-IN if raiser >= EP":
		table.ThreeBetAllIn()
	case "3-BET/ALL-IN if raiser >= MP":
		table.ThreeBetAllIn()
	case "3-BET/ALL-IN if raiser >= LATER":
		table.ThreeBetAllIn()

	case "RESTEAL/ALL-IN\n3-BET/ALL-IN":
		table.ThreeBetAllIn()
	case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= EP":
		table.ThreeBetAllIn()
	case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= MP":
		table.ThreeBetAllIn()
	case "RESTEAL/ALL-IN\n3-BET/ALL-IN if raiser >= LATER":
		table.ThreeBetAllIn()

	case "FLOP BET/ALL-IN":
		table.FlopBetAllIn()

	case "TURN BET/ALL-IN":
		table.TurnBetAllIn()
	}
}

func (table Table) RaiseFold() {
	performTwoActions(
		table.Raise, table.Fold,
		fmt.Sprintf("/tmp/croc-fold-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) RaiseAllIn() {
	performTwoActions(
		table.Raise, table.AllIn,
		fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) StealFold() {
	performTwoActions(
		table.Steal, table.Fold,
		fmt.Sprintf("/tmp/croc-fold-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) StealAllIn() {
	performTwoActions(
		table.Steal, table.AllIn,
		fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) ThreeBetFold() {
	flag := fmt.Sprintf("/tmp/croc-fold-%s-%s", table.Hero.Hand, table.Window.Id)

	if !flagFileIsOk(flag) && table.Pot <= potSaneLimitForThreeBet {
		createFlagFile(flag)
		table.ThreeBet()
	} else {
		table.Fold()
	}
}

func (table Table) ThreeBetAllIn() {
	performTwoActions(
		table.ThreeBet, table.AllIn,
		fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) FlopBetAllIn() {
	performTwoActions(
		table.Bet, table.AllIn,
		fmt.Sprintf("/tmp/croc-flop-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) TurnBetAllIn() {
	performTwoActions(
		table.Bet, table.AllIn,
		fmt.Sprintf("/tmp/croc-turn-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func performTwoActions(firstAction func(), secondAction func(), flag string) {
	if !flagFileIsOk(flag) {
		createFlagFile(flag)
		firstAction()
	} else {
		secondAction()
	}
}

func createFlagFile(name string) {
	os.Create(name)
}

func flagFileIsOk(flag string) bool {
	file, err := os.Stat(flag)

	if os.IsNotExist(err) {
		return false
	}

	if file.ModTime().Unix() < time.Now().Unix()-60 {
		return false
	}

	return true
}
