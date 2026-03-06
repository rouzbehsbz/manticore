CREATE TABLE IF NOT EXISTS "characters" (
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    nickname VARCHAR(20) NOT NULL UNIQUE,
    level INTEGER NOT NULL CHECK (level >= 1),
    xp INTEGER NOT NULL CHECK (xp >= 0),
    vitality INTEGER NOT NULL CHECK (vitality >= 0),
    intelligence INTEGER NOT NULL CHECK (intelligence >= 0),
    willpower INTEGER NOT NULL CHECK (willpower >= 0),
    dexterity INTEGER NOT NULL CHECK (dexterity >= 0),
    spirit INTEGER NOT NULL CHECK (spirit >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE RESTRICT
)
