package service

import "context"

// Channel identifies a notification delivery mechanism.
type Channel string

const (
	ChannelPush  Channel = "push"
	ChannelEmail Channel = "email"
	ChannelInApp Channel = "in_app"
)

// NotifyEvent is the single envelope passed to the Notifier port.
// Callers never reference FCM or SendGrid directly — they build an event and
// the dispatcher decides which providers to call based on Channels and user
// preferences.
type NotifyEvent struct {
	UserID   string
	Type     string            // matches the §11.2 event catalog
	Title    string
	Body     string
	Data     map[string]string // deep-link payload for push / in-app
	Channels []Channel
}

// Notifier is the port for async notification dispatch. Implementations must
// return immediately (enqueue + in-app write); they must never block the
// calling API request path. A slow or failing provider must not propagate an
// error to the caller.
type Notifier interface {
	Notify(ctx context.Context, ev NotifyEvent) error
}

// PushSender delivers push notifications to one or more FCM tokens. Invalid
// tokens reported by FCM must be pruned from device_tokens by the caller.
type PushSender interface {
	Send(ctx context.Context, tokens []string, ev NotifyEvent) (invalidTokens []string, err error)
}

// EmailSender delivers a single email message.
type EmailSender interface {
	Send(ctx context.Context, to, subject, html string) error
}
