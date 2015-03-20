package main

//@TODO: move to config
var (
	potDigitSamples     = "pot_digits/*"
	potTypeSamples      = "pot_types/*"
	potTypeWidth        = 14
	potTypeHeight       = 13
	potTypeOffsetX      = 360
	potOffsetY          = 154
	potDigitWidth       = 9
	potDigitHeight      = 13
	potCompareThreshold = 0.05
)

type Pot struct {
	Number
}

func (image Image) PotRecognize() string {
	pot := Pot{
		Number: Number{
			Digits: []ImageSnippet{},
		},
	}

	potType, err := image.NumberTypeRecognize(
		ImageSnippet{
			Width:   potTypeWidth,
			Height:  potTypeHeight,
			OffsetX: potTypeOffsetX,
			OffsetY: potOffsetY,
		},
		potTypeSamples,
		potCompareThreshold,
	)

	if err != nil {
		return err.Error()
	}

	switch potType {
	case "3":
		pot.Number.Digits = pot.GetPotImageSnippets(
			[]int{403, 412, 421},
		)
	case "4":
		pot.Number.Digits = pot.GetPotImageSnippets(
			[]int{397, 409, 418, 427},
		)
	}

	potSize := ""

	for _, potDigit := range pot.Number.Digits {
		digit, err := recognize(
			image.Crop(potDigit),
			potDigitSamples,
			potCompareThreshold,
		)

		if err != nil {
			return err.Error()
		}
		potSize += digit
	}

	return potSize
}

func (pot Pot) GetPotImageSnippets(offsets []int) []ImageSnippet {
	return getImageSnippets(
		potDigitWidth,
		potDigitHeight,
		potOffsetY,
		offsets,
	)
}
