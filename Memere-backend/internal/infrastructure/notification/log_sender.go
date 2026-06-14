// Package notification holds infrastructure-level push/email sender implementations.
// LogPushSender and LogEmailSender are the development fallbacks used when FCM /
// SendGrid keys are absent. They simply log the outbound message so Phase 5 can
// run end-to-end locally without real provider credentials.
package notification

import (
	"context"
	"log"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/service"
)

// LogPushSender implements service.PushSender by logging the notification
// instead of calling FCM. Used when FCM_SERVER_KEY is absent.
type LogPushSender struct{}

var _ service.PushSender = LogPushSender{}

func (LogPushSender) Send(_ context.Context, tokens []string, ev service.NotifyEvent) ([]string, error) {
	log.Printf("[LogPushSender] push type=%s title=%q to %d device(s)", ev.Type, ev.Title, len(tokens))
	return nil, nil
}

// LogEmailSender implements service.EmailSender by logging the message instead
// of calling SendGrid. Used when SENDGRID_API_KEY is absent.
type LogEmailSender struct{}

var _ service.EmailSender = LogEmailSender{}

func (LogEmailSender) Send(_ context.Context, to, subject, html string) error {
	log.Printf("[LogEmailSender] email to=%s subject=%q", to, subject)
	return nil
}
