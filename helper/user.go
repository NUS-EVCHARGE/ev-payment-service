package helper

import (
	"bytes"
	"encoding/json"
	dto2 "ev-payment-service/dto"
	"fmt"
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

func Getbooking(getBookingUrl string, jwtToken string, bookingId uint) ([]dto2.Booking, error) {
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
