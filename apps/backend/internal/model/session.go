package model

import (
	"time"
)

/*
CREATE TABLE IF NOT EXISTS user_sessions (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_hash TEXT NOT NULL,
    previous_refresh_hash TEXT,
    expires_at TIMESTAMPTZ,
    device_info TEXT,
    last_ip INET,
    revoked BOOLEAN DEFAULT false,
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
*/

type UserSession struct {
	ID                  int       `db:"id"`
	UserID              int       `db:"user_id"`
	RefreshHash         string    `db:"refresh_hash"`
	PreviousRefreshHash string    `db:"previous_refresh_hash"`
	ExpiresAt           time.Time `db:"expires_at"`
	// DeviceInfo          string    `db:"device_info"`
	// LastIP              net.IP    `db:"last_ip"`
	Revoked        bool      `db:"revoked"`
	LastActivityAt time.Time `db:"last_activity_at"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}
