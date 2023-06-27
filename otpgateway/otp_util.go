package otpgateway

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"

	"strangerbot/otpgateway/models"
	"strangerbot/otpgateway/smtp"
)

// generateRandomString generates a cryptographically random,
// alphanumeric string of length n.
func GenerateRandomString(totalLen int, chars string) (string, error) {
	bytes := make([]byte, totalLen)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for k, v := range bytes {
		bytes[k] = chars[v%byte(len(chars))]
	}
	return string(bytes), nil
}

type OTPTpl struct {
	Subject *template.Template
	Tpl     *template.Template
}

func LoadProviderTemplates(tplFile, subj string) (*OTPTpl, error) {

	var (
		err          error
		tpl, subjTpl *template.Template
	)

	if len(tplFile) > 0 {
		// Parse the template file.
		tpl, err = template.ParseFiles(tplFile)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s err: %v", tplFile, err)
		}
	}

	if subj != "" {
		subjTpl, err = template.New("subject").Parse(subj)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %v", err)
		}
	}

	return &OTPTpl{
		Subject: subjTpl,
		Tpl:     tpl,
	}, nil
}

type pushTpl struct {
	To        string
	Namespace string
	OTP       string
}

func OTPPush(otp models.OTP, tpl *OTPTpl, p *smtp.Emailer, rootURL string) error {

	var (
		subj = &bytes.Buffer{}
		out  = &bytes.Buffer{}

		data = pushTpl{
			Namespace: otp.Namespace,
			To:        otp.To,
			OTP:       otp.OTP,
		}
	)

	if tpl.Subject != nil {
		if err := tpl.Subject.Execute(subj, data); err != nil {
			return err
		}
	}

	if tpl.Tpl != nil {
		if err := tpl.Tpl.Execute(out, data); err != nil {
			return err
		}
	}

	return p.Push(otp, subj.String(), out.Bytes())
}
