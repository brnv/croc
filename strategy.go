package main

import (
	"errors"
	"fmt"
)

type MSSStrategy struct {
	Strategy
}

type Strategy struct {
	Table Table
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

var noLimpPotSize = 3

var laterPosition = "Later"
var strategyPositions = map[string]string{
	"EP": "EP",
	"MP": "MP",
	"CO": laterPosition,
	"BU": laterPosition,
	"SB": laterPosition,
	"BB": laterPosition,
}

// hero position's raise hands
var raisePushHandsLatePosition = []string{
	"AA", "KK", "QQ", "JJ",
	"AK", "AKs",
}
var raisePushHands = map[string][]string{
	"EP": []string{
		"AA", "KK", "QQ",
	},
	"MP": []string{
		"AA", "KK", "QQ",
		"AK", "AKs",
	},
	"CO": raisePushHandsLatePosition,
	"BU": raisePushHandsLatePosition,
	"SB": raisePushHandsLatePosition,
	"BB": raisePushHandsLatePosition,
}
var raiseFoldHandsLatePosition = []string{
	"TT", "99", "88", "77",
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9s",
	"KQ", "KQs",
}
var raiseFoldHands = map[string][]string{
	"EP": []string{
		"JJ", "TT",
		"AK", "AKs", "AQ", "AQs", "AJs",
	},
	"MP": []string{
		"JJ", "TT", "99", "88",
		"AQ", "AQs", "AJ", "AJs", "ATs",
	},
	"CO": raiseFoldHandsLatePosition,
	"BU": raiseFoldHandsLatePosition,
	"SB": raiseFoldHandsLatePosition,
	"BB": raiseFoldHandsLatePosition,
}

// raiser position's 3-bet hands
var allInHands = []string{
	"AA", "KK",
}
var threeBetAllInHandsLatePosition = []string{
	"QQ", "JJ",
	"AK", "AKs",
}
var threeBetAllInHands = map[string][]string{
	"EP": []string{"QQ"},
	"MP": []string{"QQ", "AK", "AKs"},
	"CO": threeBetAllInHandsLatePosition,
	"BU": threeBetAllInHandsLatePosition,
	"SB": threeBetAllInHandsLatePosition,
	"BB": threeBetAllInHandsLatePosition,
}
var threeBetFoldHandsLatePosition = []string{
	"TT", "99", "88",
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9s",
}
var threeBetFoldHands = map[string][]string{
	"MP": []string{
		"JJ", "TT", "99",
		"AQ", "AQs", "AJ", "AJs", "ATs",
	},
	"CO": threeBetFoldHandsLatePosition,
	"BU": threeBetFoldHandsLatePosition,
	"SB": threeBetFoldHandsLatePosition,
	"BB": threeBetFoldHandsLatePosition,
}

func (strategy Strategy) Run() {
	err := strategy.CheckInput()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Hero is %s with %s\n",
		positions[strategy.Table.Hero.Position],
		strategy.Table.Hero.Hand,
	)

	if strategy.Table.Board == "" {
		strategy.Preflop()
	} else {
		fmt.Println("flop, turn or river strategy")
	}
}

func (strategy Strategy) Preflop() {
	if !strategy.PotIsRaised() {
		heroPosition := strategy.Table.Hero.Position

		if strategyPositions[positions[heroPosition]] == laterPosition &&
			strategy.Table.Pot == noLimpPotSize {
			fmt.Println("implement steal strategy")
		} else {
			strategy.PreflopRaiseStrategy()
		}
	} else {
		strategy.PreflopThreeBetStrategy()
	}
}

func (strategy Strategy) PotIsRaised() bool {
	limpTotalSize := 0

	for _, limper := range strategy.Table.Opponents {
		limpTotalSize += limper.LimpSize
	}

	if limpTotalSize != strategy.Table.Pot {
		return true
	}

	return false
}

func (strategy Strategy) PreflopThreeBetStrategy() {
	fmt.Println("preflop 3-bet strategy decision is")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range allInHands {
		if hand == card {
			fmt.Println("ALL-IN on many raises, 3-BET on one raise")
			return
		}
	}

	for position, cards := range threeBetAllInHands {
		for _, card := range cards {
			if hand == card {
				fmt.Printf(
					"3-BET and ALL-IN after 4-BET if opponent in %s position",
					strategyPositions[position],
				)
				return
			}

		}
	}

	for position, cards := range threeBetFoldHands {
		for _, card := range cards {
			if hand == card {
				fmt.Printf(
					"3-BET and FOLD after 4-BET if opponent in %s position",
					strategyPositions[position],
				)
				return
			}

		}
	}

	fmt.Println("FOLD")
}

func (strategy Strategy) PreflopRaiseStrategy() {
	fmt.Println("preflop raise strategy decision is")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, element := range raisePushHands[position] {
		if element == hand {
			fmt.Println("RAISE and ALL-IN after 3-bet")
			return
		}
	}

	for _, element := range raiseFoldHands[position] {
		if element == hand {
			fmt.Println("RAISE and FOLD after 3-bet")
			return
		}
	}

	if position == "SB" {
		fmt.Println("FOLD or LIMP")
		return
	}

	if position == "BB" {
		fmt.Println("CHECK")
		return
	}
}

func (strategy Strategy) CheckInput() error {
	hand := strategy.Table.Hero.Hand.ShortNotation()

	if hand == "" {
		return errors.New("No hand provided")
	}

	return nil
}
