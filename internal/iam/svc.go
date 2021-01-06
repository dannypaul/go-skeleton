package iam

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/dannypaul/go-skeleton/internal/config"
	"github.com/dannypaul/go-skeleton/internal/exception"
	"github.com/dannypaul/go-skeleton/internal/notification"
	"github.com/dannypaul/go-skeleton/internal/primitive"
	"github.com/dannypaul/go-skeleton/internal/repository"
)

type UserRepo interface {
	repository.Counter
	repository.Creator
	repository.Finder
	repository.Incrementer
	repository.Patcher
	repository.Setter
}

type ChallengeRepo interface {
	repository.Creator
	repository.Deleter
	repository.Finder
	repository.Incrementer
	repository.Setter
}

type Svc interface {
	VerifySeedUser(ctx context.Context) error
	Invite(ctx context.Context, req InviteReq) (User, error)
	Challenge(ctx context.Context, challenge Challenge) (Challenge, error)
	Verify(ctx context.Context, req VerifyReq) (Session, error)
	UpdatePassword(ctx context.Context, userId primitive.Id, req UpdatePasswordReq) (bool, error)
	Login(ctx context.Context, req LoginReq) (Session, error)
	FindMe(ctx context.Context) (User, error)
	FindUser(ctx context.Context, id primitive.Id) (User, error)
}

type svc struct {
	userRepo            UserRepo
	challengeRepo       ChallengeRepo
	notificationService notification.Svc
}

func NewService(userRepo UserRepo, challengeRepo ChallengeRepo, notificationService notification.Svc) Svc {
	return svc{
		userRepo:            userRepo,
		challengeRepo:       challengeRepo,
		notificationService: notificationService,
	}
}

type Claims struct {
	UserId      primitive.Id `json:"userId"`
	UserVersion int          `json:"userVersion"`
	Verified    bool         `json:"verified,omitempty"`
	Role        Role         `json:"role"`
	jwt.StandardClaims
}

func (s svc) FindUser(ctx context.Context, id primitive.Id) (User, error) {
	_, err := VerifySession(ctx, []Role{PlatformAdmin, MerchantAdmin})
	if err != nil {
		return User{}, err
	}

	copier, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return User{}, errors.New(exception.UserNotFound)
		}
		return User{}, err
	}

	var user User
	return user, copier.Copy(&user)
}

func (s svc) VerifySeedUser(ctx context.Context) error {
	conf, _ := config.Get()
	exist, err := s.DoesEmailIdExist(ctx, conf.SeedEmailId)
	if err != nil || exist {
		return err
	}
	_, err = s.userRepo.Create(ctx, User{
		Role: PlatformAdmin,
		Name: "Root administrator",
		Identities: []Identity{{
			Type:    EMAIL,
			EmailId: conf.SeedEmailId,
		}, {
			Type:  PHONE,
			Phone: &Phone{Number: conf.SeedPhoneNumber},
		}},
	})
	return err
}

func (s svc) DoesPhoneNumberExist(ctx context.Context, phoneNumber string) (bool, error) {
	count, err := s.userRepo.Count(ctx, []repository.Filter{{Key: "identities.phone.number", Value: phoneNumber}})
	if err != nil {
		return false, fmt.Errorf("could count the number of users with the given phone number %w", err)
	}

	return count > 0, nil
}

func (s svc) DoesEmailIdExist(ctx context.Context, emailId string) (bool, error) {
	count, err := s.userRepo.Count(ctx, []repository.Filter{{Key: "identities.emailId", Value: emailId}})
	if err != nil {
		return false, fmt.Errorf("could count the number of users with the given emailId %w", err)
	}

	return count > 0, nil
}

func (s svc) FindMe(ctx context.Context) (User, error) {
	claims, ok := ctx.Value(CtxClaimsKey).(Claims)
	if !ok || claims.UserId == "" {
		return User{}, errors.New(exception.Unauthorised)
	}

	copier, err := s.userRepo.FindById(ctx, claims.UserId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return User{}, errors.New(exception.UserNotFound)
		}
		return User{}, err
	}

	var user User
	return user, copier.Copy(&user)
}

func (s svc) FindUserByIdentity(ctx context.Context, identity Identity) (User, error) {
	if identity.Type != PHONE && identity.Type != EMAIL {
		return User{}, errors.New(exception.IdentityTypeNotFound)
	}

	var copier repository.Copier
	var err error
	if identity.Type == PHONE {
		filters := []repository.Filter{{Key: "identities.phone.number", Value: identity.Phone.Number}}
		copier, err = s.userRepo.FindSingle(ctx, filters)
	}

	if identity.Type == EMAIL {
		filters := []repository.Filter{{Key: "identities.emailId", Value: identity.EmailId}}
		copier, err = s.userRepo.FindSingle(ctx, filters)
	}

	if err != nil {
		return User{}, err
	}

	var user User
	err = copier.Copy(&user)

	return user, err
}

func (s svc) Invite(ctx context.Context, inviteReq InviteReq) (User, error) {
	_, err := VerifySession(ctx, []Role{PlatformAdmin})
	if err != nil {
		return User{}, err
	}

	phoneNumberExists, err := s.DoesPhoneNumberExist(ctx, inviteReq.Phone.Number)
	if err != nil {
		return User{}, fmt.Errorf("could not check if the user exists %w", err)
	}

	if phoneNumberExists {
		return User{}, errors.New(exception.UserAlreadyExists)
	}

	copier, err := s.userRepo.Create(ctx, User{
		Role: inviteReq.Role,
		Name: inviteReq.Name,
		Identities: []Identity{
			{Type: EMAIL, EmailId: inviteReq.EmailId},
			{Type: PHONE, Phone: &inviteReq.Phone},
		},
	})
	if err != nil {
		return User{}, fmt.Errorf("could not save the user to persistence %w", err)
	}

	var user User
	err = copier.Copy(&user)
	if err != nil {
		return User{}, fmt.Errorf("could not copy the persistence response to variable %w", err)
	}

	return user, nil
}

func (s svc) UpdatePassword(ctx context.Context, userId primitive.Id, req UpdatePasswordReq) (bool, error) {
	claims, err := VerifyActionToken(ctx)
	if err != nil {
		return false, err
	}

	copier, err := s.userRepo.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			return false, errors.New(exception.UserNotFound)
		}
		return false, err
	}

	var user User
	err = copier.Copy(&user)
	if err != nil {
		return false, fmt.Errorf("could not copy the persistence response to variable %w", err)
	}

	if claims.UserVersion != user.Version || claims.Verified == false {
		return false, fmt.Errorf(exception.Unauthorised)
	}

	if user.Id != claims.UserId {
		return false, errors.New(exception.Forbidden)
	}

	user.updatePassword(req.Password)

	patchers := []repository.Patch{
		{"$set", "password", user.Password},
		{"$inc", "version", 1},
	}
	err = s.userRepo.Patch(ctx, user.Id, patchers)

	return true, err
}

func VerifyActionToken(ctx context.Context) (Claims, error) {
	claims, ok := ctx.Value(CtxClaimsKey).(Claims)
	if !ok || claims.UserId == "" || claims.Role == "" {
		return Claims{}, errors.New(exception.Unauthorised)
	}
	return claims, nil
}

func VerifySession(ctx context.Context, hasAnyRole []Role) (Claims, error) {
	claims, ok := ctx.Value(CtxClaimsKey).(Claims)
	if !ok || claims.UserId == "" {
		return Claims{}, errors.New(exception.Unauthorised)
	}

	for _, role := range hasAnyRole {
		if role == claims.Role {
			return claims, nil
		}
	}

	return Claims{}, errors.New(exception.Forbidden)
}
