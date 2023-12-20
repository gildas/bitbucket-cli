package common

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type Link struct {
	HREF url.URL `json:"href" mapstructure:"href"`
}

// MarshalJSON implements the json.Marshaler interface.
func (link Link) MarshalJSON() (data []byte, err error) {
	type surrogate Link

	data, err = json.Marshal(struct {
		surrogate
		HREF core.URL `json:"href"`
	}{
		surrogate: surrogate(link),
		HREF:      core.URL(link.HREF),
	})
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (link *Link) UnmarshalJSON(data []byte) (err error) {
	type surrogate Link

	var inner struct {
		surrogate
		HREF core.URL `json:"href"`
	}

	if err = json.Unmarshal(data, &inner); errors.Is(err, errors.JSONUnmarshalError) {
		return err
	} else if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	*link = Link(inner.surrogate)
	link.HREF = inner.HREF.AsURL()
	return
}
