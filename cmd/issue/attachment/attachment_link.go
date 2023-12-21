package attachment

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type AttachmentLink url.URL

// String returns the string representation of the URL
//
// implements the Stringer interface
func (link AttachmentLink) String() string {
	return (*url.URL)(&link).String()
}

// MarshalJSON marshals the URL into JSON
//
// implements the json.Marshaler interface
func (link AttachmentLink) MarshalJSON() (data []byte, err error) {
	type Self struct {
		HREF []core.URL `json:"href"`
	}
	type Link struct {
		Self Self `json:"self"`
	}
	data, err = json.Marshal(
		Link{
			Self: Self{HREF: []core.URL{core.URL(link)}},
		},
	)
	return data, errors.JSONUnmarshalError.Wrap(err)
}

// UnmarshalJSON unmarshals the URL from JSON
//
// implements the json.Unmarshaler interface
func (link *AttachmentLink) UnmarshalJSON(data []byte) error {
	var inner struct {
		Self struct {
			HREF []core.URL `json:"href"`
		} `json:"self"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	if len(inner.Self.HREF) == 0 {
		return errors.JSONUnmarshalError.Wrap(errors.ArgumentMissing.With("self.href"))
	}
	*link = AttachmentLink(inner.Self.HREF[0])
	return nil
}
