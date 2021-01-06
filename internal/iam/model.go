package iam

import (
	"time"

	"github.com/dannypaul/go-skeleton/internal/config"
	"github.com/dannypaul/go-skeleton/internal/exception"
	"github.com/dannypaul/go-skeleton/internal/primitive"
	"github.com/dannypaul/go-skeleton/internal/repository"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

const CtxClaimsKey = "claims"

type IdentityType string

const (
	EMAIL IdentityType = "EMAIL"
	PHONE IdentityType = "PHONE"
)

type Phone struct {
	Number string `bson:"number,omitempty" json:"number"`
}

type Identity struct {
	Id       primitive.Id `bson:"_id,omitempty" json:"id"`
	Type     IdentityType `bson:"type" json:"type"`
	Verified bool         `bson:"verified" json:"verified"`
	EmailId  string       `bson:"emailId,omitempty" json:"emailId"`
	Phone    *Phone       `bson:"phone,omitempty" json:"phone"`
}

type IdentityList []Identity

func (i IdentityList) getIdentity(identityType IdentityType) (Identity, int, error) {
	for index, identity := range i {
		if identity.Type == identityType {
			return identity, index, nil
		}
	}
	return Identity{}, 0, exception.ErrNotFound
}

type Role string

const (
	PlatformAdmin Role = "PLATFORM_ADMIN"
	MerchantAdmin Role = "MERCHANT_ADMIN"
)

type User struct {
	Id                 primitive.Id `bson:"_id,omitempty" json:"id"`
	Role               Role         `bson:"role" json:"role"`
	Name               string       `bson:"name" json:"name"`
	Version            int          `bson:"version" json:"version"`
	FailedAuthAttempts int          `bson:"failedAuthAttempts" json:"-"`
	Identities         IdentityList `bson:"identities" json:"identities"`
	Password           string       `bson:"password" json:"-"`
}

func (u User) equalsPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

func (u *User) updatePassword(newPassword string) {
	//TODO: handle the error returned by GenerateFromPassword
	bytes, _ := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
	u.Password = string(bytes)
}

type UserList struct {
	Users []User          `json:"users"`
	Page  repository.Page `json:"page"`
}

func (u User) createToken(verified bool) (string, error) {
	conf, _ := config.Get()
	claims := &Claims{
		UserId:         u.Id,
		UserVersion:    u.Version,
		Role:           u.Role,
		Verified:       verified,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().UTC().Add(conf.JwtTTL).Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.JwtSecret))
}

type Session struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
