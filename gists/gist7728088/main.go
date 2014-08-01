package gist7728088

import (
	"encoding/json"
	"errors"
)

func ParseGistId(gistJsonResponse []byte) (gistId string, err error) {
	var gistJson struct {
		Id *string
	}
	switch err := json.Unmarshal(gistJsonResponse, &gistJson); {
	case err == nil && gistJson.Id != nil:
		return *gistJson.Id, nil
	case err == nil && gistJson.Id == nil:
		return "", errors.New("gist id field missing")
	default:
		return "", err
	}
}
