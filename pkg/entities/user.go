package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		ID        primitive.ObjectID `bson:"_id"`
		UserId    string             `json:"user_id" bson:"user_id, omitempty"`
		Email     string             `json:"email" bson:"email,omitempty" validate:"required,email"`
		UserName  string             `json:"userName" bson:"userName,omitempty"`
		Password  string             `json:"password" bson:"password,omitempty" validate:"required"`
		FirstName string             `json:"firstName" bson:"firstName,omitempty" validate:"required"`
		LastName  string             `json:"lastName" bson:"lastName,omitempty" validate:"required"`
		//SecretKey                string             `json:"secretKey" bson:"secretKey,omitempty"`
		//PrivateKey               string             `json:"privateKey" bson:"privateKey,omitempty" validate:"required"`
		//CrypInitializationVector string             `json:"cryptInitializationVector" bson:"cryptInitializationVector,omitempty" validate:"required"`
		//CrypTag                  string             `json:"crypTag" bson:"crypTag,omitempty"  validate:"required"`
		//VerifyToken string `json:"verifyToken" bson:"verifyToken,omitempty"  validate:"required"`
		//CreatedAt                time.Time          `json:"created_at" bson:"createdAt,omitempty" validate:"required"`
		UpdatedAt string `json:"updated_at" bson:"updatedAt,omitempty"`
	}

	UserErrors struct {
		Param   string
		Message string
	}

	Login struct {
		Email    string `json:"email" bson:"email,omitempty" validate:"required"`
		Password string `json:"password" bson:"password,omitempty" validate:"required"`
		Otp      string `json:"otp" bson:"otp,omitempty" validate:"required,len=6"`
		//VerifyToken string `json:"verifyToken" bson:"verifyToken,omitempty"  validate:"required"`
	}
	TokenVerify struct {
		Email         string `json:"email" bson:"email,omitempty"`
		Identity      string `json:"identity" bson:"identity,omitempty" validate:"required"`
		UserId        string `json:"user_id" bson:"user_id, omitempty"`
		TokenVerifier string `json:"tokenVerifier,omitempty" bson:"tokenVerifier,omitempty"`
	}

	LoginEmailPass struct {
		Email    string `json:"email" bson:"email,omitempty" validate:"required,email"`
		Password string `json:"password" bson:"password,omitempty" validate:"required"`
	}

	VerifyEmail struct {
		Email string `json:"email" bson:"email,omitempty" validate:"required,email"`
	}

	LoginEmailOtp struct {
		VerifyType string `json:"verifyType" bson:"verifyType,omitempty"`
		Email      string `json:"email" bson:"email,omitempty" validate:"required,email"`
		Otp        string `json:"otp" bson:"otp,omitempty" validate:"required,len=6"`
	}
)
