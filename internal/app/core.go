package core

import (
	"time"
)

type SendRequest struct {
	Sender   string
	Receptor string
	Message  string
	SentAt   time.Time
}
