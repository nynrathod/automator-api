package utilities

import (
	"crypto/rand"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nynrathod/automator-api/config"
	"github.com/nynrathod/automator-api/pkg/entities"
	"strings"
	"time"
)

func MsgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "Invalid email"
	}
	return fe.Error()
}

func GenerateRandomID(length int, userType int) string {
	var charset string
	if userType == 0 {
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	}
	if userType == 1 {
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	randomID := make([]byte, length)

	_, err := rand.Read(randomID)
	if err != nil {
		return ""
	}

	// Use the specified charset
	encoded := make([]byte, length)
	for i := range encoded {
		encoded[i] = charset[int(randomID[i])%len(charset)]
	}

	// Add timestamp for both user types
	timestamp := fmt.Sprintf("%v", time.Now().UnixNano())
	remainingSpace := length - len(encoded)
	if remainingSpace > 0 {
		encoded = append(encoded, timestamp[:remainingSpace]...)
	}

	return string(encoded)
}

func GenerateJWT(additionalClaims jwt.MapClaims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	// Add additional claims
	for key, value := range additionalClaims {
		if key != "exp" {
			claims[key] = value
		}
	}

	if exp, ok := additionalClaims["exp"]; ok {
		if expValue, ok := exp.(int); ok {
			expDuration := time.Duration(expValue) * time.Minute
			expNew := time.Now().Add(expDuration).Unix()
			claims["exp"] = expNew
		} else {
			return "", nil
		}
	} else {
		exp := time.Now().Add(time.Hour * 2400).Unix()
		claims["exp"] = exp
	}

	t, err := token.SignedString([]byte(config.EnvConfigs.JWTSecrete))
	if err != nil {
		return "", err
	}

	return t, nil
}

func VerifyToken(identity string, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("error extracting claims")
		}

		// Now you can access individual claims from the map
		identityClaim, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("error extracting email claim")
		}

		if identityClaim != identity {
			return nil, fmt.Errorf("email mismatch: %s (token) vs %s (provided)", identityClaim, identity)
		}

		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the key used for signing
		return []byte(config.EnvConfigs.JWTSecrete), nil
	})

	// Check for parsing errors
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, err

	// return aa, err
}

func intToPointer(i int) *int {
	return &i
}

func ValidateUser(user *entities.User, authHeader string) []entities.UserErrors {

	fmt.Println("user", user)

	if authHeader == "" {
		return append([]entities.UserErrors{}, entities.UserErrors{
			Param:   "authorization",
			Message: "Authorization header is required",
		})
	}

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {

		fmt.Println("regerr", err)

		var apiErrors []entities.UserErrors
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			for _, fe := range validationErrors {
				//// Skip the field "verifyToken" in the loop
				//if fe.Field() == "VerifyToken" {
				//	continue
				//}
				apiErrors = append(apiErrors, entities.UserErrors{Param: strings.ToLower(fe.Field()), Message: MsgForTag(fe)})
			}
		}
		return apiErrors
	}

	//_, tokenErr := VerifyToken(user.Email, authHeader)
	//if tokenErr != nil {
	//	fmt.Println("login err", tokenErr)
	//
	//	return append([]entities.UserErrors{}, entities.UserErrors{
	//		Param:   "authorization",
	//		Message: "Authorization header is wrong",
	//		Status:  RegisterInvalidRequest,
	//	})
	//
	//}

	return nil // No validation errors
}
