package entities

import "time"

type Sms struct {
	Token     string    `json:"token" bson:"token,omitempty" validate:"required"`
	Event     string    `json:"event" bson:"event,omitempty" validate:"required"`
	Email     string    `json:"email" bson:"email,omitempty" validate:"required"`
	App       string    `json:"app" bson:"app" validate:"required"`
	Otp       string    `json:"otp" bson:"otp" validate:"required"`
	UserId    string    `json:"userId" bson:"userId,omitempty" validate:"required"`
	Message   string    `json:"message" bson:"message,omitempty" validate:"required"`
	Expiry    time.Time `json:"expiry" bson:"expiry,omitempty" validate:"required"`
	TimeStamp string    `json:"timeStamp" bson:"timeStamp,omitempty" validate:"required"`
}

type SmsForStorage struct {
	UserId    string    `json:"userId" bson:"userId"`
	App       string    `json:"app" bson:"app"`
	Otp       string    `json:"otp" bson:"otp"`
	Expiry    time.Time `json:"expiry" bson:"expiry" validate:"required"`
	TimeStamp string    `json:"timeStamp" bson:"timeStamp"`
}

type OtpRequest struct {
	Requester string `json:"requester" bson:"requester" validate:"required"`
	Recipient string `json:"recipient" bson:"recipient" validate:"required"`
	AppName   string `json:"appName" bson:"appName" validate:"required"`
	Event     string `json:"event" bson:"event" validate:"required"`
}
