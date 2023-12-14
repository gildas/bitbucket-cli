package common

import (
	"github.com/gildas/go-errors"
	"github.com/google/uuid"
)

type UUID uuid.UUID

func NewUUID() UUID {
	return UUID(uuid.New())
}

func ParseUUID(s string) (UUID, error) {
	u, err := uuid.Parse(s)
	return UUID(u), err
}

func (u UUID) IsNil() bool {
	return uuid.UUID(u) == uuid.Nil
}

func (u UUID) String() string {
	return "{" + uuid.UUID(u).String() + "}"
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.String() + `"`), nil
}

func (u *UUID) UnmarshalJSON(payload []byte) error {
	if len(payload) < 2 {
		return errors.JSONUnmarshalError.Wrap(errors.Errorf("unexpected end of JSON input"))
	}
	value := string(payload[1 : len(payload)-1])
	if len(value) == 0 {
		*u = UUID(uuid.Nil)
		return nil
	}
	parsed, err := ParseUUID(value)
	if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*u = parsed
	return nil
}
