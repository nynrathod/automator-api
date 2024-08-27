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
	StoreSms(smsData *entities.Sms) (*entities.Sms, error)
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
	// List of keywords to match exactly, as whole words
	keywords := []string{
		"shemaroome",
		"sony liv",
		"disney+ hotstar",
		"jiocinema",
		"epic on",
		"hoichoi",
		"discovery plus",
		"etv win",
		"zee5",
	}

	// Convert message to lowercase for case-insensitive matching
	message = strings.ToLower(message)

	// Define a regular expression pattern to match 4 to 5 digits
	otpPattern := regexp.MustCompile(`\b\d{4,6}\b`)

	// Check for keywords
	for _, keyword := range keywords {
		if strings.Contains(message, keyword) {
			// Extract OTP from the message
			otp := otpPattern.FindString(message)
			return keyword, otp
		}
	}
	return "", ""
}

func (s *service) StoreSms(smsData *entities.Sms) (*entities.Sms, error) {
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
		return s.repository.StoreSms(smsData)
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
