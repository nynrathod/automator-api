package presenter

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nynrathod/automator-api/pkg/entities"
	UTL "github.com/nynrathod/automator-api/utilities"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//UTL "github.com/nynrathod/automator-api/utilities"
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

// var secretKey = config.EnvConfigs.JWTSecrete

// func verifyJWT(tokenString string) (*jwt.Token, error) {
// 	// Parse the token
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// Check the signing method
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		// Provide the key used for signing
// 		return []byte(config.EnvConfigs.JWTSecrete), nil
// 	})

// 	// Check for parsing errors
// 	if err != nil {
// 		return nil, fmt.Errorf("error parsing token: %v", err)
// 	}

// 	// Check if the token is valid
// 	if !token.Valid {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	return token, nil
// }

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

//	func OtpVerificationResponse(email string, isValid bool) fiber.Map {
//		// fmt.Println("isValid", email)
//
//		additionalClaims := jwt.MapClaims{
//			"identity": email,
//			"exp":      5,
//		}
//		jwtToken, _ := UTL.GenerateJWT(additionalClaims)
//		// fmt.Println("jwtToken", jwtToken)
//		// validatedToken, err := verifyJWT(jwtToken)
//		// if err != nil {
//		// 	fmt.Println("Token validation failed:", err)
//		// 	// Handle the error, e.g., return an error response
//		// 	return nil
//		// }
//
//		// // Token is valid, you can use validatedToken.Claims to access the claims
//		// fmt.Println("Token is valid. Claims:", validatedToken.Claims)
//		if isValid {
//			return fiber.Map{
//				"status":     true,
//				"message":    "OTP is valid",
//				"token":      jwtToken,
//				"statusCode": http.StatusOK,
//			}
//		} else {
//			return fiber.Map{
//				"status":     false,
//				"message":    "Invalid OTP",
//				"statusCode": http.StatusUnauthorized,
//			}
//		}
//	}
func LoginSuccess(data *entities.User) fiber.Map {
	fmt.Println("logingdata", data)

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

	//token := jwt.New(jwt.SigningMethodHS256)
	//
	//claims := token.Claims.(jwt.MapClaims)
	//claims["email"] = email
	//claims["admin"] = true
	//claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	//
	//t, err := token.SignedString([]byte("af4777d4f6c64492b3969ddc3da9301e"))
	//if err != nil {
	//
	//}
	//
	//return fiber.Map{
	//	"status": true,
	//	"token":  t,
	//}

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
