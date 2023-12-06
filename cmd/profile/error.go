package profile

import (
	"encoding/json"
	"strings"

	"github.com/gildas/go-errors"
)

type BitBucketError struct {
	Type    string              `json:"type"`
	Message string              `json:"-"`
	Detail  map[string][]string `json:"fields"`
}

func (bberr *BitBucketError) Error() string {
	var buffer strings.Builder

	buffer.WriteString(bberr.Message)
	if len(bberr.Detail) > 0 {
		buffer.WriteString(" (")
		for field, messages := range bberr.Detail {
			buffer.WriteString(field)
			buffer.WriteString(": ")
			buffer.WriteString(strings.Join(messages, ", "))
			buffer.WriteString("; ")
		}
		buffer.WriteString(")")
	}
	return buffer.String()
}

// UnmarshalJSON unmarshals the JSON
func (bberr *BitBucketError) UnmarshalJSON(data []byte) (err error) {
	type surrogate BitBucketError
	var inner struct {
		surrogate
		Error struct {
			Message string              `json:"message"`
			Detail  map[string][]string `json:"detail"`
			Fields  map[string][]string `json:"fields"`
		} `json:"error"`
	}
	if err = json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*bberr = BitBucketError(inner.surrogate)
	bberr.Message = inner.Error.Message
	bberr.Detail = inner.Error.Detail
	return
}
