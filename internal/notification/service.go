package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Svc interface {
	VerifyPhone(ctx context.Context, phoneNumber string, otp string) error
	VerifyEmailId(ctx context.Context, emailId string, otp string) error
}

type svc struct {
}

func NewService() Svc {
	return svc{}
}

type MailGunReq struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	Html    string `json:"html"`
}

func (s svc) VerifyEmailId(ctx context.Context, emailId string, otp string) error {
	verifyUserUri, err := url.Parse("https://googe.com")
	if err != nil {
		return err
	}

	verifyUserUri.Path = path.Join(verifyUserUri.Path, "/auth/verify")
	query := verifyUserUri.Query()
	query.Set("otp", otp)
	verifyUserUri.RawQuery = query.Encode()

	domain := "mail.app-name.com"
	from := "no-reply@" + domain
	subject := "Please verify your app-name account"
	html := "<div>Hello from app-name. <a href=" + verifyUserUri.String() + ">Verify</a></div>"
	to := emailId

	mailgunUri, err := url.Parse("https://api.mailgun.net")
	if err != nil {
		return err
	}

	mailgunUri.Path = path.Join(mailgunUri.Path, "/v3/"+domain+"/messages")

	mailgunReq := MailGunReq{
		From:    from,
		To:      to,
		Subject: subject,
		Text:    html,
		Html:    html,
	}

	body, err := json.Marshal(mailgunReq)
	req, err := http.NewRequest("POST", mailgunUri.String(), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic api:key")

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}

func (s svc) VerifyPhone(ctx context.Context, phoneNumber string, otp string) error {
	uri, err := url.Parse("https://api.textlocal.in")
	if err != nil {
		return err
	}

	uri.Path = path.Join(uri.Path, "send")

	query := uri.Query()
	query.Set("apikey", "apikey")
	query.Set("numbers", phoneNumber)
	query.Set("message", "Your app-name login OTP is "+otp)
	query.Set("sender", "TXTLCL")
	uri.RawQuery = query.Encode()

	emptyBody := bytes.NewBuffer([]byte("{}"))

	req, err := http.NewRequest("POST", uri.String(), emptyBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}
