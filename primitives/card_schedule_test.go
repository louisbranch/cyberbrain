package primitives

import (
	"testing"

	"gitlab.com/luizbranco/cyberbrain/test"
)

func TestNewCardSchedule(t *testing.T) {
	deckID := ID(1)
	cardID := ID(2)

	schedule := NewCardSchedule(deckID, cardID)

	test.Equal(t, "deck id", schedule.DeckID, deckID)
	test.Equal(t, "card id", schedule.CardID, cardID)
	test.Equal(t, "next date", schedule.NextDate, days(1))
}

func TestCardSchedule_Reschedule(t *testing.T) {

	tcs := []struct {
		scenario  string
		input     CardSchedule
		correct   CardSchedule
		incorrect CardSchedule
	}{
		{
			scenario: "no previous scores",
			correct: CardSchedule{
				CurrentScore: 1,
				MaxScore:     1,
				NextDate:     days(2),
			},
			incorrect: CardSchedule{
				CurrentScore: -1,
				MaxScore:     0,
				NextDate:     days(1),
			},
		},
		{
			scenario: "winning streak",
			input: CardSchedule{
				CurrentScore: 5,
				MaxScore:     10,
			},
			correct: CardSchedule{
				CurrentScore: 6,
				MaxScore:     11,
				NextDate:     days(36),
			},
			incorrect: CardSchedule{
				CurrentScore: -1,
				MaxScore:     10,
				NextDate:     days(1),
			},
		},
		{
			scenario: "losing streak",
			input: CardSchedule{
				CurrentScore: -5,
				MaxScore:     10,
			},
			correct: CardSchedule{
				CurrentScore: 1,
				MaxScore:     10,
				NextDate:     days(2),
			},
			incorrect: CardSchedule{
				CurrentScore: -6,
				MaxScore:     10,
				NextDate:     days(1),
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {

			correct := tc.input

			correct.Reschedule(true)

			test.Equal(t, "current score", tc.correct.CurrentScore, correct.CurrentScore)
			test.Equal(t, "max score", tc.correct.MaxScore, correct.MaxScore)
			test.Equal(t, "next date", tc.correct.NextDate, correct.NextDate)

			incorrect := tc.input

			incorrect.Reschedule(false)

			test.Equal(t, "current score", tc.incorrect.CurrentScore, incorrect.CurrentScore)
			test.Equal(t, "max score", tc.incorrect.MaxScore, incorrect.MaxScore)
			test.Equal(t, "next date", tc.incorrect.NextDate, incorrect.NextDate)

		})
	}
}
