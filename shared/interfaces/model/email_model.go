package model

import (
	"encoding/json"
)

type Email struct {
	User_id uint
	Subject string
	Content string
}

func (m *Email) ToJsonString() (string, error) {
	if m != nil {
		jsonData, err := json.Marshal(&m)
		if err != nil {
			return "", err
		}

		return string(jsonData), nil
	}

	return "", nil
}
