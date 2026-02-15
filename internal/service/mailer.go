package service

import "context"

type Mailer interface {
	Send(ctx context.Context, to string, subject string, body string) error
}

// NoopMailer for development when no mail server is configured
type NoopMailer struct{}

func (m *NoopMailer) Send(ctx context.Context, to string, subject string, body string) error {
	// Just log it
	return nil
}
