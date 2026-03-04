CREATE TABLE IF NOT EXISTS "stats" (
    id SERIAL PRIMARY KEY,
    character_id INTEGER NOT NULL,

    -- Base Stats
    health INTEGER NOT NULL,
    mana INTEGER NOT NULL,
    damage INTEGER NOT NULL,
    health_regeneration INTEGER NOT NULL,
    mana_regeneration INTEGER NOT NULL,

    -- Modifier Stats
    stamina INTEGER NOT NULL, -- Increases max health
    endurance INTEGER NOT NULL, -- Increases health regeneration
    intelligence INTEGER NOT NULL, -- Increases max mana
    spirit INTEGER NOT NULL, -- Increases mana regeneration
    magic_resistance INTEGER NOT NULL, -- Reduces magical damage taken
    spell_power INTEGER NOT NULL -- Increases magical damage

    CONSTRAINT fk_character FOREIGN KEY (character_id) REFERENCES characters(id) ON DELETE RESTRICT,
);
