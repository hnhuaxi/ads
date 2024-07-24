package ads

import (
	"encoding/json"

	"github.com/stretchr/objx"
)

type Map = objx.Map

func ToMap(val interface{}) (Map, error) {
	b, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return objx.FromJSON(string(b))
}

// ToMapSlice
func ToMapSlice(val interface{}) ([]Map, error) {
	b, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	return objx.FromJSONSlice(string(b))
}
