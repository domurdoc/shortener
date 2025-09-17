CREATE TABLE IF NOT EXISTS records (
    id SERIAL PRIMARY KEY,
    "key" VARCHAR(6) NOT NULL,
    "value" VARCHAR(2048) NOT NULL
);

ALTER TABLE
    records
ADD
    CONSTRAINT unique_records_key UNIQUE (key);
