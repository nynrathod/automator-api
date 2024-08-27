package users

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/nynrathod/automator-api/api/presenter"
	"github.com/nynrathod/automator-api/config"
	"github.com/nynrathod/automator-api/pkg/entities"
	"github.com/nynrathod/automator-api/pkg/users"
	UTL "github.com/nynrathod/automator-api/utilities"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"sync"
	"text/template"
	"time"
)

type Request struct {
	to      []string
	subject string
	body    string
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

var auth smtp.Auth
var otpMap = make(map[string]string)
var otpMapMutex sync.Mutex

const otpChars = "0123456789"

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func VerifyEmail(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.VerifyEmail
		_ = c.BodyParser(&requestBody)

		_, existsErr := service.VerifyEmail(requestBody.Email, false)
		fmt.Println("existsErr", existsErr)
		if existsErr == nil {
			fmt.Println("iiierr", existsErr)
			c.Status(http.StatusConflict)
			return c.JSON(presenter.UserRegisterErrResponse(mongo.WriteException{
				WriteErrors: []mongo.WriteError{
					{Code: 11000, Message: "Duplicate key error"},
				},
			}))
		}

		response := presenter.VerifyEmailSuccess(requestBody.Email)
		return c.JSON(response)
	}
}

func Login(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		input := new(entities.Login)

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error",
				"message": "Error on login request", "data": err.Error()})
		}

		pass := input.Password
		identity := input.Email

		if isEmail(identity) {
			userModel, userErr := service.VerifyEmail(identity, true)
			fmt.Println("userErr", userErr)
			if userErr != nil {

				return c.Status(http.StatusUnauthorized).JSON(presenter.LoginError())
			}
			var userData *entities.User
			userData = &entities.User{Password: userModel.Password, Email: userModel.Email}
			fmt.Println("userdata else", userData.Password)

			if !CheckPasswordHash(pass, userData.Password) {
				response := presenter.LoginError()
				return c.Status(http.StatusUnauthorized).JSON(response)
			}

			response := presenter.LoginSuccess(userData)

			return c.JSON(response)

		}

		return nil
	}
}

func Register(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.User
		err := c.BodyParser(&requestBody)

		authHeader := c.Get("Authorization")

		// fmt.Println("reqbody: ", requestBody.VerifyToken)
		if err != nil {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "Error parsing",
			})
		}

		tokenError := UTL.ValidateUser(&requestBody, authHeader)
		if tokenError != nil {
			return c.Status(fiber.StatusBadRequest).JSON(tokenError)
		}

		_, existsErr := service.VerifyEmail(requestBody.Email, false)
		fmt.Println("existsErr", existsErr)
		if existsErr == nil {
			fmt.Println("iiierr", existsErr)
			c.Status(http.StatusConflict)
			return c.JSON(presenter.UserRegisterErrResponse(mongo.WriteException{
				WriteErrors: []mongo.WriteError{
					{Code: 11000, Message: "Duplicate key error"},
				},
			}))
		}

		result, err := service.Register(&requestBody)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.UserRegisterErrResponse(err))
		}

		return c.JSON(presenter.UserRegisterResponse(result))
	}
}

func GenerateOTP(length int) (string, error) {
	otpCharsLength := big.NewInt(int64(len(otpChars)))
	buffer := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, otpCharsLength)
		if err != nil {
			return "", err
		}
		buffer[i] = otpChars[randomIndex.Int64()]
	}
	return string(buffer), nil
}

func SendOtp(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.User
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.UserRegisterErrResponse(err))
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			print("anythingempty")
			return c.Status(http.StatusForbidden).JSON(presenter.OtpSendError(http.StatusForbidden))
		}

		_, tokenErr := UTL.VerifyToken(requestBody.Email, authHeader)
		if tokenErr != nil {
			fmt.Println("incorrec")
			return c.Status(http.StatusForbidden).JSON(presenter.OtpSendError(http.StatusForbidden))
		}

		otp, _ := GenerateOTP(6)

		// Store the OTP in the map with the user's email
		otpMapMutex.Lock()
		otpMap[requestBody.Email] = otp
		otpMapMutex.Unlock()

		templateData := struct {
			Otp string
			URL string
		}{
			Otp: otp,
			URL: "http://yuezers.app",
		}
		r := NewRequest([]string{requestBody.Email}, "OTP verification from yuezers", "")
		r.ParseTemplate("api/handlers/template.html", templateData)
		if err := r.ParseTemplate("api/handlers/template.html", templateData); err == nil {
			_, _ = r.SendEmail()
		}
		fmt.Println("coming here")
		go func() {
			time.Sleep(5 * time.Minute)
			otpMapMutex.Lock()
			delete(otpMap, requestBody.Email)
			otpMapMutex.Unlock()
		}()

		return c.JSON(presenter.OtpSendResponse(requestBody.Email))
	}
}

func VerifyOtp(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.LoginEmailOtp
		err := c.BodyParser(&requestBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(presenter.UserRegisterErrResponse(err))
		}

		submittedOtp := requestBody.Otp
		storedOTP, ok := otpMap[requestBody.Email]

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			print("anythingempty")
			return c.Status(http.StatusForbidden).JSON(presenter.OtpVerificationError(http.StatusForbidden))
		}

		_, tokenErr := UTL.VerifyToken(requestBody.Email, authHeader)
		if tokenErr != nil {
			fmt.Println("incorrec")
			return c.Status(http.StatusForbidden).JSON(presenter.OtpVerificationError(http.StatusForbidden))
		}

		if ok && storedOTP == submittedOtp {
			otpMapMutex.Lock()
			delete(otpMap, requestBody.Email)
			otpMapMutex.Unlock()
			response := presenter.OtpVerificationSuccess(requestBody.Email, requestBody.VerifyType)
			if requestBody.VerifyType == "login" {
				response["isLogin"] = true
			}
			return c.JSON(response)
		} else {
			response := presenter.OtpVerificationError(http.StatusUnauthorized)
			return c.Status(http.StatusUnauthorized).JSON(response)
		}

	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte("From: Team Uoozer\r\n" + "To: " + strings.Join(r.to, ", ") + "\r\n" + subject + mime + "\n" + r.body)

	addr := "smtp.gmail.com:587"
	fmt.Println("onmail", config.EnvConfigs.SmtpEmail, config.EnvConfigs.SmtpToken)
	auth = smtp.PlainAuth("", config.EnvConfigs.SmtpEmail, config.EnvConfigs.SmtpToken, "smtp.gmail.com")

	if err := smtp.SendMail(addr, auth, config.EnvConfigs.SmtpEmail, r.to, msg); err != nil {
		fmt.Println("mailerrr", err)
		return false, err
	}
	fmt.Println("mail sent")
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}
