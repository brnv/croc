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

	fmt.Printf("Hero have %s\n", strategy.Table.Hero.Hand)

	if len(strategy.Table.Board.Cards) == 0 {
		strategy.Preflop()
	} else {
		boardCardsNum := len(strategy.Table.Board.Cards)

		if boardCardsNum == 3 {
			strategy.Flop()
		} else if boardCardsNum == 4 {
			strategy.Turn()
		} else if boardCardsNum == 5 {
			fmt.Println("river strategy decision is")
			fmt.Println("monster, overpair, top pair: BET/RAISE or BET/CALL")
			fmt.Println("anything else: CHECK/FOLD")
		}
	}
}

func (strategy Strategy) Turn() {
	fmt.Println("turn strategy decision is")

	hero := strategy.Table.Hero
	board := strategy.Table.Board

	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.String() != "" {
		if completedCombination.OverPair {
			fmt.Println("overpair: BET and ALL-IN on reraise or RERAISE opponents raise")
			return
		}

		if completedCombination.Three ||
			completedCombination.Triplet ||
			completedCombination.TwoPairs {
			fmt.Println("monster: BET and ALL-IN on reraise or RERAISE opponents raise")
			return
		}

		if completedCombination.TopPair {
			fmt.Println("top pair: C-BET and FOLD on reraise or FOLD on opponents raise")
			return
		}
	}

	fmt.Println("Monster draw: BET and ALL-IN on reraise or RERAISE opponents raise")
	fmt.Println("anything else: CHECK/FOLD")
	return
}

func (strategy Strategy) Flop() {
	if strategy.PotIsRaised() {
		strategy.Postflop()
	} else {
		fmt.Println("implement freeplay strategy")
	}
}

func (strategy Strategy) Postflop() {
	fmt.Println("postflop strategy decision is")

	hero := strategy.Table.Hero
	board := strategy.Table.Board

	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.String() != "" {
		if completedCombination.OverPair {
			fmt.Println("overpair: BET and ALL-IN on reraise or RERAISE opponents raise")
			return
		}

		if completedCombination.Three ||
			completedCombination.Triplet ||
			completedCombination.TwoPairs {
			fmt.Println("monster: BET and ALL-IN on reraise or RERAISE opponents raise")
			return
		}

		if completedCombination.TopPair {
			fmt.Println("top pair: C-BET and FOLD on reraise or FOLD on opponents raise")
			return
		}
	}

	emptyCombination := hero.Hand.GetEmptyCombination(board)

	if emptyCombination.String() != "" {
		if emptyCombination.OverCards {
			fmt.Println("overcards")
			fmt.Println("Dry board and 1 opponent: C-BET and FOLD on reraise or FOLD on opponents raise")
			fmt.Println("Draw board or 2+ opponents: CHECK/FOLD")
			return
		}
	}

	fmt.Println("Monster draw: BET and ALL-IN on reraise or RERAISE opponents raise")
	fmt.Println("gotshot, trash on 2+ opponents: CHECK/FOLD")
	fmt.Println("anything else: C-BET and FOLD on reraise or FOLD on opponents raise")
	return
}

func (strategy Strategy) Preflop() {
	heroPosition := strategy.Table.Hero.Position

	fmt.Printf("Hero is %s\n", positions[heroPosition])

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
	fmt.Println("preflop steal strategy decision is")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range stealAllInHands {
		if hand == card {
			fmt.Println("STEAL and ALL-IN on opponents resteal")
			return
		}
	}

	for _, card := range stealFoldHands[position] {
		if hand == card {
			fmt.Println("STEAL and FOLD on opponents resteal")
			return
		}
	}

	fmt.Println("FOLD")
}

func (strategy Strategy) PreflopReStealStrategy() {
	fmt.Println("preflop re-steal strategy decision is")

	position := positions[strategy.Table.Hero.Position]

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range stealAllInHands {
		if hand == card {
			fmt.Println("RESTEAL and ALL-IN on opponents resteal")
			return
		}
	}

	for _, card := range reStealFoldHands[position] {
		if hand == card {
			fmt.Println("RESTEAL and FOLD on opponents resteal")
			return
		}
	}

	fmt.Println("FOLD")
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
					"3-BET and ALL-IN after %s position opponents 4-BET\n",
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
					"3-BET and FOLD after %s position opponents 4-BET\n",
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
		fmt.Println("implement hand equity and decide to")
		fmt.Println("LIMP or FOLD")
		return
	}

	if position == "BB" {
		fmt.Println("CHECK")
		return
	}

	fmt.Println("FOLD")
}

func (strategy Strategy) CheckInput() error {
	hand := strategy.Table.Hero.Hand.ShortNotation()

	if hand == "" {
		return errors.New("No hand provided")
	}

	return nil
}
