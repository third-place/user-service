package service

import (
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/third-place/user-service/internal/model"
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
	fromMail = mail.NewEmail("Example User", "test@example.com")
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

func (m *MailService) SendEmailVerify(otp *model.Otp) (*rest.Response, error) {
	name := m.getSenderName(otp)
	link := m.createVerifyLink(otp)
	return m.client.Send(
		mail.NewSingleEmail(
			fromMail,
			"Third place: email verification",
			mail.NewEmail(name, otp.User.Email),
			"Your Third place verification code is "+otp.Code+"\n"+
				"Copy and paste the link to verify your email address now: "+link,
			"<p>Your Third place verification code is "+otp.Code+"</p>"+
				"<p><a href=\""+link+"\">Click here to verify your email address</a></p>",
		),
	)
}

func (m *MailService) SendPasswordReset(otp *model.Otp) (*rest.Response, error) {
	name := m.getSenderName(otp)
	link := m.createPasswordResetLink(otp)
	return m.client.Send(
		mail.NewSingleEmail(
			fromMail,
			"Third place: password reset request",
			mail.NewEmail(name, otp.User.Email),
			"Your Third place verification code is "+otp.Code+"\n"+
				"Copy and paste the link to verify your email address now: "+link,
			"<p>Your Third place verification code is "+otp.Code+"</p>"+
				"<p><a href=\""+link+"\">Click here to reset your password</a></p>",
		),
	)
}

func (m *MailService) getSenderName(otp *model.Otp) string {
	user := otp.User
	name := "New User"
	if user.Name != "" {
		name = user.Name
	}
	return name
}

func (m *MailService) createVerifyLink(otp *model.Otp) string {
	return "https://thirdplaceapp.com/otp/?email=" + otp.User.Email + "&code=" + otp.Code
}

func (m *MailService) createPasswordResetLink(otp *model.Otp) string {
	return "https://thirdplaceapp.com/forgot-password/?email=" + otp.User.Email + "&code=" + otp.Code
}
