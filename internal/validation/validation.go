package validation

import (
	"errors"
	"github.com/dannypaul/go-skeleton/internal/exception"
	"regexp"
)

func ValidatePhone(phoneNumber string) error {
	if len(phoneNumber) != 10 {
		return errors.New(exception.PhoneNumberInvalid)
	}

	for i, c := range phoneNumber {
		if i == 0 {
			if c < '6' || c > '9' {
				return errors.New(exception.PhoneNumberInvalid)
			}
		}
		if c < '0' || c > '9' {
			return errors.New(exception.PhoneNumberInvalid)
		}
	}

	return nil
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func ValidateEmailId(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New(exception.EmailIdInvalid)
	}
	return nil
}
