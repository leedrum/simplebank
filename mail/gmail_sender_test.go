package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGmailSender(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	sender := NewGmailSender(
		"leedrum",
		"developer.levantung@gmail.com",
		"app_password",
	)
	err := sender.SendEmail(
		"Hello",
		"Hello, this is a test email",
		[]string{"developer.levantung1@gmail.com"},
		nil, nil,
		[]string{"../README.md"},
	)
	require.NoError(t, err)
}
