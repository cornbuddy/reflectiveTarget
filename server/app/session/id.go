package session

import (
	"github.com/google/uuid"
)

type SessionID = uuid.UUID

func IDFromString(id string) (*SessionID, error) {
	return nil, nil
}

func MakeID() SessionID {
	// I'm concerned about panic here
	return uuid.Must(uuid.NewV7())
}
