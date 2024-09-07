package presenter

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nynrathod/automator-api/pkg/entities"
	UTL "github.com/nynrathod/automator-api/utilities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Email        string             `json:"email" bson:"email,omitempty"`
	Password     string             `json:"password" bson:"password,omitempty"`
	MobileNumber string             `json:"mobileNumber"`
	FirstName    string             `json:"firstName"`
	LastName     string             `json:"lastName"`
	//CreatedAt time.Time          `json:"created_at"`
	SecretKey string `json:"-"`
}

func UserRegisterResponse(data *entities.User) *fiber.Map {

	additionalClaims := jwt.MapClaims{
		"email": data.Email,
		"exp":   1,
	}
	jwtToken, _ := UTL.GenerateJWT(additionalClaims)

	return &fiber.Map{
		"status": UTL.RegisterSuccess,
		"token":  jwtToken,
	}
}

func UserProfileResponse(result *entities.User) *fiber.Map {
	//Expiry of 1 year
	additionalClaims := jwt.MapClaims{
		"exp":          525600,
		"id":           result.ID,
		"mobileNumber": result.MobileNumber,
		"firstName":    result.FirstName,
		"lastName":     result.LastName,
		"email":        result.Email,
	}

	jwtToken, _ := UTL.GenerateJWT(additionalClaims)

	return &fiber.Map{
		"status": UTL.ProfileSuccess,
		"data": fiber.Map{
			"id":           result.ID,
			"mobileNumber": result.MobileNumber,
			"firstName":    result.FirstName,
			"lastName":     result.LastName,
			"email":        result.Email,
			"token":        jwtToken,
		},
	}
}

func UserRegisterErrResponse(err error) *fiber.Map {
	var writeException mongo.WriteException
	if errors.As(err, &writeException) {
		for _, writeError := range writeException.WriteErrors {
			if writeError.Code == 11000 {
				return &fiber.Map{
					"status": UTL.RegisterExists,
					"error":  "Email already exists",
				}
			}
		}
	}

	return &fiber.Map{
		"status": UTL.RegisterSuccess,
	}
}

func VerifyEmailSuccess(email string) interface{} {

	additionalClaims := jwt.MapClaims{
		"email": email,
		"exp":   1,
	}
	jwtToken, _ := UTL.GenerateJWT(additionalClaims)

	return fiber.Map{
		"status": true,
		"email":  email,
		"token":  jwtToken,
	}
}

func VerifyEmailError(v interface{}) interface{} {
	if v == mongo.ErrNoDocuments {
		return fiber.Map{
			"status":     false,
			"statusCode": http.StatusNotFound,
		}
	} else {
		return fiber.Map{
			"status":     false,
			"statusCode": http.StatusInternalServerError,
		}
	}
}

func VerifyTokenResponse(isValid bool) fiber.Map {

	if isValid {
		return fiber.Map{
			"status":     true,
			"statusCode": http.StatusOK,
		}
	} else {
		return fiber.Map{
			"status":     false,
			"message":    "Invalid request",
			"statusCode": http.StatusUnauthorized,
		}
	}
}

func LoginSuccess(data *entities.User) fiber.Map {
	additionalClaims := jwt.MapClaims{
		"email": data.Email,
		"exp":   1,
	}
	jwtToken, _ := UTL.GenerateJWT(additionalClaims)

	return fiber.Map{
		"status": UTL.LoginSuccess,
		"token":  jwtToken,
	}

}

func LoginError() fiber.Map {
	return fiber.Map{
		"status":  UTL.LoginWrongCredentials,
		"message": "Wrong email or password",
	}
}

func OtpVerificationSuccess(email, verifyType string) fiber.Map {
	exp := 3
	status := UTL.OTPSignupSuccess

	if verifyType == "login" {
		exp = 1
		status = UTL.OTPLoginSuccess
	}

	additionalClaims := jwt.MapClaims{
		"email": email,
		"exp":   exp,
	}

	jwtToken, _ := UTL.GenerateJWT(additionalClaims)

	return fiber.Map{
		"status": status,
		"token":  jwtToken,
	}
}

func OtpVerificationError(err int) fiber.Map {
	if err == http.StatusUnauthorized {
		return fiber.Map{
			"status": UTL.OTPInvalidOTP,
		}
	}
	return fiber.Map{
		"status": UTL.OTPInvalidrequest,
	}
}

func OtpSendError(err int) fiber.Map {
	if err == http.StatusUnauthorized {
		return fiber.Map{
			"status": UTL.OTPInvalidOTP,
		}
	}
	return fiber.Map{
		"status": UTL.OTPInvalidrequest,
	}
}

func OtpSendResponse(email string) fiber.Map {
	additionalClaims := jwt.MapClaims{
		"email": email,
		"exp":   1,
	}

	jwtToken, _ := UTL.GenerateJWT(additionalClaims)
	return fiber.Map{
		"status": true,
		"token":  jwtToken,
	}
}

func GetProfileError(err int) fiber.Map {
	if err == http.StatusUnauthorized {
		return fiber.Map{
			"status": UTL.ProfileInvalidRequest,
		}
	}
	return fiber.Map{
		"status": UTL.ProfileNotFound,
	}
}
