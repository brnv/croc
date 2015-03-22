package main

import (
	"errors"
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

var allInHands = []string{"AA", "KK"}

var latePositionsThreeBetAllInHands = []string{"QQ", "JJ", "AK"}
var threeBetAllInHands = map[string][]string{
	"EP": []string{"QQ"},
	"MP": []string{"QQ", "AK"},
	"CO": latePositionsThreeBetAllInHands,
	"BU": latePositionsThreeBetAllInHands,
	"SB": latePositionsThreeBetAllInHands,
	"BB": latePositionsThreeBetAllInHands,
}

var latePositionsThreeBetFoldHands = []string{
	"TT", "99", "88", "AQ", "AJ", "AT", "A9s",
}
var threeBetFoldHands = map[string][]string{
	"MP": []string{"JJ", "TT", "99", "AQ", "AJ", "ATs"},
	"CO": latePositionsThreeBetFoldHands,
	"BU": latePositionsThreeBetFoldHands,
	"SB": latePositionsThreeBetFoldHands,
	"BB": latePositionsThreeBetFoldHands,
}

func (strategy Strategy) Run() {
	err := strategy.Check()

	if err != nil {
		fmt.Println("bad input")
		return
	}

	fmt.Printf("players position is %s\n", positions[strategy.Table.Hero.Position])

	fmt.Printf("players hand is %s\n", strategy.Table.Hero.Hand.ShortNotification())

	if strategy.Table.Board == "" {
		strategy.Preflop()
	} else {
		fmt.Println("flop, turn or river strategy")
	}
}

func (strategy Strategy) Preflop() {

	if strategy.OpponentsWereRaising() {
		strategy.PreflopThreeBetStrategy()
	} else {
		strategy.PreflopRaiseStrategy()
	}
}

func (strategy Strategy) OpponentsWereRaising() bool {
	limpTotalSize := 0

	for _, limper := range strategy.Table.Opponents {
		limpTotalSize += limper.LimpSize

	}

	if limpTotalSize == strategy.Table.Pot {
		return false
	}

	return true
}

func (strategy Strategy) PreflopThreeBetStrategy() {
	fmt.Println("preflop 3-bet strategy decision is")

	hand := strategy.Table.Hero.Hand.ShortNotification()

	for _, card := range allInHands {
		if hand == card {
			fmt.Println("ALL-IN if many raises, 3-BET if one")
			return
		}
	}

	for position, cards := range threeBetAllInHands {
		for _, card := range cards {
			if hand == card {
				fmt.Printf(
					"3-BET and ALL-IN after 4-BET if opponent in %s", position,
				)
				return
			}

		}
	}

	for position, cards := range threeBetFoldHands {
		for _, card := range cards {
			if hand == card {
				fmt.Printf(
					"3-BET and FOLD after 4-BET if opponent in %s", position,
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

	hand := strategy.Table.Hero.Hand.ShortNotification()

	for _, element := range raiseFoldHands[position] {
		if element == hand {
			fmt.Println("RAISE and FOLD after 3-bet")
			return
		}
	}

	for _, element := range raisePushHands[position] {
		if element == hand {
			fmt.Println("RAISE and ALL-IN after 3-bet")
			return
		}
	}

	if position == "BB" {
		fmt.Println("CHECK and FOLD after opponents raise")
	} else {
		fmt.Println("FOLD")
	}
}

func (strategy Strategy) Check() error {
	hand := strategy.Table.Hero.Hand.ShortNotification()

	if hand == "" {
		return errors.New("No hand provided")
	}

	return nil
}
