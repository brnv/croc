package main

import "fmt"

//@TODO: move to config
var (
	cardSamples          = "/tmp/croc/cards/*"
	cardWidth            = 46
	cardHeight           = 30
	handLeftCardOffsetX  = 346
	handRightCardOffsetX = 396
	handCardOffsetY      = 341 // 9 players, 340 for 6 players
	handCompareThreshold = 0.05
)

type Card struct {
	Value string
	Suit  string
}

type HandCard struct {
	Card
	ImageSnippet
}

type Hand struct {
	Cards []HandCard
}

func (table *Table) HandRecognize() {
	table.Hero.Hand = Hand{
		Cards: []HandCard{
			HandCard{ImageSnippet: ImageSnippet{
				Width:   cardWidth,
				Height:  cardHeight,
				OffsetX: handLeftCardOffsetX,
				OffsetY: handCardOffsetY,
			}},
			HandCard{ImageSnippet: ImageSnippet{
				Width:   cardWidth,
				Height:  cardHeight,
				OffsetX: handRightCardOffsetX,
				OffsetY: handCardOffsetY,
			}},
		}}

	recognized, err := recognize(
		table.Image.Crop(table.Hero.Hand.Cards[0].ImageSnippet),
		cardSamples,
		handCompareThreshold,
	)

	if err == nil {
		table.Hero.Hand.Cards[0].Value = fmt.Sprintf("%c", recognized[0])
		table.Hero.Hand.Cards[0].Suit = fmt.Sprintf("%c", recognized[1])
	}

	recognized, err = recognize(
		table.Image.Crop(table.Hero.Hand.Cards[1].ImageSnippet),
		cardSamples,
		handCompareThreshold,
	)

	if err == nil {
		table.Hero.Hand.Cards[1].Value = fmt.Sprintf("%c", recognized[0])
		table.Hero.Hand.Cards[1].Suit = fmt.Sprintf("%c", recognized[1])
	}
}

var cardStrength = map[string]int{
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"T": 10,
	"J": 11,
	"Q": 12,
	"K": 13,
	"A": 14,
}

func (hand Hand) ShortNotation() string {
	if len(hand.Cards) == 0 {
		return ""
	}

	if hand.Cards[0].Value == "" || hand.Cards[1].Value == "" {
		return ""
	}

	short := ""

	if cardStrength[hand.Cards[0].Value] > cardStrength[hand.Cards[1].Value] {
		short = fmt.Sprintf("%s%s", hand.Cards[0].Value, hand.Cards[1].Value)
	} else {
		short = fmt.Sprintf("%s%s", hand.Cards[1].Value, hand.Cards[0].Value)
	}

	if hand.Cards[0].Suit == hand.Cards[1].Suit {
		short += "s"
	}

	return short
}

func (hand Hand) String() string {
	return fmt.Sprintf("%s%s%s%s",
		hand.Cards[0].Value, hand.Cards[0].Suit,
		hand.Cards[1].Value, hand.Cards[1].Suit,
	)
}

type CompletedCombination struct {
	OverPair bool
	Three    bool

	TopPair  bool
	TwoPairs bool
	Triplet  bool
}

func (combination CompletedCombination) String() string {
	if combination.Triplet || combination.Three {
		return "three of a kind"
	}

	if combination.TwoPairs {
		return "two pairs"
	}

	if combination.OverPair {
		return "over pair"
	}

	if combination.TopPair {
		return "top pair"
	}

	return ""
}

type DrawCombination struct {
	Oesd          bool
	Gotshot       bool
	DoubleGotshot bool
	FlushDraw     bool
	MonsterDraw   bool
}

type EmptyCombination struct {
	OverCards bool
	Trash     bool
}

func (combination EmptyCombination) String() string {
	if combination.OverCards {
		return "over cards"
	}

	if combination.Trash {
		return "trash"
	}

	return ""
}

func (combination CompletedCombination) CheckTopPair(hand Hand, board Board) bool {
	strongestBoardCard := board.GetStrontestBoardCard()

	for _, handCard := range hand.Cards {
		if strongestBoardCard == handCard.Value {
			return true
		}
	}

	return false

}

func checkOverCards(hand Hand, board Board) bool {
	strongestBoardCard := board.GetStrontestBoardCard()

	for _, handCard := range hand.Cards {
		if cardStrength[strongestBoardCard] >= cardStrength[handCard.Value] {
			return false
		}
	}

	return true
}

func (combination CompletedCombination) CheckTwoPairs(hand Hand, board Board) bool {
	pairsCount := 0

	for _, handCard := range hand.Cards {
		for _, boardCard := range board.Cards {
			if handCard.Value == boardCard.Value {
				pairsCount++
				break
			}
		}
	}

	if pairsCount == 2 {
		return true
	}

	return false
}

func (combination CompletedCombination) CheckThree(hand Hand, board Board) bool {
	for _, boardCard := range board.Cards {
		if boardCard.Value == hand.Cards[0].Value {
			return true
		}
	}

	return false
}

func (combination CompletedCombination) CheckTriplet(hand Hand, board Board) bool {
	count := 1
	for _, handCard := range hand.Cards {
		count = 1

		for _, boardCard := range board.Cards {
			if boardCard.Value == handCard.Value {
				count++
			}
		}

		if count == 3 {
			return true
		}
	}
	return false
}

func (hand Hand) GetCompletedCombination(board Board) CompletedCombination {
	combination := CompletedCombination{}

	if hand.Cards[0].Value == hand.Cards[1].Value {
		combination.OverPair = checkOverCards(hand, board)
		combination.Three = combination.CheckThree(hand, board)
	} else {
		combination.TopPair = combination.CheckTopPair(hand, board)
		combination.TwoPairs = combination.CheckTwoPairs(hand, board)
		combination.Triplet = combination.CheckTriplet(hand, board)
	}

	return combination
}

func (hand Hand) GetEmptyCombination(board Board) EmptyCombination {
	combination := EmptyCombination{}

	if checkOverCards(hand, board) {
		combination.OverCards = true
	}

	return combination
}
