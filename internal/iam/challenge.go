package iam

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dannypaul/go-skeleton/config"
	"github.com/dannypaul/go-skeleton/internal/exception"
	"github.com/dannypaul/go-skeleton/internal/kit"
	"github.com/dannypaul/go-skeleton/internal/primitive"
	"github.com/dannypaul/go-skeleton/internal/repository"
	"github.com/dannypaul/go-skeleton/internal/validation"
)

type Challenge struct {
	Id                      primitive.Id `bson:"_id,omitempty" json:"id"`
	CreatedAt               time.Time    `bson:"createdAt" json:"-"`
	UpdatedAt               time.Time    `bson:"updatedAt" json:"-"`
	IdentityType            IdentityType `bson:"identityType" json:"identityType"`
	OTP                     string       `bson:"otp" json:"-"`
	FailedVerificationCount int          `bson:"failedVerificationCount" json:"-"`

	// Phone
	Phone Phone `bson:"phone,omitempty" json:"phone"`

	// Email
	EmailId string `bson:"emailId,omitempty" json:"emailId"`
}

func (c Challenge) Validate() error {
	if c.IdentityType == PHONE {
		if err := validation.ValidatePhone(c.Phone.Number); err != nil {
			return err
		}
		c.EmailId = ""
	}

	if c.IdentityType == EMAIL {
		if err := validation.ValidateEmailId(c.EmailId); err != nil {
			return err
		}
		c.Phone = Phone{}
	}

	return nil
}

func (s svc) FindChallengeByIdentity(ctx context.Context, identity Identity) (Challenge, error) {
	if identity.Type != PHONE && identity.Type != EMAIL {
		return Challenge{}, errors.New(exception.IdentityTypeNotFound)
	}

	var copier repository.Copier
	var err error
	if identity.Type == PHONE {
		copier, err = s.challengeRepo.FindSingle(ctx, []repository.Filter{{Key: "phone.number", Value: identity.Phone.Number}})
	}

	if identity.Type == EMAIL {
		copier, err = s.challengeRepo.FindSingle(ctx, []repository.Filter{{Key: "emailId", Value: identity.EmailId}})
	}

	if err != nil {
		return Challenge{}, err
	}

	var challenge Challenge
	err = copier.Copy(&challenge)

	return challenge, err
}

func (s svc) Challenge(ctx context.Context, req Challenge) (Challenge, error) {
	if err := req.Validate(); err != nil {
		return Challenge{}, err
	}
	conf, _ := config.Get()

	now := time.Now()
	otp, err := kit.GenerateOTP(6)
	if err != nil {
		return Challenge{}, err
	}

	identity := Identity{
		Type:    req.IdentityType,
		EmailId: req.EmailId,
		Phone:   &req.Phone,
	}
	challenge, err := s.FindChallengeByIdentity(ctx, identity)
	if errors.Is(err, exception.ErrNotFound) {
		req.CreatedAt = now
		req.UpdatedAt = now
		req.OTP = otp

		copier, err := s.challengeRepo.Create(ctx, req)
		if err != nil {
			return Challenge{}, fmt.Errorf("could not save the challenge request to persistence %w", err)
		}

		var createdChallenge Challenge
		err = copier.Copy(&createdChallenge)
		if err != nil {
			return Challenge{}, fmt.Errorf("could not copy the challenge to local variable %w", err)
		}

		challenge = createdChallenge
	}

	if !challenge.UpdatedAt.Equal(challenge.CreatedAt) {
		duration := now.Sub(challenge.UpdatedAt)
		if duration.Seconds() < conf.ChallengeTTL.Seconds() {
			return Challenge{}, errors.New(exception.TooManyChallengeRequests)
		}
	}

	err = s.challengeRepo.SetAllById(ctx, challenge.Id, []repository.KeyValue{
		{"updatedAt", time.Now()},
		{"otp", otp},
	})
	if err != nil {
		return Challenge{}, err
	}

	if req.IdentityType == PHONE {
		err = s.notificationService.VerifyPhone(ctx, challenge.Phone.Number, otp)
		if err != nil {
			return Challenge{}, fmt.Errorf("could not send the verification OTP to device %w", err)
		}
	}

	if req.IdentityType == EMAIL {
		err = s.notificationService.VerifyEmailId(ctx, challenge.EmailId, otp)
		if err != nil {
			return Challenge{}, fmt.Errorf("could not send the verification OTP to emailId %w", err)
		}
	}

	return challenge, nil
}

func (s svc) Verify(ctx context.Context, req VerifyReq) (Session, error) {
	identity := Identity{
		Type:    req.IdentityType,
		EmailId: req.EmailId,
		Phone:   &req.Phone,
	}
	challenge, err := s.FindChallengeByIdentity(ctx, identity)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return Session{}, errors.New(exception.ChallengeNotFound)
		}
		return Session{}, fmt.Errorf("could not find challenge request for the given emailId %w", err)
	}

	if challenge.FailedVerificationCount >= 3 {
		_, err = s.challengeRepo.Delete(ctx, challenge.Id)
		if err != nil {
			return Session{}, fmt.Errorf("could not delete the challenge request after failed verification limit exceeded %w", err)
		}
		return Session{}, errors.New(exception.FailedVerificationLimitExceeded)
	}

	if challenge.OTP != req.OTP {
		err = s.challengeRepo.IncrementById(ctx, challenge.Id, "failedVerificationCount", 1)
		return Session{}, errors.New(exception.VerificationFailed)
	}

	user, err := s.FindUserByIdentity(ctx, identity)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return Session{}, errors.New(exception.UserNotRegistered)
		}
		return Session{}, fmt.Errorf("could not find the user by identity, %w", err)
	}

	user.Version += 1

	_, identityIndex, err := user.Identities.getIdentity(req.IdentityType)
	patchers := []repository.Patch{
		{"$set", "identities." + strconv.Itoa(identityIndex) + ".verified", true},
		{"$set", "failedAuthAttempts", 0},
		{"$inc", "version", 1},
	}
	err = s.userRepo.Patch(ctx, user.Id, patchers)
	if err != nil {
		return Session{}, fmt.Errorf("could not update the user %w", err)
	}

	token, err := user.createToken(true)
	if err != nil {
		return Session{}, fmt.Errorf("could not create token for the user %w", err)
	}

	_, err = s.challengeRepo.Delete(ctx, challenge.Id)
	if err != nil {
		return Session{}, fmt.Errorf("could not delete the challenge request %w", err)
	}

	return Session{User: user, Token: token}, err
}
