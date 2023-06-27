package models

import (
	"encoding/json"
	"time"
)

type OTP struct {
	Namespace   string          `redis:"namespace" json:"namespace" db:"namespace"`
	ID          string          `redis:"id" json:"id" db:"id"`
	To          string          `redis:"to" json:"to" db:"to"`
	ChannelDesc string          `redis:"channel_description" json:"channel_description" db:"channel_desc"`
	AddressDesc string          `redis:"address_description" json:"address_description" db:"address_desc"`
	Extra       json.RawMessage `redis:"extra" json:"extra" db:"extra"`
	Provider    string          `redis:"provider" json:"provider" db:"provider"`
	OTP         string          `redis:"otp" json:"otp" db:"otp"`
	MaxAttempts int             `redis:"max_attempts" json:"max_attempts" db:"max_attempts"`
	Attempts    int             `redis:"attempts" json:"attempts" db:"attempts"`
	Closed      bool            `redis:"closed" json:"closed" db:"closed"`
	TTL         time.Duration   `redis:"-" json:"-" db:"ttl"`
	TTLSeconds  float64         `redis:"-" json:"ttl" db:"ttl_seconds"`
}
