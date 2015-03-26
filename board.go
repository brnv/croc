package main

import (
	"fmt"
)

//@TODO: move to config
var (
	boardOffsetY          = 181
	boardCompareThreshold = 0.05
)

type Board struct {
	Cards []Card
}

func (image Image) BoardRecognize() Board {
	board := Board{}

	boardCards := board.GetBoardImageSnippets(
		[]int{264, 318, 372, 426, 480},
	)

	for _, boardCard := range boardCards {
		card, err := recognize(
			image.Crop(boardCard),
			cardSamples,
			boardCompareThreshold,
		)

		if err != nil {
			continue
		}

		board.Cards = append(board.Cards, Card{
			Value: fmt.Sprintf("%c", card[0]),
			Suit:  fmt.Sprintf("%c", card[1]),
		})
	}

	return board
}

func (board Board) GetBoardImageSnippets(offsets []int) []ImageSnippet {
	return getImageSnippets(
		cardWidth,
		cardHeight,
		boardOffsetY,
		offsets,
	)
}

func (board Board) GetStrontestBoardCard() string {
	strongestBoardCard := board.Cards[0].Value

	for _, boardCard := range board.Cards {
		if cardStrength[boardCard.Value] > cardStrength[strongestBoardCard] {
			strongestBoardCard = boardCard.Value
		}
	}

	return strongestBoardCard
}
