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

	// Convert message to lowercase for case-insensitive matching
	message = strings.ToLower(message)

	// Define a regular expression pattern to match 4 to 5 digits
	otpPattern := regexp.MustCompile(`\b\d{4,6}\b`)

	// Check for keywords
	for keyword, mappedValue := range keywordMappings {
		if strings.Contains(message, keyword) {
			// Extract OTP from the message
			otp := otpPattern.FindString(message)
			return mappedValue, otp // Return the mapped value and the OTP
		}
	}
	return "", ""
}

func (s *service) StoreSms(smsData *entities.Sms, uuid string) (*entities.Sms, error) {
	now := time.Now()
	expiryTime := now.Add(4 * time.Minute)
	smsData.Expiry = expiryTime

	// Check if the SMS message contains any of the keywords
	matchedKeyword, otp := analyzeSms(smsData.Message)
	if matchedKeyword != "" {
		smsData.App = matchedKeyword // Store the matched keyword in smsData.App
		if otp != "" {
			smsData.Otp = otp // Store the extracted OTP
		}
	}
	// Check if both App and Otp are set before storing the SMS data
	if smsData.App != "" && smsData.Otp != "" {
		fmt.Println("smsApp", smsData.App)
		fmt.Println("smsOtp", smsData.Otp)
		return s.repository.StoreSms(smsData, uuid)
	}

	// If either App or Otp is empty, return an error or handle it appropriately
	return nil, fmt.Errorf("failed to store SMS: missing keyword or OTP")
}

func (s *service) ProcessOtpRequest(requestData *entities.OtpRequest, uuid context.Context) (*entities.OtpRequest, error) {
	fmt.Println("requestData", requestData)
	//filter := bson.D{{"mobile_number", requestData.Recipient}}
	//err := collection.FindOne(context.Background(), filter).Decode(&result)
	_, _ = s.repository.ProcessOtpRequest(requestData, uuid)
	return nil, nil
}
