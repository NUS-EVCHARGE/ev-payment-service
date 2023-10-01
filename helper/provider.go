package helper

import (
	"bytes"
	"encoding/json"
	"github.com/NUS-EVCHARGE/ev-provider-service/dto"
	"io"
	"net/http"
)

func GetProvider(getProviderURL string, jwtToken string) (dto.Provider, error) {

	var (
		provider   dto.Provider
		httpClient = http.Client{}
	)

	url := getProviderURL
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return provider, err
	}

	req.Header.Add("Authentication", jwtToken)
	respReader, err := httpClient.Do(req)
	if err != nil {
		return provider, err
	}

	respByte, err := io.ReadAll(respReader.Body)
	if err != nil {
		return provider, err
	}

	err = json.Unmarshal(respByte, &provider)
	if err != nil {
		return provider, err
	}

	return provider, nil
}
