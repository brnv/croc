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

func (strategy *Strategy) IsGoodHand() bool {
	hero := strategy.Table.Hero
	board := strategy.Table.Board
	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.OverPair ||
		completedCombination.Three ||
		completedCombination.Triplet ||
		completedCombination.TwoPairs ||
		completedCombination.StrongTopPair {
		strategy.Messages = append(strategy.Messages, "good hand")
		return true
	}

	return false
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
