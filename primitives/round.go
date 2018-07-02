package primitives

import "time"

type Round struct {
	MetaID        ID        `db:"id"`
	MetaVersion   int       `db:"version"`
	MetaCreatedAt time.Time `db:"created_at"`
	MetaUpdatedAt time.Time `db:"updated_at"`

	PracticeID ID           `db:"practice_id"`
	CardIDs    []ID         `db:"card_ids"`
	PromptMode PracticeMode `db:"prompt_mode"`
	GuessMode  PracticeMode `db:"guess_mode"`
	Prompt     string       `db:"prompt"`
	Guess      string       `db:"guess"`
	Options    []string     `db:"options"`
	Answer     string       `db:"answer"`
	Correct    bool         `db:"correct"`
	Done       bool         `db:"done"`
	Caption    string       `db:"caption"`
}

func (r Round) ID() ID {
	return r.MetaID
}

func (r Round) Type() string {
	return "round"
}

func (r *Round) SetID(id ID) {
	r.MetaID = id
}

func (r *Round) SetVersion(v int) {
	r.MetaVersion = v
}

func (r *Round) SetCreatedAt(t time.Time) {
	r.MetaCreatedAt = t
}

func (r *Round) SetUpdatedAt(t time.Time) {
	r.MetaUpdatedAt = t
}

func (r *Round) GuessAnswer(answer string) bool {
	correct := r.Answer == answer

	r.Guess = answer
	r.Correct = correct
	r.Done = true

	return correct
}
