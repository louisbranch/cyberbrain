package html

import (
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	"gitlab.com/luizbranco/srs/primitives"
	"gitlab.com/luizbranco/srs/web"
)

const CHECKED = "on"

func NewUserFromForm(form url.Values, auth primitives.Authenticator) (*primitives.User, error) {
	u := &primitives.User{}

	u.Name = form.Get("name")
	u.Email = form.Get("email")

	password := form.Get("password")

	if u.Name == "" {
		return nil, errors.New("user name cannot be empty")
	}

	if u.Email == "" {
		return nil, errors.New("user email cannot be empty")
	}

	if len(password) < 6 {
		return nil, errors.New("user password cannot be less than 6 characters")
	}

	hash, err := auth.Create(password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash user password")
	}

	u.PasswordHash = hash

	return u, nil
}

func NewDeckFromForm(form url.Values) (*primitives.Deck, error) {
	d := &primitives.Deck{}

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

func NewCardFromForm(deck primitives.Deck, form url.Values) (*primitives.Card, error) {
	c := &primitives.Card{
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

func NewTagFromForm(deck primitives.Deck, form url.Values) (*primitives.Tag, error) {
	t := &primitives.Tag{
		DeckID: deck.ID(),
		Name:   form.Get("name"),
	}
	return t, nil
}

func NewPracticeFromForm(deck primitives.Deck, tags []primitives.Tag, form url.Values, ub web.URLBuilder) (*primitives.Practice, error) {
	rounds := form.Get("rounds")
	n, err := strconv.Atoi(rounds)
	if err != nil {
		return nil, errors.Wrap(err, "invalid number of rounds")
	}

	p := &primitives.Practice{
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

		for _, t := range tags {
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

	p.PromptMode = primitives.PracticeImages
	p.GuessMode = primitives.PracticeDefinitions

	return p, nil
}
