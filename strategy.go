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

func (strategy Strategy) PreflopRaiseSituation() bool {
	if !strategy.PotIsRaised() &&
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
	"QQ", "JJ",
}
var raiseFoldHandsLatePosition = []string{
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs", "A9s",
	"KQ", "KQs",
	"TT", "99", "88", "77",
}
var raiseFoldHands = map[string][]string{
	"EP": []string{
		"AQ", "AQs", "AJs",
		"TT",
	},
	"MP": []string{
		"AQ", "AQs", "AJ", "AJs", "ATs",
		"TT", "99", "88",
	},
	"CO": raiseFoldHandsLatePosition,
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

var threeBetFoldMPHands = []string{
	"AQ", "AQs", "AJ", "AJs", "AT", "ATs",
	"TT", "99", "88", "77",
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

	if strategy.Table.Pot > potSaneLimitForThreeBet {
		return "MANUAL"
	}

	for _, card := range raiseWaitPlayerHands {
		if hand == card {
			return "RAISE/MANUAL"
		}
	}

	for _, card := range threeBetFoldMPHands {
		if hand == card {
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

	if completedCombination.OverPair ||
		completedCombination.Three ||
		completedCombination.Triplet ||
		completedCombination.TwoPairs {
		return "MANUAL"
	}

	if completedCombination.StrongTopPair {
		return "MANUAL"
	}

	if completedCombination.TopPair {
		if strategy.Table.Pot <= 10 {
			return "FLOP CHECK/FOLD"
		} else if strategy.Table.Pot <= 35 {
			return "FLOP C-BET/FOLD"
		}
	}

	for _, card := range contBetPairs {
		if hand == card && strategy.Table.Pot <= 35 {
			return "FLOP C-BET/FOLD"
		}
	}

	emptyCombination := hero.Hand.GetEmptyCombination(board)

	if emptyCombination.String() != "" {
		if emptyCombination.OverCards {
			if strategy.Table.Pot <= 15 {
				return "FLOP C-BET/FOLD"
			}
		}
	}

	//@TODO: automate below logic

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

	return "MANUAL"
}

func (strategy *Strategy) TurnDecision() string {
	strategy.Messages = append(strategy.Messages, "turn")

	hero := strategy.Table.Hero
	board := strategy.Table.Board
	completedCombination := hero.Hand.GetCompletedCombination(board)

	if completedCombination.OverPair ||
		completedCombination.Three ||
		completedCombination.Triplet ||
		completedCombination.TwoPairs {
		return "MANUAL"
	}

	emptyCombination := hero.Hand.GetEmptyCombination(board)

	if emptyCombination.String() != "" {
		if emptyCombination.OverCards {
			return "CHECK/FOLD"
		}
	}

	//@TODO: automate below logic

	if completedCombination.TopPair {
		if strategy.Table.Pot <= 10 {
			// assuming we are on freeplay
			return "CHECK/FOLD"
		}

		return "MANUAL"
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
