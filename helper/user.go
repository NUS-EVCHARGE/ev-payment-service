package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/NUS-EVCHARGE/ev-user-service/dto"
	"io"
	"net/http"
)

func GetUser(getUserUrl string, jwtToken string) (dto.User, error) {
	var (
		user       dto.User
		httpClient = http.Client{}
	)
	req, err := http.NewRequest("GET", getUserUrl, bytes.NewReader([]byte("")))
	if err != nil {
		return user, err
	}

	req.Header.Add("Authentication", jwtToken)

	respReader, err := httpClient.Do(req)
	if err != nil {
		return user, err
	}

	respByte, err := io.ReadAll(respReader.Body)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(respByte, &user)
	if err != nil {
		return user, err
	}
	if user.Email == "" {
		var errGetUserResp = map[string]interface{}{}
		err = json.Unmarshal(respByte, &errGetUserResp)
		return user, fmt.Errorf(errGetUserResp["message"].(string))
	}
	return user, nil
}
