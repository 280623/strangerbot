package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"strangerbot/otpgateway"
	"strangerbot/otpgateway/models"
	"strangerbot/otpgateway/smtp"
)

var (
	ErrTooManyAttempts = errors.New("too many attempts")
)

const (
	alphaChars    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	numChars      = "0123456789"
	alphaNumChars = alphaChars + numChars
)

type OTPMaster struct {
	Namespace string
	Store     otpgateway.Store
	OTPTpl    *otpgateway.OTPTpl
	Emailer   *smtp.Emailer

	otpTTL         time.Duration
	otpMaxAttempts int
	OTPMaxLen      int
}

var OTPMasterIns *OTPMaster

func NewOTPMaster(namespace string, store otpgateway.Store, OTPTpl *otpgateway.OTPTpl, emailer *smtp.Emailer, otpTTL time.Duration, otpMaxAttempts int, OTPMaxLen int) *OTPMaster {
	return &OTPMaster{
		Namespace:      namespace,
		Store:          store,
		OTPTpl:         OTPTpl,
		Emailer:        emailer,
		otpTTL:         otpTTL,
		otpMaxAttempts: otpMaxAttempts,
		OTPMaxLen:      OTPMaxLen,
	}
}

func (o *OTPMaster) SendOTP(id string, to string) error {

	// validate to
	err := o.Emailer.ValidateAddress(to)
	if err != nil {
		return err
	}

	// generate a random ID
	if len(id) == 0 {
		id, err = otpgateway.GenerateRandomString(32, alphaNumChars)
		if err != nil {
			return err
		}
	}

	// Check if the OTP attempts have exceeded the quota.
	otp, err := o.Store.Check(o.Namespace, id, false)
	if err != nil && err != otpgateway.ErrNotExist {
		log.Printf("error checking OTP status: %v\n", err)
		return err
	}

	// There's an existing OTP that's locked.
	if err != otpgateway.ErrNotExist && isLocked(otp) {
		return ErrTooManyAttempts
	}

	// generate a random OTP
	otpVal, err := otpgateway.GenerateRandomString(o.OTPMaxLen, numChars)
	if err != nil {
		return err
	}

	// Create the OTP.
	newOTP, err := o.Store.Set(o.Namespace, id, models.OTP{
		OTP:         otpVal,
		To:          to,
		ChannelDesc: "",
		AddressDesc: "",
		Extra:       []byte(""),
		Provider:    "smtp",
		TTL:         o.otpTTL,
		MaxAttempts: o.otpMaxAttempts,
	})

	if err != nil {
		log.Printf("error setting OTP: %v\n", err)
		return err
	}

	if o.Emailer != nil {
		if err := otpgateway.OTPPush(newOTP, o.OTPTpl, o.Emailer, ""); err != nil {
			return err
		}
	}

	return nil
}

func (o *OTPMaster) VerifyOTP(id string, otp string) (models.OTP, error) {

	if len(otp) == 0 {
		return models.OTP{}, errors.New(fmt.Sprintf("ID should be min %d chars", o.OTPMaxLen))
	}

	// Check the OTP.
	out, err := o.Store.Check(o.Namespace, id, true)
	if err != nil {
		if err != otpgateway.ErrNotExist {
			return out, fmt.Errorf("error checking OTP: %v", err)
		}
		return out, err
	}

	errMsg := ""
	if isLocked(out) {
		errMsg = fmt.Sprintf("Too many attempts. Please retry after %0.f seconds.",
			out.TTL.Seconds())
	} else if out.OTP != otp {
		errMsg = "OTP does not match"
	}

	// There was an error.
	if errMsg != "" {
		return out, errors.New(errMsg)
	}

	if err := o.Store.Close(o.Namespace, id); err != nil {
		return out, err
	}

	out.Closed = true

	return out, nil
}

// isLocked tells if an OTP is locked after exceeding attempts.
func isLocked(otp models.OTP) bool {
	if otp.Attempts >= otp.MaxAttempts {
		return true
	}
	return false
}
