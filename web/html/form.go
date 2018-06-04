package html

import (
	"net/url"
	"strconv"

	"github.com/luizbranco/srs"
	"github.com/luizbranco/srs/web"
	"github.com/pkg/errors"
)

const CHECKED = "on"

func NewDeckFromForm(form url.Values) (*srs.Deck, error) {
	d := &srs.Deck{}

	d.UserID = 1 // FIXME
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

	for _, f := range form["sound_urls"] {
		if f != "" {
			c.SoundURLs = append(c.SoundURLs, f)
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

func NewTagFromForm(deck srs.Deck, form url.Values) (*srs.Tag, error) {
	t := &srs.Tag{
		DeckID: deck.ID(),
		Name:   form.Get("name"),
	}
	return t, nil
}

func NewPracticeFromForm(deck srs.Deck, form url.Values, ub web.URLBuilder) (*srs.Practice, error) {
	rounds := form.Get("rounds")
	n, err := strconv.Atoi(rounds)
	if err != nil {
		return nil, errors.Wrap(err, "invalid number of rounds")
	}

	p := &srs.Practice{
		DeckID:      deck.ID(),
		TotalRounds: n,
	}

	tagID := form.Get("tags")

	if tagID != "" {
		id, err := ub.ParseID(tagID)
		if err != nil {
			errors.Wrap(err, "invalid tag id")
		}

		found := false

		for _, t := range deck.Tags {
			if t.ID() == id {
				found = true
				break
			}
		}

		if !found {
			return nil, errors.Errorf("invalid tag id %s", tagID)
		}

		p.TagID = &id
	}

	p.PromptMode = srs.PracticeImages
	p.GuessMode = srs.PracticeDefinitions

	return p, nil
}
