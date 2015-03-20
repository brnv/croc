package main

//@TODO: move to config
var (
	boardOffsetY          = 181
	boardCompareThreshold = 0.05
)

type Board struct {
	Cards []ImageSnippet
}

func (image Image) BoardRecognize() string {
	board := Board{
		Cards: []ImageSnippet{},
	}

	board.Cards = board.GetBoardImageSnippets(
		[]int{264, 318, 372, 426, 480},
	)

	boardCards := ""

	for _, boardCard := range board.Cards {
		card, err := recognize(
			image.Crop(boardCard),
			cardSamples,
			boardCompareThreshold,
		)

		if err != nil {
			return err.Error()
		}

		boardCards += card
	}

	return boardCards
}

func (board Board) GetBoardImageSnippets(offsets []int) []ImageSnippet {
	return getImageSnippets(
		cardWidth,
		cardHeight,
		boardOffsetY,
		offsets,
	)
}
