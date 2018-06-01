package psql

var tableQueries = []string{
	`
		CREATE TABLE IF NOT EXISTS decks(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			name TEXT NOT NULL CHECK(name <> ''),
			description TEXT,
			image_url TEXT,
			card_fields TEXT[]
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS cards(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			definitions TEXT[],
			image_urls TEXT[],
			audio_urls TEXT[]
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS tags(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			name TEXT NOT NULL UNIQUE CHECK(name <> '')
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS card_tags(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			card_id INTEGER NOT NULL REFERENCES cards ON DELETE CASCADE,
			tag_id INTEGER NOT NULL REFERENCES tags ON DELETE CASCADE,
			UNIQUE (card_id, tag_id)
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS practices(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			state TEXT NOT NULL CHECK(state <> ''),
			rounds INTEGER NOT NULL CHECK(rounds > 0)
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS practice_rounds(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			card_id INTEGER NOT NULL REFERENCES cards ON DELETE CASCADE,
			practice_id INTEGER NOT NULL REFERENCES practices ON DELETE CASCADE,
			round INTEGER NOT NULL CHECK(round > 0),
			expect TEXT NOT NULL CHECK(expect <> ''),
			answer TEXT,
			correct BOOLEAN,
			UNIQUE (card_id, practice_id)
		);
		`,
}
