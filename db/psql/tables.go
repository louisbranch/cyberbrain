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
			fields TEXT[] NOT NULL CHECK (cardinality(fields) > 0)
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS cards(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			definitions TEXT[] NOT NULL CHECK (cardinality(definitions) > 0),
			image_urls TEXT[],
			sound_urls TEXT[]
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
			mode INTEGER NOT NULL DEFAULT 0,
			field INTEGER NOT NULL DEFAULT 0,
			tag_id INTEGER,
			total_rounds INTEGER NOT NULL CHECK(total_rounds > 0),
			done BOOLEAN NOT NULL DEFAULT FALSE
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS rounds(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

			practice_id INTEGER NOT NULL REFERENCES practices ON DELETE RESTRICT,
			mode INTEGER NOT NULL DEFAULT 0,
			card_ids INTEGER[] NOT NULL CHECK (cardinality(card_ids) > 0),
			options TEXT[] NOT NULL CHECK (cardinality(options) > 0),
			guess TEXT,
			correct BOOLEAN NOT NULL DEFAULT FALSE,
			done BOOLEAN NOT NULL DEFAULT FALSE
		);
		`,
}
