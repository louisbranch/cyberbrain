package html

import (
	"net/url"
	"strconv"

	"github.com/luizbranco/srs"
	"github.com/pkg/errors"
)

func NewDeckFromForm(form url.Values) (*srs.Deck, error) {
	d := &srs.Deck{}

	d.Name = form.Get("name")
	d.Description = form.Get("description")
	d.ImageURL = form.Get("image_url")

	for _, cf := range form["fields"] {
		if cf != "" {
			d.Fields = append(d.Fields, cf)
		}
	}

	if d.Name == "" {
		return nil, errors.New("deck name cannot be empty")
	}

	if len(d.Fields) == 0 {
		return nil, errors.New("deck fields cannot be empty")
	}

	return d, nil
}

func NewCardFromForm(deck srs.Deck, form url.Values) (*srs.Card, error) {
	c := &srs.Card{
		DeckID: deck.ID(),
	}

	for _, f := range form["image_urls"] {
		if f != "" {
			c.ImageURLs = append(c.ImageURLs, f)
		}
	}

	for _, f := range form["audio_urls"] {
		if f != "" {
			c.AudioURLs = append(c.AudioURLs, f)
		}
	}

	for _, f := range form["definitions"] {
		if f != "" {
			c.Definitions = append(c.Definitions, f)
		}
	}

	if len(c.ImageURLs) == 0 {
		return nil, errors.New("card image cannot be empty")
	}

	if len(c.Definitions) != len(deck.Fields) {
		return nil, errors.New("card definition numbers must be the same as deck field definitions")
	}

	return c, nil
}

func NewTagFromForm(deckID srs.ID, form url.Values) (*srs.Tag, error) {
	t := &srs.Tag{
		DeckID: deckID,
		Name:   form.Get("name"),
	}
	return t, nil
}

func NewPracticeFromForm(deckID srs.ID, form url.Values) (*srs.Practice, error) {
	rounds := form.Get("rounds")
	n, err := strconv.Atoi(rounds)
	if err != nil {
		return nil, errors.Wrap(err, "invalid number of rounds")
	}

	p := &srs.Practice{
		DeckID: deckID,
		Rounds: n,
	}

	return p, nil
}
