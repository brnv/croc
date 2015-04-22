package main

func (strategy *Strategy) PreflopDecision() string {
	strategy.Messages = append(strategy.Messages, "preflop")

	if strategy.PreflopRaiseSituation() {
		return strategy.PreflopRaiseDecision()
	}

	if strategy.PreflopStealSituation() {
		return strategy.PreflopStealDecision()
	}

	if strategy.PreflopThreeBetSituation() {
		return strategy.PreflopThreeBetDecision()
	}

	if strategy.PreflopRestealSituation() {
		return strategy.PreflopRestealDecision()
	}

	return "MANUAL"
}

func (strategy Strategy) PreflopRaiseSituation() bool {
	if strategy.Table.PotIsRaised() {
		return false
	}

	if strategy.PreflopStealSituation() {
		return false
	}

	return true
}

func (strategy *Strategy) PreflopRaiseDecision() string {
	strategy.Messages = append(strategy.Messages, "raise")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, manualHand := range manualHands {
		if hand == manualHand {
			return "MANUAL"
		}
	}

	position := positions[strategy.Table.Hero.Position]

	for _, raiseFoldHand := range raiseFoldHands[position] {
		if hand == raiseFoldHand {
			return "RAISE/FOLD"
		}
	}

	if position == "BB" {
		return "CHECK"
	}

	return "FOLD"
}

func (strategy Strategy) PreflopStealSituation() bool {
	heroPosition := strategy.Table.Hero.Position

	defaultPotSize := strategy.Table.BigBlindSize + strategy.Table.BigBlindSize/2

	if strategy.Table.Pot != defaultPotSize {
		return false
	}

	if strategyPositions[positions[heroPosition]] != laterPosition {
		return false
	}

	return true
}

func (strategy *Strategy) PreflopStealDecision() string {
	strategy.Messages = append(strategy.Messages, "steal")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, manualHand := range manualHands {
		if hand == manualHand {
			return "MANUAL"
		}
	}

	position := positions[strategy.Table.Hero.Position]

	for _, card := range stealFoldHands[position] {
		if hand == card {
			return "RAISE/FOLD"
		}
	}

	if position == "BB" {
		return "CHECK"
	}

	return "FOLD"
}

func (strategy Strategy) PreflopRestealSituation() bool {
	heroPosition := strategy.Table.Hero.Position

	if !strategy.Table.PotIsRaised() {
		return false
	}

	if positions[heroPosition] != "SB" ||
		positions[heroPosition] != "BB" {
		return false
	}

	avgStealSizePot := 5 * strategy.Table.BigBlindSize

	if strategy.Table.Pot > avgStealSizePot {
		return false

	}
	if !strategy.Table.IsStealerPosition(
		strategy.Table.GetFirstRaiserPosition()) {
		return false
	}

	return true
}

func (strategy *Strategy) PreflopRestealDecision() string {
	strategy.Messages = append(strategy.Messages, "resteal")

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, manualHand := range manualHands {
		if hand == manualHand {
			return "MANUAL"
		}
	}

	for _, card := range restealFoldHands {
		if hand == card {
			return "RAISE/FOLD"
		}
	}

	return "FOLD"
}

func (strategy Strategy) PreflopThreeBetSituation() bool {
	if !strategy.Table.PotIsRaised() {
		return false
	}

	if strategy.PreflopRestealSituation() {
		return false
	}

	return true
}

func (strategy *Strategy) PreflopThreeBetDecision() string {
	strategy.Messages = append(strategy.Messages, "3-bet")

	raiserPosition := positions[strategy.Table.GetFirstRaiserPosition()]
	strategy.Messages = append(strategy.Messages, "raiser in "+raiserPosition)

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, manualHand := range manualHands {
		if hand == manualHand {
			return "MANUAL"
		}
	}

	potSaneLimitForThreeBet := 9 * strategy.Table.BigBlindSize

	for _, card := range threeBetHands[strategyPositions[raiserPosition]] {
		if hand == card {
			if strategy.Table.Pot > potSaneLimitForThreeBet {
				return "MANUAL"
			}

			return "RAISE/MANUAL"
		}
	}

	return "FOLD"
}
