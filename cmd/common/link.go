package common

import (
	"encoding/json"
	"net/url"

	"github.com/gildas/go-core"
	"github.com/gildas/go-errors"
)

type Link struct {
	Name   string  `json:"name,omitempty" mapstructure:"name"`
	HREF   url.URL `json:"-"`
	GitRef string  `json:"-"`
}

// MarshalJSON implements the json.Marshaler interface.
func (link Link) MarshalJSON() (data []byte, err error) {
	type surrogate Link

	if len(link.GitRef) > 0 {
		data, err = json.Marshal(struct {
			surrogate
			GitRef string `json:"href"`
		}{
			surrogate: surrogate(link),
			GitRef:    link.GitRef,
		})
	} else {
		data, err = json.Marshal(struct {
			surrogate
			HREF core.URL `json:"href"`
		}{
			surrogate: surrogate(link),
			HREF:      core.URL(link.HREF),
		})
	}
	return data, errors.JSONMarshalError.Wrap(err)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (link *Link) UnmarshalJSON(data []byte) (err error) {
	type surrogate Link

	var header struct {
		Name string `json:"name"`
	}
	if err = json.Unmarshal(data, &header); errors.Is(err, errors.JSONUnmarshalError) {
		return err
	} else if err != nil {
		return errors.JSONUnmarshalError.Wrap(err)
	}
	switch header.Name {
	case "ssh":
		var inner struct {
			surrogate
			GitRef string `json:"href"`
		}
		if err = json.Unmarshal(data, &inner); errors.Is(err, errors.JSONUnmarshalError) {
			return err
		} else if err != nil {
			return errors.JSONUnmarshalError.Wrap(err)
		}
		*link = Link(inner.surrogate)
		link.GitRef = inner.GitRef
	default:
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
	}
	return
}
