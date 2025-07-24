package email

import (
	"log"

	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog"
)

type ResendConfig struct {
	ApiKey string
	From   string
}

type ResendClient struct {
	client *resend.Client
	config ResendConfig
	logger *zerolog.Logger
}

func NewResendClient(cfg ResendConfig, logger *zerolog.Logger) *ResendClient {
	log.Println("Initializing Resend client with API key:", cfg.ApiKey)
	return &ResendClient{
		client: resend.NewClient(cfg.ApiKey),
		config: cfg,
		logger: logger,
	}
}

func (r *ResendClient) SendEmail(to, subject, html string) error {
	params := &resend.SendEmailRequest{
		From:    r.config.From,
		To:      []string{to},
		Subject: subject,
		Html:    html,
	}

	sent, err := r.client.Emails.Send(params)
	if err != nil {
		if r.logger != nil {
			r.logger.Error().Err(err).
				Str("to", to).
				Str("subject", subject).
				Str("from", r.config.From).
				Msg("Failed to send email")
		}
		return err
	}
	if r.logger != nil {
		r.logger.Info().
			Str("to", to).
			Str("subject", subject).
			Str("from", r.config.From).
			Str("id", sent.Id).
			Msg("Email sent successfully")
	}
	return nil
}
