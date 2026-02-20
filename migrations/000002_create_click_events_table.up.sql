CREATE TABLE IF NOT EXISTS click_events (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(50) NOT NULL REFERENCES urls(short_code) ON DELETE CASCADE,
    user_agent TEXT,
    ip INET,
    referer TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_click_events_short_code') THEN
        CREATE INDEX idx_click_events_short_code ON click_events(short_code);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_click_events_created_at') THEN
        CREATE INDEX idx_click_events_created_at ON click_events(created_at);
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_click_events_short_code_created_at') THEN
        CREATE INDEX idx_click_events_short_code_created_at ON click_events(short_code, created_at);
    END IF;
END $$;