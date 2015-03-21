package main

import (
	"fmt"
)

type MSSStrategy struct {
	Strategy
}

type Strategy struct {
	Table  Table
	Action string
	Bet    string
}

var latePositionsRaisePushHands = []string{"AA", "KK", "QQ", "JJ", "AK"}

var raisePushHands = map[string][]string{
	"EP": []string{"AA", "KK", "QQ"},
	"MP": []string{"AA", "KK", "QQ", "AK"},
	"CO": latePositionsRaisePushHands,
	"BU": latePositionsRaisePushHands,
	"SB": latePositionsRaisePushHands,
	"BB": latePositionsRaisePushHands,
}

var latePositionsRaiseFoldHands = []string{
	"TT", "99", "88", "77", "AQ", "AJ", "AT", "A9s", "KQ",
}

var raiseFoldHands = map[string][]string{
	"EP": []string{"JJ", "TT", "AK", "AQ", "AJs"},
	"MP": []string{"JJ", "TT", "99", "88", "AQ", "AJ", "ATs"},
	"CO": latePositionsRaiseFoldHands,
	"BU": latePositionsRaiseFoldHands,
	"SB": latePositionsRaiseFoldHands,
	"BB": latePositionsRaiseFoldHands,
}

func (strategy Strategy) Run() {
	err := strategy.Check()

	if err != nil {
		fmt.Println("bad input")
		return
	}

	if strategy.Table.Board == "" {
		strategy.Preflop()
	} else {
		fmt.Println("flop, turn or river strategy")
	}
}

func (strategy Strategy) Preflop() {
	position := strategy.Table.Hero.Position
	hand := strategy.Table.Hero.Hand.FoldedNotification()

	if hand == "" {
		fmt.Println("no hand provided")
		return
	}

	for _, element := range raiseFoldHands[position] {
		if element == hand {
			fmt.Println("raise fold")
			return
		}
	}

	for _, element := range raisePushHands[position] {
		if element == hand {
			fmt.Println("raise push")
			return
		}
	}

	if position == "BB" {
		fmt.Println("check fold")
	} else {
		fmt.Println("fold")
	}
}

func (strategy Strategy) Check() error {
	return nil
}
