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

const avgStealSizePot = 9

func (strategy Strategy) PreflopRestealSituation() bool {
	heroPosition := strategy.Table.Hero.Position

	if positions[heroPosition] == "BB" &&
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

var pushHands = []string{
	"AA", "KK",
}
var raiseWaitPlayerHands = []string{
	"AK", "AKs",
	"QQ", "JJ", "TT", "99",
}

var raiseFoldHandsLatePosition = []string{
	"AQ", "AQs", "AJs",

	"AJ", "KQ",
	"KQs", "KJs", "ATs",
	"88", "77",

	"AT", "A9", "A8",
	"KJ", "QJ", "KT",

	"A9s", "A8s", "A7s", "A6s", "A5s",
	"KTs", "K9s", "QJs", "QTs", "JTs",
	"T9s",
	"66", "55", "44", "33", "22",

	"A7", "A6", "A5", "A4", "A3", "A2",
	"K9", "K8",
	"QT", "Q9",
	"JT",
	"T9", "98",

	"A4s", "A3s", "A2s",
	"K8s", "K7s", "K6s",
	"Q9s", "Q8s", "J9s",
	"98s", "87s", "76s", "65s",
}

var raiseFoldHands = map[string][]string{
	"EP": []string{
		"AQ", "AQs", "AJs",
	},
	"MP": []string{
		"AQ", "AQs", "AJs",

		"AJ", "KQ",
		"KQs", "KJs", "ATs",
		"88", "77",
	},
	"CO": []string{
		"AQ", "AQs", "AJs",

		"AJ", "KQ",
		"KQs", "KJs", "ATs",
		"88", "77",

		"AT", "A9", "A8",
		"KJ", "QJ", "KT",

		"A9s", "A8s", "A7s", "A6s", "A5s",
		"KTs", "K9s", "QJs", "QTs", "JTs",
		"T9s",
		"66", "55", "44", "33", "22",
	},
	"BU": raiseFoldHandsLatePosition,
	"SB": raiseFoldHandsLatePosition,
	"BB": raiseFoldHandsLatePosition,
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

var stealWaitPlayerHands = []string{
	"AK", "AKs",
	"QQ", "JJ", "TT",
}
var stealFoldHandsBUandSB = []string{
	"99", "88", "77", "66", "55", "44", "33", "22",
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
		"99", "88", "77", "66", "55", "44", "33", "22",
		"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9s", "A8s", "A7s",
		"KQ", "KQs", "KJ", "KJs", "KTs",
		"QJ", "QJs", "QTs",
		"JTs",
	},
	"BU": stealFoldHandsBUandSB,
	"SB": stealFoldHandsBUandSB,
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

var restealFoldHands = []string{
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9", "A9s",
	"99", "88",
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

var threeBetFoldMPHands = map[string][]string{
	"EP": []string{
		"QQ",
	},
	"MP": []string{
		"QQ", "TT",
		"AQ", "AQs", "AK", "AKs",
	},
	"LATER": []string{
		"QQ", "TT", "99", "88",
		"AQ", "AQs", "AK", "AKs",
	},
}

const potSaneLimitForThreeBet = 18

func (strategy *Strategy) PreflopThreeBetDecision() string {
	strategy.Messages = append(strategy.Messages, "3-bet")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, card := range pushHands {
		if hand == card {
			return "RAISE/ALL-IN"
		}
	}

	raiserPosition := positions[strategy.Table.GetFirstRaiserPosition()]

	strategy.Messages = append(strategy.Messages, "raiser in "+raiserPosition)

	for _, card := range threeBetFoldMPHands[strategyPositions[raiserPosition]] {
		if hand == card {
			if strategy.Table.Pot > potSaneLimitForThreeBet {
				return "MANUAL"
			}
			return "RAISE/FOLD"
		}
	}

	return "FOLD"
}

var contBetPairs = []string{
	"JJ", "TT",
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
