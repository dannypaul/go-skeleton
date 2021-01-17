package exception

import (
	"errors"
	"net/http"
)

const (
	// User
	UserAlreadyExists               = "userAlreadyExists"
	UserNotFound                    = "userNotFound"
	UserNotRegistered               = "userNotRegistered"
	CredentialsInvalid              = "credentialsInvalid"
	FailedLoginLimitExceeded        = "failedLoginLimitExceeded"
	ChallengeNotFound               = "challengeNotFound"
	FailedVerificationLimitExceeded = "failedVerificationLimitExceeded"
	UserVerificationIncomplete      = "userVerificationIncomplete"
	VerificationFailed              = "verificationFailed"
	TooManyChallengeRequests        = "tooManyChallengeRequests"
	EmailIdInvalid                  = "emailIdInvalid"
	PhoneNumberInvalid              = "phoneNumberInvalid"

	// Cron
	MinuteIsInvalid    = "minuteIsInvalid"
	HourIsInvalid      = "hourIsInvalid"
	DayOfWeekIsInvalid = "dayOfWeekIsInvalid"

	// Identity
	IdentityTypeNotFound = "identityTypeNotFound"

	// Generic
	Unauthorised        = "unauthorised"
	Forbidden           = "forbidden"
	NotFound            = "notFound"
	InternalServerError = "internalServerError"
	Conflict            = "conflict"
	IdInvalid           = "idInvalid"
)

var ErrNotFound = errors.New(NotFound)
var ErrConflict = errors.New(Conflict)
var ErrIdInvalid = errors.New(IdInvalid)

var messages = map[string]string{
	// User
	UserAlreadyExists:               "User with the given phone number already exists",
	UserNotFound:                    "User not found",
	UserNotRegistered:               "User not registered",
	CredentialsInvalid:              "You have entered an invalid username or password",
	FailedLoginLimitExceeded:        "Exceeded failed login limit",
	ChallengeNotFound:               "Challenge not found",
	FailedVerificationLimitExceeded: "Exceeded failed verification limit",
	UserVerificationIncomplete:      "Complete user verification before attempting to login",
	VerificationFailed:              "Verification failed",
	TooManyChallengeRequests:        "Too many challenge requests",
	EmailIdInvalid:                  "Invalid email ID",
	PhoneNumberInvalid:              "Invalid phone number",

	// Cron
	MinuteIsInvalid:    "Invalid minute",
	HourIsInvalid:      "Invalid hour",
	DayOfWeekIsInvalid: "Invalid day of week",

	// Identity
	IdentityTypeNotFound: "Identity type not found",

	// Generic
	Unauthorised:        "Unauthorized",
	Forbidden:           "Forbidden",
	NotFound:            "Not Found",
	InternalServerError: "Internal Server Error",
	Conflict:            "Conflict",
	IdInvalid:           "Id Invalid",
}

var httpStatus = map[string]int{
	// User
	UserAlreadyExists:               http.StatusConflict,
	UserNotFound:                    http.StatusNotFound,
	UserNotRegistered:               http.StatusNotFound,
	CredentialsInvalid:              http.StatusUnauthorized,
	FailedLoginLimitExceeded:        http.StatusForbidden,
	ChallengeNotFound:               http.StatusNotFound,
	FailedVerificationLimitExceeded: http.StatusForbidden,
	UserVerificationIncomplete:      http.StatusForbidden,
	VerificationFailed:              http.StatusForbidden,
	TooManyChallengeRequests:        http.StatusTooManyRequests,
	EmailIdInvalid:                  http.StatusBadRequest,
	PhoneNumberInvalid:              http.StatusBadRequest,

	// Identity
	IdentityTypeNotFound: http.StatusNotFound,

	// Cron
	MinuteIsInvalid:    http.StatusBadRequest,
	HourIsInvalid:      http.StatusBadRequest,
	DayOfWeekIsInvalid: http.StatusBadRequest,

	// Generic
	Unauthorised:        http.StatusUnauthorized,
	Forbidden:           http.StatusForbidden,
	NotFound:            http.StatusNotFound,
	InternalServerError: http.StatusInternalServerError,
	Conflict:            http.StatusConflict,
	IdInvalid:           http.StatusUnprocessableEntity,
}

func Message(code string) string {
	if m, ok := messages[code]; ok {
		return m
	}
	return http.StatusText(http.StatusInternalServerError)
}

func HttpStatus(code string) int {
	if h, ok := httpStatus[code]; ok {
		return h
	}
	return http.StatusInternalServerError
}
