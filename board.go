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

func (table *Table) BoardRecognize() {
	boardCards := table.Board.GetBoardImageSnippets(
		[]int{264, 318, 372, 426, 480},
	)

	_, err := recognize(
		table.Image.Crop(boardCards[0]),
		"/tmp/croc/first_board_card_empty_zoom",
		0.05,
	)

	if err == nil {
		return
	}

	_, err = recognize(
		table.Image.Crop(boardCards[0]),
		"/tmp/croc/first_board_card_empty",
		0.05,
	)

	if err == nil {
		return
	}

	for _, boardCard := range boardCards {
		card, err := recognize(
			table.Image.Crop(boardCard),
			cardSamples,
			boardCompareThreshold,
		)

		if err != nil {
			continue
		}

		table.Board.Cards = append(table.Board.Cards, Card{
			Value: fmt.Sprintf("%c", card[0]),
			Suit:  fmt.Sprintf("%c", card[1]),
		})
	}
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
