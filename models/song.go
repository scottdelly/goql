package models

import (
	"time"
)

type Song struct {
	Model
	Duration time.Duration
}
