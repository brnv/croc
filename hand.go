package main

//@TODO: move to config
var (
	cardSamples          = "cards/*"
	cardWidth            = 46
	cardHeight           = 30
	handLeftCardOffsetX  = 346
	handRightCardOffsetX = 396
	handCardOffsetY      = 341 // 9 players, 340 for 6 players
)

type Card struct {
	ImageSnippet
}

type Hand struct {
	LeftCard  Card
	RightCard Card
}

func (image Image) HandRecognize() {
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

	leftCard, err := recognize(image.Crop(hand.LeftCard.ImageSnippet), cardSamples)
	if err != nil {
		log.Notice("%v", err.Error())
	}
	log.Notice("%v\n", leftCard)

	rightCard, err := recognize(image.Crop(hand.RightCard.ImageSnippet), cardSamples)
	if err != nil {
		log.Notice("%v", err.Error())
	}
	log.Notice("%v\n", rightCard)
}
