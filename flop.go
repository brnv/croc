package main

func (strategy *Strategy) PotStealIsOk() bool {
	if strategy.Table.Pot <= 5*strategy.Table.BigBlindSize {
		return false
	}

	if len(strategy.Table.Opponents) > 1 {
		return false
	}

	if strategy.Table.Pot >= 9*strategy.Table.BigBlindSize {
		return false
	}

	if !strategy.Table.CheckButtonIsVisible() {
		return false
	}

	strategy.Messages = append(strategy.Messages, "pot steal")

	return true
}

func (strategy *Strategy) FlopDecision() string {
	strategy.Messages = append(strategy.Messages, "flop")

	if strategy.IsGoodHand() {
		return "FLOP RAISE/ALL-IN"
	}

	hand := strategy.Table.Hero.Hand.ShortNotation()

	for _, flopAllInHand := range flopAllInHands {
		if hand == flopAllInHand {
			return "FLOP RAISE/ALL-IN"
		}
	}

	if !strategy.IsGoodHand() && strategy.PotStealIsOk() {
		return "FLOP C-BET/FOLD"
	}

	strategy.PrintReminders()

	return "CHECK/FOLD"
}
