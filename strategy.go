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
	noLimpPotSize = 3
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
	"AA", "KK", "QQ", "JJ",
	"AK", "AKs",
}
var stealFoldHandsBUandSB = []string{
	"JJ", "TT", "99", "88", "77", "66", "55", "44", "33", "22",
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9", "A9s",
	"A8", "A8s", "A7", "A7s", "A6s", "A5s", "A4s", "A3s", "A2s",
	"KQ", "KQs", "KJ", "KJs", "KT", "KTs",
	"QJ", "QJs", "QT", "QTs",
	"JT", "JTs",
	"T9", "T9s",
	"98s",
	"87s",
	"76s",
}
var stealFoldHands = map[string][]string{
	"CO": []string{
		"TT", "99", "88", "77", "66", "55", "44", "33", "22",
		"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9s", "A8s", "A7s",
		"KQ", "KQs", "KJ", "KJs", "KTs",
		"QJ", "QJs", "QTs",
		"JTs",
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

func (strategy *Strategy) Run() string {
	err := strategy.CheckInput()

	if err != nil {
		return err.Error()
	}

	boardCardsCount := len(strategy.Table.Board.Cards)

	if boardCardsCount == 0 {
		return strategy.Preflop()
	}

	if boardCardsCount == 3 {
		strategy.Flop()
		return ""
	} else if boardCardsCount == 4 {
		strategy.Turn()
		return ""
	}

	strategy.River()
	return ""
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

	for _, limper := range strategy.Table.Limpers {
		limpTotalSize += limper.BetSize
	}

	if limpTotalSize != strategy.Table.Pot {
		return true
	}

	return false
}

func (strategy *Strategy) Preflop() string {
	strategy.Messages = append(strategy.Messages, "preflop")

	decision := ""

	heroPosition := strategy.Table.Hero.Position

	if !strategy.PotIsRaised() {
		if strategyPositions[positions[heroPosition]] == laterPosition &&
			strategy.Table.Pot == noLimpPotSize {
			decision = strategy.PreflopStealStrategy()
		} else {
			decision = strategy.PreflopRaiseStrategy()
		}
	} else {
		if strategyPositions[positions[heroPosition]] == laterPosition {
			decision += fmt.Sprintf("%s\n", strategy.PreflopReStealStrategy())
		}
		if decision == "FOLD\n" {
			decision = strategy.PreflopThreeBetStrategy()
		} else {
			decision += strategy.PreflopThreeBetStrategy()
		}
	}

	return decision
}

func (strategy *Strategy) PreflopStealStrategy() string {
	strategy.Messages = append(strategy.Messages, "steal")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range stealAllInHands {
		if hand == card {
			return "STEAL/ALL-IN"
		}
	}

	for _, card := range stealFoldHands[position] {
		if hand == card {
			return "STEAL/FOLD"
		}
	}

	if position == "BB" {
		return "CHECK"
	}

	return "FOLD"
}

func (strategy *Strategy) PreflopRaiseStrategy() string {
	strategy.Messages = append(strategy.Messages, "raise")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, element := range raisePushHands[position] {
		if element == hand {
			return "RAISE/ALL-IN"
		}
	}

	for _, element := range raiseFoldHands[position] {
		if element == hand {
			return "RAISE/FOLD"
		}
	}

	if position == "SB" {
		//return "LIMP or FOLD"
	}

	if position == "BB" {
		return "CHECK"
	}

	return "FOLD"
}

func (strategy *Strategy) PreflopReStealStrategy() string {
	strategy.Messages = append(strategy.Messages, "resteal")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range stealAllInHands {
		if hand == card {
			return "RESTEAL/ALL-IN"
		}
	}

	for _, card := range reStealFoldHands[position] {
		if hand == card {
			return "RESTEAL/FOLD"
		}
	}

	return "FOLD"
}

func (strategy *Strategy) PreflopThreeBetStrategy() string {
	strategy.Messages = append(strategy.Messages, "3-bet")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range allInHands {
		if hand == card {
			//return "ALL-IN on 2+ raises or 3-BET on 1"
			return "3-BET/ALL-IN"
		}
	}

	for position, cards := range threeBetAllInHands {
		for _, card := range cards {
			if hand == card {
				return fmt.Sprintf(
					"3-BET/ALL-IN if raiser >= %s",
					strategyPositions[position],
				)
			}
		}
	}

	for position, cards := range threeBetFoldMPHands {
		for _, card := range cards {
			if hand == card {
				return fmt.Sprintf(
					"3-BET/FOLD if raiser >= %s",
					strategyPositions[position],
				)
			}
		}
	}

	for position, cards := range threeBetFoldLATERHands {
		for _, card := range cards {
			if hand == card {
				return fmt.Sprintf(
					"3-BET/FOLD if raiser >= %s",
					strategyPositions[position],
				)
			}
		}
	}

	return "FOLD"
}

func (strategy *Strategy) Flop() {
	fmt.Println("FLOP")

	strategy.Messages = append(strategy.Messages, "flop")
	hero := strategy.Table.Hero
	board := strategy.Table.Board

	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.String() != "" {
		if completedCombination.OverPair ||
			completedCombination.Three ||
			completedCombination.Triplet ||
			completedCombination.TwoPairs {
			fmt.Println("BET/ALL-IN or RERAISE;")
			return
		}

		if completedCombination.TopPair {
			fmt.Println("C-BET/FOLD or FOLD;")
			fmt.Println("freeplay: CHECK/FOLD;")
			return
		}
	}

	emptyCombination := hero.Hand.GetEmptyCombination(board)

	if emptyCombination.String() != "" {
		if emptyCombination.OverCards {
			fmt.Println("overcards: 1 opponent: C-BET/FOLD or FOLD;")
		}
	}

	fmt.Println("monster draw: BET/ALL-IN or RERAISE;")

	fmt.Println("draws: C-BET/FOLD or FOLD, on freeplay: CHECK/FOLD;")

	fmt.Println(
		fmt.Sprintf(
			"draws: if win_size / call_size / " +
				"[monster/3, flush/4, oesd/5, overcards/7, pair/8] > 1:" +
				" CALL;",
		),
	)

	fmt.Println("gotshot, 2+ opponents: CHECK/FOLD;")
}

func (strategy *Strategy) Turn() {
	strategy.Messages = append(strategy.Messages, "turn")
	hero := strategy.Table.Hero
	board := strategy.Table.Board

	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.String() != "" {
		if completedCombination.OverPair ||
			completedCombination.Three ||
			completedCombination.Triplet ||
			completedCombination.TwoPairs {
			fmt.Println("BET/ALL-IN or RERAISE;")
			return
		}

		if completedCombination.TopPair {
			fmt.Println("C-BET/FOLD or FOLD;")
			fmt.Println("freeplay: CHECK/FOLD;")
			return
		}
	}

	fmt.Println("monster draw: BET/ALL-IN or RERAISE;")

	fmt.Println("draw: CHECK/FOLD;")

	fmt.Println(
		fmt.Sprintf(
			"draws: if win_size / call_size / " +
				"[monster/1, flush/2, oesd/2, overcards/3, pair/4] > 1:" +
				" CALL;",
		),
	)

	return
}

func (strategy *Strategy) River() {
	strategy.Messages = append(strategy.Messages, "river")
	fmt.Println("monster, overpair, top pair: BET/RAISE or BET/CALL;")
	fmt.Println("anything else: CHECK/FOLD;")
}
