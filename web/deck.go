package web

type Deck struct {
	ID          uint   `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type Decks []Deck

func (d *Deck) Type() string {
	return "decks"
}

func (d *Deck) SetID(id uint) {
	d.ID = id
}

func (d *Decks) NewRecord() Record {
	return &Deck{}
}

func (d *Decks) Append(r Record) {
	deck := r.(*Deck)
	*d = append(*d, *deck)
}
