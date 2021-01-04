package iam

type VerifyReq struct {
	IdentityType IdentityType `json:"identityType"`

	EmailId string `json:"emailId"`
	Phone   Phone  `json:"phone"`

	OTP string `json:"otp"`
}

type InviteReq struct {
	Name    string `json:"name"`
	Phone   Phone  `json:"phone"`
	EmailId string `json:"emailId"`
	Role    Role   `json:"role"`
}

type UpdatePasswordReq struct {
	IdentityType IdentityType `json:"identityType"`
	Password     string       `json:"password"`
}

type LoginReq struct {
	IdentityType IdentityType `json:"identityType"`

	EmailId string `json:"emailId"`
	Phone   Phone  `json:"phone"`

	Password string `json:"password"`
}
