package profile

import (
	"encoding/json"
	"strings"

	"github.com/gildas/go-errors"
)

type BitBucketError struct {
	Type    string              `json:"type"`
	Message string              `json:"-"`
	Detail  string              `json:"-"`
	Fields  map[string][]string `json:"-"`
}

func (bberr *BitBucketError) Error() string {
	var buffer strings.Builder

	buffer.WriteString(bberr.Message)
	if len(bberr.Detail) > 0 {
		buffer.WriteString(": ")
		buffer.WriteString(bberr.Detail)
	}
	if len(bberr.Fields) > 0 {
		buffer.WriteString(" (")
		for field, messages := range bberr.Fields {
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

	var innerType1 struct {
		surrogate
		Error struct {
			Message string              `json:"message"`
			Fields  map[string][]string `json:"detail"`
		} `json:"error"`
	}
	if err = json.Unmarshal(data, &innerType1); err == nil && len(innerType1.Error.Fields) > 0 {
		*bberr = BitBucketError(innerType1.surrogate)
		bberr.Message = innerType1.Error.Message
		bberr.Fields = innerType1.Error.Fields
		return
	}

	var innerType2 struct {
		surrogate
		Error struct {
			Message string            `json:"message"`
			Detail  string            `json:"detail"`
			Fields  map[string]string `json:"fields"`
		} `json:"error"`
	}

	if err = json.Unmarshal(data, &innerType2); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*bberr = BitBucketError(innerType2.surrogate)
	bberr.Message = innerType2.Error.Message
	bberr.Detail = innerType2.Error.Detail
	if len(innerType2.Error.Fields) > 0 {
		bberr.Fields = make(map[string][]string)
		for field, message := range innerType2.Error.Fields {
			bberr.Fields[field] = []string{message}
		}
	}
	return
}
