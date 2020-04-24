package json

import (
	"encoding/json"
)

func Serialize(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	return string(bytes), err
}

func Deserialize(data string, target interface{}) error {
	return json.Unmarshal([]byte(data), target)
}
