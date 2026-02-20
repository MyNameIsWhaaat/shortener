DROP TABLE IF EXISTS click_events CASCADE;

CREATE TABLE click_events (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(50) NOT NULL REFERENCES urls(short_code) ON DELETE CASCADE,
    user_agent TEXT,
    ip INET,
    referer TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_click_events_short_code ON click_events(short_code);
CREATE INDEX idx_click_events_created_at ON click_events(created_at);