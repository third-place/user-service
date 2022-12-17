package service

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/third-place/user-service/internal/entity"
	"os"
)

type MailClient interface {
	Send(email *mail.SGMailV3) (*rest.Response, error)
}

type TestMailClient struct{}

func (t *TestMailClient) Send(email *mail.SGMailV3) (*rest.Response, error) {
	return nil, nil
}

type MailService struct {
	client MailClient
}

var fromMail *mail.Email

func init() {
	fromMail = mail.NewEmail("ThirdplaceBot", "info@thirdplaceapp.com")
}

func CreateMailService() *MailService {
	return &MailService{
		client: sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY")),
	}
}

func CreateTestMailService() *MailService {
	return &MailService{
		client: &TestMailClient{},
	}
}

func (m *MailService) SendVerificationEmail(user *entity.User) (*rest.Response, error) {
	name := m.getSenderName(user)
	link := m.createVerifyLink(user)
	return m.client.Send(
		mail.NewSingleEmail(
			fromMail,
			"Third place: email verification",
			mail.NewEmail(name, user.Email),
			"Your Third place verification code is "+user.OTP+"\n"+
				"Copy and paste the link to verify your email address now: "+link,
			"<p>Your Third place verification code is "+user.OTP+"</p>"+
				"<p><a href=\""+link+"\">Click here to verify your email address</a></p>",
		),
	)
}

func (m *MailService) SendPasswordResetEmail(user *entity.User) (*rest.Response, error) {
	name := m.getSenderName(user)
	link := m.createPasswordResetLink(user)
	return m.client.Send(
		mail.NewSingleEmail(
			fromMail,
			"Third place: password reset request",
			mail.NewEmail(name, user.Email),
			"Your Third place verification code is "+user.OTP+"\n"+
				"Copy and paste the link to verify your email address now: "+link,
			"<p>Your Third place verification code is "+user.OTP+"</p>"+
				"<p><a href=\""+link+"\">Click here to reset your password</a></p>",
		),
	)
}

func (m *MailService) getSenderName(user *entity.User) string {
	name := "New User"
	if user.Name != "" {
		name = user.Name
	}
	return name
}

func (m *MailService) createVerifyLink(user *entity.User) string {
	return "https://thirdplaceapp.com/otp/?email=" + user.Email + "&code=" + user.OTP
}

func (m *MailService) createPasswordResetLink(user *entity.User) string {
	return "https://thirdplaceapp.com/forgot-password/?email=" + user.Email + "&code=" + user.OTP
}
