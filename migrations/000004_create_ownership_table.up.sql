CREATE TABLE IF NOT EXISTS ownership (
    user_id INTEGER REFERENCES users (id),
    record_id INTEGER REFERENCES records (id),
    UNIQUE (user_id, record_id)
)
