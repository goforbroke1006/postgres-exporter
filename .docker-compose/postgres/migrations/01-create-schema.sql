CREATE TABLE IF NOT EXISTS quotes
(
    id         SERIAL PRIMARY KEY,
    bid        DECIMAL(18, 6) DEFAULT 0.0   NOT NULL,
    ask        DECIMAL(18, 6) DEFAULT 0.0   NOT NULL,
    updated_at TIMESTAMP      DEFAULT NOW() NOT NULL
);
