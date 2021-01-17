package kit

import "testing"

func TestOtpLength(t *testing.T) {
	otpLength := 5
	otp, _ := GenerateOTP(otpLength)
	if len(otp) != otpLength {
		t.Errorf("OTP length was incorrect, got: %d, want: %d.", len(otp), otpLength)
	}
}
