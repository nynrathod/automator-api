package sms

import (
	"fmt"
	"github.com/nynrathod/automator-api/pkg/entities"
	"golang.org/x/net/context"
	"regexp"
	"strings"
	"time"
)

type Service interface {
	StoreSms(smsData *entities.Sms, uuid string) (*entities.Sms, error)
	ProcessOtpRequest(requestData *entities.OtpRequest, uuid context.Context) (*entities.OtpRequest, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func analyzeSms(message string) (string, string) {
	// Map of keywords to their corresponding values
	keywordMappings := map[string]string{
		"shemaroome":      "shemarooMe",
		"sony liv":        "sonyLiv",
		"disney+ hotstar": "disneyPlusHotstar",
		"jiocinema":       "jioCinema",
		"epic on":         "epicOn",
		"hoichoi":         "hoichoi",
		"discovery plus":  "discoveryPlus",
		"etv win":         "etvWin",
		"zee5":            "zee5",
	}

	message = strings.ToLower(message)
	otpPattern := regexp.MustCompile(`\b\d{4,6}\b`)

	for keyword, mappedValue := range keywordMappings {
		if strings.Contains(message, keyword) {
			otp := otpPattern.FindString(message)
			return mappedValue, otp
		}
	}
	return "", ""
}

func (s *service) StoreSms(smsData *entities.Sms, uuid string) (*entities.Sms, error) {
	now := time.Now()
	expiryTime := now.Add(4 * time.Minute)
	smsData.Expiry = expiryTime

	matchedKeyword, otp := analyzeSms(smsData.Message)
	if matchedKeyword != "" {
		smsData.App = matchedKeyword
		if otp != "" {
			smsData.Otp = otp
		}
	}

	if smsData.App != "" && smsData.Otp != "" {
		//fmt.Println("smsApp", smsData.App)
		//fmt.Println("smsOtp", smsData.Otp)
		return s.repository.StoreSms(smsData, uuid)
	}

	return nil, fmt.Errorf("failed to store SMS: missing keyword or OTP")
}

func (s *service) ProcessOtpRequest(requestData *entities.OtpRequest, uuid context.Context) (*entities.OtpRequest, error) {
	_, _ = s.repository.ProcessOtpRequest(requestData, uuid)
	return nil, nil
}
