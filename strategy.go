package main

import (
	"errors"
	"fmt"
)

type Strategy struct {
	Table    Table
	Messages []string
}

var positions = map[int]string{
	1: "SB",
	2: "BB",
	3: "EP",
	4: "EP",
	5: "MP",
	6: "MP",
	7: "MP",
	8: "CO",
	9: "BU",
}

const (
	laterPosition = "LATER"
)

var strategyPositions = map[string]string{
	"EP": "EP",
	"MP": "MP",
	"CO": laterPosition,
	"BU": laterPosition,
	"SB": laterPosition,
	"BB": laterPosition,
}

func (strategy *Strategy) Run() string {
	err := strategy.CheckInput()

	if err != nil {
		return err.Error()
	}

	boardCardsCount := len(strategy.Table.Board.Cards)

	if boardCardsCount == 0 {
		return strategy.PreflopDecision()
	}

	if boardCardsCount == 3 {
		return strategy.FlopDecision()
	} else if boardCardsCount == 4 {
		return strategy.TurnDecision()
	}

	return strategy.RiverDecision()
}

func (strategy *Strategy) PreflopDecision() string {
	strategy.Messages = append(strategy.Messages, "preflop")

	if strategy.PreflopRaiseSituation() {
		return strategy.PreflopRaiseDecision()
	}

	if strategy.PreflopStealSituation() {
		return strategy.PreflopStealDecision()
	}

	if strategy.PreflopThreeBetSituation() {
		return strategy.PreflopThreeBetDecision()
	}

	if strategy.PreflopRestealSituation() {
		return strategy.PreflopRestealDecision()
	}

	return "MANUAL"
}

func (strategy Strategy) PreflopRaiseSituation() bool {
	if !strategy.Table.PotIsRaised() &&
		!strategy.PreflopStealSituation() {
		return true
	}

	return false
}

const defaultPotSize = 3

func (strategy Strategy) PreflopStealSituation() bool {
	heroPosition := strategy.Table.Hero.Position

	if strategyPositions[positions[heroPosition]] == laterPosition &&
		strategy.Table.Pot == defaultPotSize {
		return true
	}

	return false
}

const avgStealSizePot = 12

var stealerPosition = map[int]bool{
	8: true,
	9: true,
	1: true,
	2: true,
}

func (strategy Strategy) PreflopRestealSituation() bool {
	heroPosition := strategy.Table.Hero.Position

	if positions[heroPosition] == "BB" &&
		stealerPosition[strategy.Table.GetFirstRaiserPosition()] &&
		strategy.Table.PotIsRaised() &&
		strategy.Table.Pot <= avgStealSizePot {
		return true
	}

	return false
}

func (strategy Strategy) PreflopThreeBetSituation() bool {
	if !strategy.PreflopRaiseSituation() &&
		!strategy.PreflopStealSituation() &&
		!strategy.PreflopRestealSituation() {
		return true
	}

	return false
}

func (strategy *Strategy) PreflopRaiseDecision() string {
	strategy.Messages = append(strategy.Messages, "raise")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range pushHands {
		if hand == card {
			return "RAISE/ALL-IN"
		}
	}

	for _, card := range raiseWaitPlayerHands {
		if hand == card {
			return "RAISE/MANUAL"
		}
	}

	for _, card := range callHands {
		if hand == card {
			return "RAISE/MANUAL"
		}
	}

	for _, element := range raiseFoldHands[position] {
		if element == hand {
			return "RAISE/FOLD"
		}
	}

	if position == "BB" {
		return "CHECK"
	}

	return "FOLD"
}

func (strategy *Strategy) PreflopStealDecision() string {
	strategy.Messages = append(strategy.Messages, "steal")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range pushHands {
		if hand == card {
			return "RAISE/ALL-IN"
		}
	}

	for _, card := range stealWaitPlayerHands {
		if hand == card {
			return "RAISE/MANUAL"
		}
	}

	for _, card := range stealFoldHands[position] {
		if hand == card {
			return "RAISE/FOLD"
		}
	}

	if position == "BB" {
		return "CHECK"
	}

	return "FOLD"
}

func (strategy *Strategy) PreflopRestealDecision() string {
	strategy.Messages = append(strategy.Messages, "resteal")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range pushHands {
		if hand == card {
			return "RAISE/ALL-IN"
		}
	}

	for _, card := range stealWaitPlayerHands {
		if hand == card {
			return "RAISE/MANUAL"
		}
	}

	for _, card := range restealFoldHands {
		if hand == card {
			return "RAISE/FOLD"
		}
	}

	return "FOLD"
}

const potSaneLimitForThreeBet = 18

func (strategy *Strategy) PreflopThreeBetDecision() string {
	strategy.Messages = append(strategy.Messages, "3-bet")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	raiserPosition := positions[strategy.Table.GetFirstRaiserPosition()]

	strategy.Messages = append(strategy.Messages, "raiser in "+raiserPosition)

	for _, card := range pushHands {
		if hand == card {
			return "RAISE/ALL-IN"
		}
	}

	for _, card := range callHands {
		if hand == card {
			return "MANUAL CALL"
		}
	}

	for _, card := range threeBetHands[strategyPositions[raiserPosition]] {
		if hand == card {
			if strategy.Table.Pot > potSaneLimitForThreeBet {
				return "MANUAL"
			}
			return "RAISE/MANUAL"
		}
	}

	return "FOLD"
}

func (strategy *Strategy) FlopDecision() string {
	strategy.Messages = append(strategy.Messages, "flop")

	hero := strategy.Table.Hero
	board := strategy.Table.Board
	completedCombination := hero.Hand.GetCompletedCombination(board)
	hand := strategy.Table.Hero.Hand.ShortNotation()

	if completedCombination.OverPair {
		strategy.Messages = append(strategy.Messages, "overpair")

		return "MANUAL"
	}

	if completedCombination.Three || completedCombination.Triplet {
		strategy.Messages = append(strategy.Messages, "three")

		return "MANUAL"
	}

	if completedCombination.TwoPairs {
		strategy.Messages = append(strategy.Messages, "two pairs")

		return "MANUAL"
	}

	if completedCombination.StrongTopPair {
		strategy.Messages = append(strategy.Messages, "strong top pair")

		return "MANUAL"
	}

	if completedCombination.TopPair {
		strategy.Messages = append(strategy.Messages, "top pair")

		if strategy.Table.Pot <= 10 {
			return "FLOP CHECK/FOLD"
		} else if strategy.Table.Pot <= 35 {
			return "FLOP C-BET/FOLD"
		}
	}

	for _, card := range contBetPairs {
		if hand == card && strategy.Table.Pot <= 35 {
			strategy.Messages = append(strategy.Messages, "pair")

			return "FLOP C-BET/FOLD"
		}
	}

	emptyCombination := hero.Hand.GetEmptyCombination(board)

	if emptyCombination.String() != "" {
		if emptyCombination.OverCards {
			strategy.Messages = append(strategy.Messages, "overcards")

			if strategy.Table.Pot <= 16 {
				return "FLOP C-BET/FOLD"
			}
		}
	}

	strategy.PrintReminders()

	return "MANUAL"
}

func (strategy *Strategy) TurnDecision() string {
	strategy.Messages = append(strategy.Messages, "turn")

	hero := strategy.Table.Hero
	board := strategy.Table.Board
	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.OverPair {
		strategy.Messages = append(strategy.Messages, "overpair")

		return "MANUAL"
	}

	if completedCombination.Three || completedCombination.Triplet {
		strategy.Messages = append(strategy.Messages, "three")

		return "MANUAL"
	}

	if completedCombination.TwoPairs {
		strategy.Messages = append(strategy.Messages, "two pairs")

		return "MANUAL"
	}

	if completedCombination.StrongTopPair {
		strategy.Messages = append(strategy.Messages, "strong top pair")

		return "MANUAL"
	}

	if completedCombination.TopPair {
		strategy.Messages = append(strategy.Messages, "top pair")
		if strategy.Table.Pot <= 10 {
			return "CHECK/FOLD"
		}

		return "MANUAL"
	}

	emptyCombination := hero.Hand.GetEmptyCombination(board)

	if emptyCombination.String() != "" {
		if emptyCombination.OverCards {
			strategy.Messages = append(strategy.Messages, "overcards")
			return "CHECK/FOLD"
		}
	}

	strategy.PrintReminders()

	return "MANUAL"
}

func (strategy *Strategy) RiverDecision() string {
	strategy.Messages = append(strategy.Messages, "river")

	hero := strategy.Table.Hero
	board := strategy.Table.Board
	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.OverPair ||
		completedCombination.Three ||
		completedCombination.Triplet ||
		completedCombination.TopPair ||
		completedCombination.TwoPairs {
		return "MANUAL"
	}

	//@TODO: automate this logic

	fmt.Println("monster: BET/RAISE or BET/CALL;")
	fmt.Println("anything else: CHECK/FOLD;")

	return "MANUAL"
}

func (strategy Strategy) CheckInput() error {
	hand := strategy.Table.Hero.Hand.ShortNotation()

	if hand == "" {
		return errors.New("no hand provided")
	}

	return nil
}

func (strategy Strategy) PrintReminders() {
	//@TODO: automate this logic

	fmt.Println("monster draw: BET/ALL-IN or RERAISE;")

	fmt.Println(
		fmt.Sprintf(
			"draws: if win_size / call_size / odds > 1: CALL;",
		),
	)
}
