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

		_, err := service.VerifyEmail(requestBody.Email)
		if err != nil {
			response := presenter.VerifyEmailError(err)
			c.Status(response.(fiber.Map)["statusCode"].(int))

			return c.JSON(fiber.Map{"status": response.(fiber.Map)["status"], "error": err.Error()})
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
		var userData *entities.User
		userModel, err := new(entities.User), *new(error)

		if isEmail(identity) {
			userModel, err = service.VerifyEmail(identity)
		} else {
			//userModel, err = service.VerifyUserName(identity)
		}

		if userModel == nil {
			//fmt.Println("userdata", userModel.Password)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": false, "message": "User not found",
				"data": err.Error()})
		} else {

			userData = &entities.User{
				UserId:   userModel.UserId,
				Email:    userModel.Email,
				Password: userModel.Password,
				//CrypInitializationVector: userModel.CrypInitializationVector,
				FirstName: userModel.FirstName,
				LastName:  userModel.LastName,
				//PrivateKey:               userModel.PrivateKey,
				//CrypTag:                  userModel.CrypTag,
			}
		}
		fmt.Println("userdata else", userData.Password)
		if !CheckPasswordHash(pass, userData.Password) {
			response := presenter.LoginError(userData)
			statusCode := response["statusCode"].(int)
			delete(response, "statusCode")
			return c.Status(statusCode).JSON(response)
		}

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		response := presenter.LoginSuccess(userData)

		return c.JSON(response)
	}
}

func Register(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.User
		err := c.BodyParser(&requestBody)

		authHeader := c.Get("Authorization")

		// fmt.Println("reqbody: ", requestBody.VerifyToken)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenter.UserRegisterErrResponse(err))
		}

		existingUser, err := service.VerifyEmail(requestBody.Email)
		if err == nil && existingUser != nil {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"error": "Email is already in use",
			})
		}

		apiErrors := UTL.ValidateUser(&requestBody, authHeader)
		if apiErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": apiErrors,
			})
		}

		result, err := service.Register(&requestBody)

		if err != nil {
			c.Status(http.StatusBadRequest)
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

		go func() {
			time.Sleep(5 * time.Minute)
			otpMapMutex.Lock()
			delete(otpMap, requestBody.Email)
			otpMapMutex.Unlock()
		}()

		return c.JSON("sa")
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

		if ok && storedOTP == submittedOtp {
			otpMapMutex.Lock()
			delete(otpMap, requestBody.Email)
			otpMapMutex.Unlock()
			response := presenter.OtpVerificationSuccess(requestBody.Email)
			if requestBody.VerifyType == "login" {
				//delete(response, "token")
				response["isLogin"] = true
			}
			return c.JSON(response)
		} else {
			response := presenter.OtpVerificationError(requestBody.Email)
			return c.Status(http.StatusUnauthorized).JSON(response)
		}

	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte("From: Team Uoozer\r\n" + "To: " + strings.Join(r.to, ", ") + "\r\n" + subject + mime + "\n" + r.body)

	addr := "smtp.gmail.com:587"

	auth = smtp.PlainAuth("", config.EnvConfigs.SmtpEmail, config.EnvConfigs.SmtpToken, "smtp.gmail.com")

	if err := smtp.SendMail(addr, auth, config.EnvConfigs.SmtpEmail, r.to, msg); err != nil {
		return false, err
	}
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
