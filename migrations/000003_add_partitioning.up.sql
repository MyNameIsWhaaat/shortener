DROP TABLE IF EXISTS click_events CASCADE;

CREATE TABLE click_events (
    id BIGSERIAL,
    short_code VARCHAR(50) NOT NULL REFERENCES urls(short_code) ON DELETE CASCADE,
    user_agent TEXT,
    ip INET,
    referer TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

CREATE TABLE click_events_2026_02 PARTITION OF click_events
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');

CREATE TABLE click_events_2026_03 PARTITION OF click_events
    FOR VALUES FROM ('2026-03-01') TO ('2026-04-01');

CREATE INDEX idx_click_events_short_code ON click_events(short_code);
CREATE INDEX idx_click_events_created_at ON click_events(created_at);