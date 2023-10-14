package helper

import (
	"bytes"
	"encoding/json"
	dto2 "ev-payment-service/dto"
	"fmt"
	providerDTO "github.com/NUS-EVCHARGE/ev-provider-service/dto"
	"github.com/NUS-EVCHARGE/ev-user-service/dto"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
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

func GetBooking(getBookingUrl string, jwtToken string, bookingId uint) ([]dto2.Booking, error) {
	var (
		booking    []dto2.Booking
		httpClient = http.Client{}
		url        string
	)
	if bookingId == 0 {
		url = getBookingUrl
	} else {
		url = getBookingUrl + "/" + strconv.Itoa(int(bookingId))
	}

	logrus.WithField("url", url).Info("url")
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return booking, err
	}

	req.Header.Add("Authentication", jwtToken)
	logrus.WithField("req", req).Info("req")
	respReader, err := httpClient.Do(req)
	if err != nil {
		return booking, err
	}

	respByte, err := io.ReadAll(respReader.Body)
	if err != nil {
		return booking, err
	}

	err = json.Unmarshal(respByte, &booking)
	if err != nil {
		return booking, err
	}

	if len(booking) > 0 {
		if booking[0].Email == "" {
			var errGetBookingResp = map[string]interface{}{}
			err = json.Unmarshal(respByte, &errGetBookingResp)
			return booking, fmt.Errorf(errGetBookingResp["message"].(string))
		}
	}

	return booking, nil
}

func GetChargerById(getChargerUrl string, jwtToken string, chargerId uint) ([]providerDTO.Charger, error) {
	var (
		charger    []providerDTO.Charger
		httpClient = http.Client{}
		url        string
	)

	url = getChargerUrl + "/charger/" + strconv.Itoa(int(chargerId))

	logrus.WithField("url", url).Info("url")
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return charger, err
	}

	req.Header.Add("Authentication", jwtToken)
	logrus.WithField("req", req).Info("req")
	respReader, err := httpClient.Do(req)
	if err != nil {
		return charger, err
	}

	respByte, err := io.ReadAll(respReader.Body)
	if err != nil {
		return charger, err
	}

	err = json.Unmarshal(respByte, &charger)
	if err != nil {
		return charger, err
	}

	return charger, nil
}

func GetRateById(getRateURL string, jwtToken string, rateId uint) ([]providerDTO.Rates, error) {
	var (
		rate       []providerDTO.Rates
		httpClient = http.Client{}
		url        string
	)

	url = getRateURL + "/rates/" + strconv.Itoa(int(rateId))

	logrus.WithField("url", url).Info("url")
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return rate, err
	}

	req.Header.Add("Authentication", jwtToken)
	respReader, err := httpClient.Do(req)
	if err != nil {
		return rate, err
	}

	respByte, err := io.ReadAll(respReader.Body)
	if err != nil {
		return rate, err
	}

	err = json.Unmarshal(respByte, &rate)
	if err != nil {
		return rate, err
	}

	return rate, nil
}
