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

var noLimpPotSize = 7

var laterPosition = "LATER"
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
var threeBetFoldMPHands = map[string][]string{
	"MP": []string{
		"JJ", "TT", "99",
		"AQ", "AQs", "AJ", "AJs", "ATs",
	},
}
var threeBetFoldHandsLatePosition = []string{
	"TT", "99", "88",
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9s",
}
var threeBetFoldLATERHands = map[string][]string{
	"CO": threeBetFoldHandsLatePosition,
	"BU": threeBetFoldHandsLatePosition,
	"SB": threeBetFoldHandsLatePosition,
	"BB": threeBetFoldHandsLatePosition,
}

// steal hands
var stealAllInHands = []string{
	"AA", "KK", "QQ", "JJ", "TT",
	"AK", "AKs",
}
var stealFoldHandsBUandSB = []string{
	"99", "88", "77", "66", "55", "44", "33", "22",
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9", "A9s",
	"A8s", "A7s", "A6s", "A5s", "A4s", "A3s", "A2s",
	"KQ", "KQs", "KJ", "KJs", "KT", "KTs",
	"QJ", "QJs", "QT", "QTs",
	"JT", "JTs",
}
var stealFoldHands = map[string][]string{
	"CO": []string{
		"99", "88", "77", "66", "55",
		"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9s",
		"KQ", "KQs", "KJs",
		"QJs",
	},
	"BU": stealFoldHandsBUandSB,
	"SB": stealFoldHandsBUandSB,
}
var reStealFoldHands = map[string][]string{
	"BB": []string{
		"99", "88",
		"AQ", "AQs", "AJ", "AJs",
	},
}

func (strategy Strategy) Run() {
	err := strategy.CheckInput()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(strategy.Table.Board.Cards) == 0 {
		strategy.Preflop()
	} else {
		boardCardsNum := len(strategy.Table.Board.Cards)

		if boardCardsNum == 3 {
			strategy.Flop()
		} else if boardCardsNum == 4 {
			strategy.Turn()
		} else if boardCardsNum == 5 {
			strategy.River()
		}
	}
}

func (strategy Strategy) CheckInput() error {
	hand := strategy.Table.Hero.Hand.ShortNotation()

	if hand == "" {
		return errors.New("no hand provided")
	}

	return nil
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

func (strategy Strategy) Preflop() {
	heroPosition := strategy.Table.Hero.Position

	if !strategy.PotIsRaised() {
		if strategyPositions[positions[heroPosition]] == laterPosition &&
			strategy.Table.Pot == noLimpPotSize {
			strategy.PreflopStealStrategy()
		} else {
			strategy.PreflopRaiseStrategy()
		}
	} else {
		if strategyPositions[positions[heroPosition]] == laterPosition {
			strategy.PreflopReStealStrategy()
		}
		strategy.PreflopThreeBetStrategy()
	}
}

func (strategy Strategy) PreflopStealStrategy() {
	fmt.Print("steal: ")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range stealAllInHands {
		if hand == card {
			fmt.Println("STEAL/ALL-IN")
			return
		}
	}

	for _, card := range stealFoldHands[position] {
		if hand == card {
			fmt.Println("STEAL/FOLD")
			return
		}
	}

	fmt.Println("FOLD")
}

func (strategy Strategy) PreflopRaiseStrategy() {
	fmt.Print("raise: ")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, element := range raisePushHands[position] {
		if element == hand {
			fmt.Println("RAISE/ALL-IN")
			return
		}
	}

	for _, element := range raiseFoldHands[position] {
		if element == hand {
			fmt.Println("RAISE/FOLD")
			return
		}
	}

	if position == "SB" {
		fmt.Println("LIMP or FOLD")
		return
	}

	if position == "BB" {
		fmt.Println("CHECK")
		return
	}

	fmt.Println("FOLD")
}

func (strategy Strategy) PreflopReStealStrategy() {
	fmt.Print("resteal: ")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range stealAllInHands {
		if hand == card {
			fmt.Println("RESTEAL/ALL-IN")
			return
		}
	}

	for _, card := range reStealFoldHands[position] {
		if hand == card {
			fmt.Println("RESTEAL/FOLD")
			return
		}
	}

	fmt.Println("FOLD")
}

func (strategy Strategy) PreflopThreeBetStrategy() {
	fmt.Print("3-bet: ")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range allInHands {
		if hand == card {
			fmt.Println("ALL-IN on 2+ raises or 3-BET on 1")
			return
		}
	}

	for position, cards := range threeBetAllInHands {
		for _, card := range cards {
			if hand == card {
				fmt.Printf(
					"3-BET/ALL-IN if raiser >= %s\n",
					strategyPositions[position],
				)
				return
			}
		}
	}

	for position, cards := range threeBetFoldMPHands {
		for _, card := range cards {
			if hand == card {
				fmt.Printf(
					"3-BET/FOLD if raiser >= %s\n",
					strategyPositions[position],
				)
				return
			}
		}
	}

	for position, cards := range threeBetFoldLATERHands {
		for _, card := range cards {
			if hand == card {
				fmt.Printf(
					"3-BET/FOLD if raiser >= %s\n",
					strategyPositions[position],
				)
				return
			}
		}
	}

	fmt.Println("FOLD")
}

func (strategy Strategy) Flop() {
	hero := strategy.Table.Hero
	board := strategy.Table.Board

	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.String() != "" {
		if completedCombination.OverPair ||
			completedCombination.Three ||
			completedCombination.Triplet ||
			completedCombination.TwoPairs {
			fmt.Println("BET/ALL-IN or RERAISE")
			return
		}

		if completedCombination.TopPair {
			fmt.Println("C-BET/FOLD or FOLD")
			fmt.Println("freeplay: CHECK/FOLD")
			return
		}
	}

	emptyCombination := hero.Hand.GetEmptyCombination(board)

	if emptyCombination.String() != "" {
		if emptyCombination.OverCards {
			fmt.Println("1 opponent: C-BET/FOLD or FOLD")
			fmt.Println("2+ opponents: CHECK/FOLD")
			return
		}
	}

	fmt.Println("monster draw: BET/ALL-IN or RERAISE")

	fmt.Println("draws: C-BET/FOLD or FOLD, on freeplay: CHECK/FOLD")

	fmt.Println("gotshot, 2+ opponents: CHECK/FOLD")
}

func (strategy Strategy) Turn() {
	hero := strategy.Table.Hero
	board := strategy.Table.Board

	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.String() != "" {
		if completedCombination.OverPair ||
			completedCombination.Three ||
			completedCombination.Triplet ||
			completedCombination.TwoPairs {
			fmt.Println("BET/ALL-IN or RERAISE")
			return
		}

		if completedCombination.TopPair {
			fmt.Println("C-BET/FOLD or FOLD")
			fmt.Println("freeplay: CHECK/FOLD")
			return
		}
	}

	fmt.Println("monster draw: BET/ALL-IN or RERAISE")
	fmt.Println("draw: CHECK/FOLD")
	return
}

func (strategy Strategy) River() {
	fmt.Println("monster, overpair, top pair: BET/RAISE or BET/CALL")
	fmt.Println("anything else: CHECK/FOLD")
}
