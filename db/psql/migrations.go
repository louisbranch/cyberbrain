package psql

import (
	"time"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/cyberbrain/primitives"
)

func createCardSchedules(db *Database) error {
	q := `SELECT cards.id, cards.deck_id FROM cards
	LEFT JOIN card_schedules ON card_schedules.card_id = cards.id
	WHERE card_schedules.card_id IS NULL;`

	rows, err := db.DB.Query(q)
	if err != nil {
		return errors.Wrapf(err, "failed to query records %q", q)
	}
	defer rows.Close()

	var schedules []*primitives.CardSchedule

	for rows.Next() {
		var cardID, deckID primitives.ID

		err := rows.Scan(&cardID, &deckID)
		if err != nil {
			return errors.Wrapf(err, "failed to scan records %q", q)
		}

		schedule := &primitives.CardSchedule{
			NextDate: time.Now(),
			DeckID:   deckID,
			CardID:   cardID,
		}

		schedules = append(schedules, schedule)
	}

	err = rows.Err()
	if err != nil {
		return errors.Wrapf(err, "failed to query records %q", q)
	}

	for _, s := range schedules {
		err := db.Create(s)
		if err != nil {
			return errors.Wrapf(err, "failed to create card schedule %v", s)
		}
	}

	return nil
}
