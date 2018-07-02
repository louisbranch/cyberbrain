package psql

var tableQueries = []string{
	`
		CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			name TEXT NOT NULL CHECK(name <> ''),
			email TEXT NOT NULL UNIQUE CHECK(email <> ''),
			password_hash TEXT NOT NULL CHECK(password_hash <> ''),
			image_url TEXT
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS sessions(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			user_id INTEGER NOT NULL REFERENCES users ON DELETE CASCADE
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS decks(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			user_id INTEGER NOT NULL REFERENCES users ON DELETE CASCADE,
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
			image_url TEXT NOT NULL CHECK(image_url <> ''),
			sound_url TEXT
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS tags(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			name TEXT NOT NULL CHECK(name <> ''),
			UNIQUE (deck_id, name)
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
			prompt_mode INTEGER NOT NULL DEFAULT 0,
			guess_mode INTEGER NOT NULL DEFAULT 0,
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
			card_ids INTEGER[] NOT NULL CHECK (cardinality(card_ids) > 0),
			prompt_mode INTEGER NOT NULL DEFAULT 0,
			guess_mode INTEGER NOT NULL DEFAULT 0,
			options TEXT[],
			prompt TEXT NOT NULL CHECK(prompt <> ''),
			guess TEXT,
			answer TEXT NOT NULL CHECK(answer <> ''),
			correct BOOLEAN NOT NULL DEFAULT FALSE,
			done BOOLEAN NOT NULL DEFAULT FALSE
		);
		`,
	`
		CREATE TABLE IF NOT EXISTS jobs(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			run_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			name TEXT NOT NULL CHECK(name <> ''),
			state TEXT NOT NULL CHECK(state <> ''),
			args BYTEA,
			error TEXT,
			tries INTEGER NOT NULL DEFAULT 0
		);
		`,

	` ALTER TABLE cards ADD COLUMN IF NOT EXISTS caption TEXT;`,
	` ALTER TABLE rounds ADD COLUMN IF NOT EXISTS caption TEXT;`,
	`
		CREATE TABLE IF NOT EXISTS card_schedules(
			id SERIAL PRIMARY KEY,
			version INTEGER NOT NULL DEFAULT 1,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			next_date TIMESTAMPTZ NOT NULL,
			deck_id INTEGER NOT NULL REFERENCES decks ON DELETE CASCADE,
			card_id INTEGER NOT NULL REFERENCES cards ON DELETE CASCADE,
			current_score INTEGER NOT NULL DEFAULT 0,
			max_score INTEGER NOT NULL DEFAULT 0
		);
		`,
}
