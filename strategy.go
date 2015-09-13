package main

import (
	"errors"
	"fmt"
)

type Strategy struct {
	Table    Table
	Messages []string
	Decision string
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

func (strategy *Strategy) Run() error {
	err := strategy.CheckInput()

	if err != nil {
		return err
	}

	if strategy.Table.IsPreflop() {
		strategy.Decision = strategy.PreflopDecision()
	}

	if strategy.Table.IsFlop() {
		strategy.Decision = strategy.FlopDecision()
	}

	if strategy.Table.IsTurn() {
		strategy.Decision = strategy.TurnDecision()
	}

	if strategy.Table.IsRiver() {
		strategy.Decision = strategy.RiverDecision()
	}

	return nil
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

	//hand := strategy.Table.Hero.Hand.ShortNotation()

	//strategy.PrintReminders()

	if strategy.IsGoodHand() {
		return "TURN RAISE/ALL-IN"
	}

	return "CHECK/FOLD"
}

func (strategy *Strategy) RiverDecision() string {
	strategy.Messages = append(strategy.Messages, "river")

	//hand := strategy.Table.Hero.Hand.ShortNotation()

	//@TODO: automate this logic

	//fmt.Println("monster: BET/RAISE or BET/CALL;")
	//fmt.Println("anything else: CHECK/FOLD;")

	if strategy.IsGoodHand() {
		return "RIVER CHECK/CALL"
	}

	return "CHECK/FOLD"
}

func (strategy Strategy) CheckInput() error {
	if strategy.Table.Hero.Hand.ShortNotation() == "" {
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
