package main

import (
	"fmt"
)

//@TODO: move to config
var (
	cardSamples          = "cards/*"
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
	ImageSnippet
}

type Hand struct {
	CardLeft  Card
	CardRight Card
}

func (image Image) HandRecognize() Hand {
	hand := Hand{
		CardLeft: Card{
			ImageSnippet: ImageSnippet{
				Width:   cardWidth,
				Height:  cardHeight,
				OffsetX: handLeftCardOffsetX,
				OffsetY: handCardOffsetY,
			},
		},
		CardRight: Card{
			ImageSnippet: ImageSnippet{
				Width:   cardWidth,
				Height:  cardHeight,
				OffsetX: handRightCardOffsetX,
				OffsetY: handCardOffsetY,
			},
		},
	}

	recognized, err := recognize(
		image.Crop(hand.CardLeft.ImageSnippet),
		cardSamples,
		handCompareThreshold,
	)

	if err == nil {
		hand.CardLeft.Value = fmt.Sprintf("%c", recognized[0])
		hand.CardLeft.Suit = fmt.Sprintf("%c", recognized[1])
	}

	recognized, err = recognize(
		image.Crop(hand.CardRight.ImageSnippet),
		cardSamples,
		handCompareThreshold,
	)

	if err == nil {
		hand.CardRight.Value = fmt.Sprintf("%c", recognized[0])
		hand.CardRight.Suit = fmt.Sprintf("%c", recognized[1])
	}

	return hand
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
	if hand.CardLeft.Value == "" || hand.CardRight.Value == "" {
		return ""
	}

	short := ""

	if cardStrength[hand.CardLeft.Value] > cardStrength[hand.CardRight.Value] {
		short = fmt.Sprintf("%s%s", hand.CardLeft.Value, hand.CardRight.Value)
	} else {
		short = fmt.Sprintf("%s%s", hand.CardRight.Value, hand.CardLeft.Value)
	}

	if hand.CardLeft.Suit == hand.CardRight.Suit {
		short += "s"
	}

	return short
}

func (hand Hand) String() string {
	return fmt.Sprintf("%s%s%s%s",
		hand.CardLeft.Value, hand.CardRight.Value,
		hand.CardLeft.Suit, hand.CardRight.Suit,
	)
}
