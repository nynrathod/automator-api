package presenter

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nynrathod/automator-api/pkg/entities"
	UTL "github.com/nynrathod/automator-api/utilities"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//UTL "github.com/nynrathod/automator-api/utilities"
	"net/http"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	//CreatedAt time.Time          `json:"created_at"`
	SecretKey string `json:"-"`
}

func UserRegisterResponse(data *entities.User) *fiber.Map {
	return &fiber.Map{
		"id":        data.ID,
		"email":     data.Email,
		"firstName": data.FirstName,
		"lastName":  data.LastName,
		//"privateKey": data.PrivateKey,
		//"iv":         data.CrypInitializationVector,
		//"tag":        data.CrypTag,
		"status":     true,
		"statusCode": http.StatusOK,
	}
}

func UserProfileResponse(data *entities.User) *fiber.Map {
	fmt.Println("pr", data)
	if data == nil {
		// User not found, return an error response
		return &fiber.Map{
			"status":     false,
			"error":      "User not found",
			"statusCode": http.StatusNotFound,
		}
	}
	user := fiber.Map{
		"ID":        data.ID,
		"email":     data.Email,
		"firstName": data.FirstName,
		"lastName":  data.LastName,
		//"token":     "",
		//"privateKey": data.PrivateKey,
		//"iv":         data.CrypInitializationVector,
		//"tag":        data.CrypTag,
	}
	return &fiber.Map{
		"data": user,
	}
}

func UserRegisterErrResponse(err error) *fiber.Map {
	return &fiber.Map{
		"error": err.Error(),
	}
}

func VerifyEmailSuccess(email string) interface{} {
	return fiber.Map{
		"status": true,
		"email":  email,
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
	return fiber.Map{
		"status":  true,
		"message": "Password is valid",
		//"statusCode": http.StatusOK,
		"userId":    data.UserId,
		"email":     data.Email,
		"firstName": data.FirstName,
		"lastName":  data.LastName,
		//"privateKey": data.PrivateKey,
		//"iv":         data.CrypInitializationVector,
		//"tag":        data.CrypTag,
	}

}

func LoginError(data *entities.User) fiber.Map {
	return fiber.Map{
		"status":     false,
		"message":    "Invalid password",
		"statusCode": http.StatusUnauthorized,
	}
}

func OtpVerificationSuccess(email string) fiber.Map {
	additionalClaims := jwt.MapClaims{
		"email": email,
		"exp":   1,
	}
	jwtToken, _ := UTL.GenerateJWT(additionalClaims)
	return fiber.Map{
		"status":  true,
		"message": "OTP is valid",
		"token":   jwtToken,
	}
}

func OtpVerificationError(email string) fiber.Map {
	return fiber.Map{
		"status":  false,
		"message": "Invalid OTP",
	}
}
