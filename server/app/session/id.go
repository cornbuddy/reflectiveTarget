package session

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type SessionID = uuid.UUID

var errBadVersion = errors.New("bad version of UUID")

const wantUUIDVersion = 7

func IDFromString(id string) (*SessionID, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ver := uuid.Version()
	if ver != wantUUIDVersion {
		return nil, fmt.Errorf("%w: %d", errBadVersion, ver)
	}

	return &uuid, nil
}

func MakeID() SessionID {
	// I'm concerned about panic here
	return uuid.Must(uuid.NewV7())
}
