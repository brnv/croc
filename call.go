package main

//@TODO: move to config
var (
	callDigitSamples           = "/tmp/croc/call_digits/*"
	callTypeSamples            = "/tmp/croc/call_types/*"
	raiseTypeSamples           = "/tmp/croc/raise_types/*"
	callTypeWidth              = 10
	callTypeHeight             = 16
	callTypeOffsetX            = 557
	raiseTypeOffsetX           = 692
	callOffsetY                = 518
	callDigitWidth             = 11
	callDigitHeight            = 16
	callCompareThreshold       = 0.26
	callDigitsTwoCallButton    = []int{575, 586}
	callDigitsThreeCallButton  = []int{570, 581, 592}
	callDigitsFourCallButton   = []int{}
	callDigitsTwoRaiseButton   = []int{711, 722}
	callDigitsThreeRaiseButton = []int{705, 717, 729}
	callDigitsFourRaiseButton  = []int{700, 711, 723, 735}
)

type Call struct {
	Number
}

func (image Image) CallRecognize() string {
	call := Call{
		Number: Number{
			Digits: []ImageSnippet{},
		},
	}

	callType, err := image.NumberTypeRecognize(
		ImageSnippet{
			Width:   callTypeWidth,
			Height:  callTypeHeight,
			OffsetX: callTypeOffsetX,
			OffsetY: callOffsetY,
		},
		callTypeSamples,
		callCompareThreshold,
	)

	if err != nil {
		return err.Error()
	}

	callDigitsTwo := callDigitsTwoCallButton
	callDigitsThree := callDigitsThreeCallButton
	callDigitsFour := callDigitsFourCallButton

	if callType == "" {
		callType, err = image.NumberTypeRecognize(
			ImageSnippet{
				Width:   callTypeWidth,
				Height:  callTypeHeight,
				OffsetX: raiseTypeOffsetX,
				OffsetY: callOffsetY,
			},
			raiseTypeSamples,
			callCompareThreshold,
		)

		if err != nil {
			return err.Error()
		}

		callDigitsTwo = callDigitsTwoRaiseButton
		callDigitsThree = callDigitsThreeRaiseButton
		callDigitsFour = callDigitsFourRaiseButton
	}

	switch callType {
	case "2":
		call.Number.Digits = call.GetCallImageSnippets(
			callDigitsTwo,
		)
	case "3":
		call.Number.Digits = call.GetCallImageSnippets(
			callDigitsThree,
		)
	case "3_1":
		call.Number.Digits = call.GetCallImageSnippets(
			[]int{569, 580, 592},
		)
	case "4":
		call.Number.Digits = call.GetCallImageSnippets(
			callDigitsFour,
		)
	}

	callSize := ""

	for _, callDigit := range call.Number.Digits {
		digit, err := recognize(
			image.Crop(callDigit),
			callDigitSamples,
			callCompareThreshold,
		)

		if err != nil {
			log.Notice("%v", err.Error())
			return callSize
		}

		callSize += digit
	}

	return callSize
}

func (call Call) GetCallImageSnippets(offsets []int) []ImageSnippet {
	return getImageSnippets(
		callDigitWidth,
		callDigitHeight,
		callOffsetY,
		offsets,
	)
}
