package iam

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/dannypaul/go-skeleton/internal/primitive"
	"github.com/dannypaul/go-skeleton/internal/rest"
)

func Router(svc Svc) *chi.Mux {
	resource := resource{svc}

	router := chi.NewRouter()

	router.Post("/challenge", resource.challenge)
	router.Post("/challenges/{challengeId}/resend", resource.challenge)
	router.Post("/verify", resource.verify)

	router.Post("/login", resource.login)

	router.Get("/users/me", resource.findMe)
	router.Get("/users/{userId}", resource.findUser)
	router.Put("/users/{userId}/password", resource.updatePassword)

	return router
}

type resource struct {
	svc Svc
}

func (res resource) challenge(w http.ResponseWriter, r *http.Request) {
	var req Challenge
	if rest.DecodeReq(w, r, &req) != nil {
		return
	}

	startedVerification, err := res.svc.Challenge(r.Context(), req)
	rest.EncodeRes(w, r, startedVerification, err)
}

func (res resource) login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if rest.DecodeReq(w, r, &req) != nil {
		return
	}

	session, err := res.svc.Login(r.Context(), req)
	rest.EncodeRes(w, r, session, err)
}

func (res resource) verify(w http.ResponseWriter, r *http.Request) {
	var req VerifyReq
	if rest.DecodeReq(w, r, &req) != nil {
		return
	}

	session, err := res.svc.Verify(r.Context(), req)
	rest.EncodeRes(w, r, session, err)
}

func (res resource) updatePassword(w http.ResponseWriter, r *http.Request) {
	var req UpdatePasswordReq
	if rest.DecodeReq(w, r, &req) != nil {
		return
	}

	session, err := res.svc.UpdatePassword(r.Context(), primitive.Id(chi.URLParam(r, "userId")), req)
	rest.EncodeRes(w, r, session, err)
}

func (res resource) findMe(w http.ResponseWriter, r *http.Request) {
	user, err := res.svc.FindMe(r.Context())
	rest.EncodeRes(w, r, user, err)
}

func (res resource) findUser(w http.ResponseWriter, r *http.Request) {
	user, err := res.svc.FindUser(r.Context(), primitive.Id(chi.URLParam(r, "userId")))
	rest.EncodeRes(w, r, user, err)
}
