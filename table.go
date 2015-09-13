package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/seletskiy/tplutil"
)

const (
	tableTpl = `
	[hand {{.Hero.Hand}}, 
	position {{.Hero.Position}}, 
	board {{.Board}}, 
	chips {{.Hero.Chips}}, 
	pot {{.Pot}}]`

	playersCount = 9
)

type Table struct {
	Window Window
	Image  Image
	Hero
	Board
	Opponents      []Opponent
	Pot            int
	ButtonPosition int
	Errors         []string
	BigBlindSize   int
}

type Image struct {
	Path string
}

type Hero struct {
	Chips    int
	Hand     Hand
	Position int
}

type ImageSnippet struct {
	Width   int
	Height  int
	OffsetX int
	OffsetY int
}

func (table *Table) Recognize() {
	wg := &sync.WaitGroup{}

	wg.Add(5)

	go func() {
		table.HandRecognize()
		wg.Done()
	}()

	go func() {
		table.ButtonRecognize()
		table.HeroPositionRecognize()
		wg.Done()
	}()

	go func() {
		table.OpponentsRecognize()
		table.RaisersRecognize()
		wg.Done()
	}()

	go func() {
		table.PotRecognize()
		wg.Done()
	}()

	go func() {
		table.BoardRecognize()
		wg.Done()
	}()

	wg.Wait()
}

func (table Table) SitOut() {
	table.Window.Click(12, 410)
	table.Fold()
}

func (table Table) Fold() {
	table.Window.Click(400, 505)
}

func (table Table) Check() {
	if table.CheckButtonIsVisible() {
		table.Window.Click(540, 505)
	}
}

func (table Table) Call() {
	table.Window.Click(540, 505)
}

func (table Table) Raise() {
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

func (table Table) ContBet() {
	table.Window.Click(600, 440)
	table.Window.Click(680, 505)
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

func (table Table) GetButtonRelativePosition(offset int) int {
	if table.ButtonPosition < offset {
		return offset - table.ButtonPosition
	} else {
		return playersCount + offset - table.ButtonPosition
	}
}

func (table *Table) HeroPositionRecognize() {
	table.Hero.Position = table.GetButtonRelativePosition(5)
}

func (table Table) GetFirstRaiserPosition() int {
	lowestPosition := playersCount + 1

	for _, opponent := range table.Opponents {
		opponentPosition := table.GetButtonRelativePosition(opponent.Index)
		if opponent.Raiser && opponentPosition < lowestPosition {
			lowestPosition = opponentPosition
		}
	}

	return lowestPosition
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

func (table Table) CheckButtonIsVisible() bool {
	fastFoldButton := ImageSnippet{
		120, 23, 522, 494,
	}

	_, err := recognize(
		table.Image.Crop(fastFoldButton),
		"/tmp/croc/button_check",
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

const (
	sitOutTopChipsAmount    = 180
	sitOutBottomChipsAmount = 10
)

func (table Table) PerformAutomatedActions(decision string) {
	switch decision {

	case "CHECK":
		table.Check()

	case "FOLD":
		table.HeroChipsRecognize()

		if table.Chips != 0 &&
			(table.Hero.Chips >= sitOutTopChipsAmount ||
				table.Hero.Chips <= sitOutBottomChipsAmount) {
			table.SitOut()
		} else {
			table.Fold()
		}

	case "RAISE/ALL-IN":
		table.RaiseAllIn()

	case "FLOP RAISE/ALL-IN":
		table.FlopRaiseAllIn()

	case "TURN RAISE/ALL-IN":
		table.TurnRaiseAllIn()

	case "RIVER CHECK/CALL":
		table.RiverCheckCall()

	case "RAISE/FOLD":
		table.RaiseFold()

	case "RAISE/MANUAL":
		table.RaisePlayerMove()

	case "FLOP CHECK/FOLD":
		table.FlopCheckFold()

	case "CHECK/FOLD":
		table.CheckFold()

	case "FLOP C-BET/FOLD":
		table.ContBetFold("flop")

	case "FLOP C-BET/MANUAL":
		table.ContBetManualMove("flop")
	}
}

func (table Table) RiverCheckCall() {
	performTwoActions(
		table.Check, table.Call,
		fmt.Sprintf("/tmp/croc-river-check-call-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) TurnRaiseAllIn() {
	performTwoActions(
		table.Raise, table.AllIn,
		fmt.Sprintf("/tmp/croc-turn-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) FlopRaiseAllIn() {
	performTwoActions(
		table.Raise, table.AllIn,
		fmt.Sprintf("/tmp/croc-flop-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) RaiseAllIn() {
	performTwoActions(
		table.Raise, table.AllIn,
		fmt.Sprintf("/tmp/croc-allin-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) ContBetManualMove(street string) {
	performTwoActions(
		table.ContBet, table.WaitFold,
		fmt.Sprintf(
			"/tmp/croc-%s-c-bet-fold-%s-%s",
			street,
			table.Hero.Hand,
			table.Window.Id),
	)
}

func (table Table) RaisePlayerMove() {
	performTwoActions(
		table.Raise, table.WaitFold,
		fmt.Sprintf("/tmp/croc-wait-player-move-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) RaiseFold() {
	performTwoActions(
		table.Raise, table.Fold,
		fmt.Sprintf("/tmp/croc-fold-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) WaitFold() {
	flag := fmt.Sprintf(
		"/tmp/croc-wait-fold-%s-%s",
		table.Hero.Hand,
		table.Window.Id,
	)

	file, err := os.Stat(flag)

	if os.IsNotExist(err) {
		createFlagFile(flag)
	} else if file.ModTime().Unix() < time.Now().Unix()-20 {
		table.Fold()
	}
}

func (table Table) ContBetFold(street string) {
	performTwoActions(
		table.ContBet, table.Fold,
		fmt.Sprintf(
			"/tmp/croc-%s-c-bet-fold-%s-%s",
			street,
			table.Hero.Hand,
			table.Window.Id),
	)
}

func (table Table) FlopCheckFold() {
	performTwoActions(
		table.Check, table.Fold,
		fmt.Sprintf("/tmp/croc-flop-check-fold-%s-%s", table.Hero.Hand, table.Window.Id),
	)
}

func (table Table) CheckFold() {
	performTwoActions(
		table.Check, table.Fold,
		fmt.Sprintf("/tmp/croc-check-fold-%s-%s", table.Hero.Hand, table.Window.Id),
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

func (table Table) PotIsRaised() bool {
	limpTotalSize := 0

	for _, opponent := range table.Opponents {
		if opponent.Limper {
			limpTotalSize += opponent.ChipsInPot
		}
	}

	if positions[table.Hero.Position] == "BB" {
		limpTotalSize += table.BigBlindSize
	}

	if limpTotalSize != table.Pot &&
		limpTotalSize+1 != table.Pot {
		return true
	}

	return false
}

func (table Table) IsStealerPosition(position int) bool {
	if position == 8 ||
		position == 9 ||
		position == 1 {
		return true
	}

	return false
}

func (table Table) IsPreflop() bool {
	if len(table.Board.Cards) == 0 {
		return true
	}

	return false
}

func (table Table) IsFlop() bool {
	if len(table.Board.Cards) == 3 {
		return true
	}

	return false
}

func (table Table) IsTurn() bool {
	if len(table.Board.Cards) == 4 {
		return true
	}

	return false
}

func (table Table) IsRiver() bool {
	if len(table.Board.Cards) == 5 {
		return true
	}

	return false
}
