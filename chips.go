package main

//@TODO: move to config
var (
	chipsDigitSamples     = "chips_digits/*"
	chipsTypeSamples      = "chips_types/*"
	chipsTypeWidth        = 12
	chipsTypeHeight       = 13
	chipsTypeOffsetX      = 349
	chipsOffsetY          = 408
	chipsDigitWidth       = 9
	chipsDigitHeight      = 13
	chipsCompareThreshold = 0.2
)

type Chips struct {
	Number
}

func (image Image) HeroChipsRecognize() string {
	chips := Chips{
		Number: Number{
			Digits: []ImageSnippet{},
		},
	}

	chipsType, err := image.NumberTypeRecognize(
		ImageSnippet{
			Width:   chipsTypeWidth,
			Height:  chipsTypeHeight,
			OffsetX: chipsTypeOffsetX,
			OffsetY: chipsOffsetY,
		},
		chipsTypeSamples,
		chipsCompareThreshold,
	)

	if err != nil {
		log.Notice("%v", err.Error())
	}

	switch chipsType {
	case "3":
		chips.Number.Digits = chips.GetChipsImageSnippets(
			[]int{367, 376, 385},
		)
	case "4":
		chips.Number.Digits = chips.GetChipsImageSnippets(
			[]int{361, 373, 382, 391},
		)
	}

	chipsCount := ""

	for _, chipsDigit := range chips.Number.Digits {
		digit, err := recognize(
			image.Crop(chipsDigit),
			chipsDigitSamples,
			chipsCompareThreshold,
		)

		if err != nil {
			log.Notice("%v", err.Error())
		}
		chipsCount += digit
	}

	return chipsCount
}

func (chips Chips) GetChipsImageSnippets(offsets []int) []ImageSnippet {
	return getImageSnippets(
		chipsDigitWidth,
		chipsDigitHeight,
		chipsOffsetY,
		offsets,
	)
}
