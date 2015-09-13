package main

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

		_, offset := getSubimageOffset(
			"/tmp/croc/cards",
			table.Image.Crop(boardCard),
		)

		if offset >= 0 {
			table.Board.Cards = append(table.Board.Cards, Card{
				Value: GetValueByOffset(offset),
				Suit:  GetSuitByOffset(offset),
			})
		}
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

func (board Board) GetStrongestBoardCard() string {
	strongestBoardCard := board.Cards[0].Value

	for _, boardCard := range board.Cards {
		if cardStrength[boardCard.Value] > cardStrength[strongestBoardCard] {
			strongestBoardCard = boardCard.Value
		}
	}

	return strongestBoardCard
}
