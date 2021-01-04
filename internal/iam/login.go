package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/dannypaul/go-skeleton/internal/exception"
)

func (s svc) Login(ctx context.Context, req LoginReq) (Session, error) {
	identity := Identity{
		Type:    req.IdentityType,
		EmailId: req.EmailId,
		Phone:   &req.Phone,
	}
	user, err := s.FindUserByIdentity(ctx, identity)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return Session{}, errors.New(exception.UserNotFound)
		}
		return Session{}, err
	}

	if user.FailedAuthAttempts >= 3 {
		return Session{}, errors.New(exception.FailedLoginLimitExceeded)
	}

	if user.Password == "" {
		return Session{}, errors.New(exception.UserVerificationIncomplete)
	}

	if !user.equalsPassword(req.Password) {
		err = s.userRepo.IncrementById(ctx, user.Id, "failedAuthAttempts", 1)
		return Session{}, errors.New(exception.CredentialsInvalid)
	}

	if user.FailedAuthAttempts > 0 {
		err = s.userRepo.SetById(ctx, user.Id, "failedAuthAttempts", 0)
		if err != nil {
			return Session{}, fmt.Errorf("could not reset the failed authentication attempt count %w", err)
		}
	}

	token, err := user.createToken(false)
	if err != nil {
		return Session{}, fmt.Errorf("could not create token for the user %w", err)
	}

	return Session{User: user, Token: token}, nil
}
