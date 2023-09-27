package custom_types

import (
	"encoding/json"
	"time"
)

type CustomDate time.Time

var _ json.Unmarshaler = &CustomDate{}

func (cd *CustomDate) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*cd = CustomDate(t)
	return nil
}
