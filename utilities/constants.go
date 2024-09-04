package utilities

const (

	// Login
	LoginSuccess          = "LOGIN-0001"
	LoginWrongCredentials = "LOGIN-0002"
	LoginInvalidEmail     = "LOGIN-0003"
	LoginInvalidPass      = "LOGIN-0004"
	LoginUserNotFound     = "LOGIN-0005"
	LoginInternalError    = "LOGIN-0005"

	// OTP verification
	OTPLoginSuccess   = "OTP-0001"
	OTPSignupSuccess  = "OTP-0002"
	OTPInvalidOTP     = "OTP-0003"
	OTPInvalidrequest = "OTP-0004"
	OTPInternalError  = "OTP-0005"

	// Register
	RegisterSuccess        = "REGISTER-0001"
	RegisterInvalidRequest = "REGISTER-0002"
	RegisterExists         = "REGISTER-0003"
	RegisterInternalError  = "REGISTER-0004"

	//	Profile
	ProfileSuccess        = "PROFILE-0001"
	ProfileInvalidRequest = "PROFILE-0002"
	ProfileNotFound       = "PROFILE-0003"
	ProfileInternalError  = "PROFILE-0004"
)
