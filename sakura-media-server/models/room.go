package models

import (
	"sync"
)

type Room struct {
	ID           string
	Participants map[string]*Participant
	Mutex        sync.RWMutex
}
