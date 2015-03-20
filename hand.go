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
	ImageSnippet
}

type Hand struct {
	LeftCard  Card
	RightCard Card
}

func (image Image) HandRecognize() string {
	hand := Hand{
		LeftCard: Card{
			ImageSnippet: ImageSnippet{
				Width:   cardWidth,
				Height:  cardHeight,
				OffsetX: handLeftCardOffsetX,
				OffsetY: handCardOffsetY,
			},
		},
		RightCard: Card{
			ImageSnippet: ImageSnippet{
				Width:   cardWidth,
				Height:  cardHeight,
				OffsetX: handRightCardOffsetX,
				OffsetY: handCardOffsetY,
			},
		},
	}

	leftCard, err := recognize(
		image.Crop(hand.LeftCard.ImageSnippet),
		cardSamples,
		handCompareThreshold,
	)

	if err != nil {
		return err.Error()
	}

	rightCard, err := recognize(
		image.Crop(hand.RightCard.ImageSnippet),
		cardSamples,
		handCompareThreshold,
	)

	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%v%v", leftCard, rightCard)
}
