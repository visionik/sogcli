//go:build integration

package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visionik/sogcli/internal/imap"
	"github.com/visionik/sogcli/internal/smtp"
)

// End-to-end test: send email via SMTP, read via IMAP

func TestEndToEndSendReceive(t *testing.T) {
	// Send via SMTP
	smtpClient := smtp.NewClient(smtp.Config{
		Host:     "localhost",
		Port:     3025,
		NoTLS:    true,
		Email:    "e2e@test.com",
		Password: "test",
	})

	testSubject := "E2E Test " + time.Now().Format("15:04:05")
	msg := &smtp.Message{
		From:    "e2e@test.com",
		To:      []string{"e2e@test.com"},
		Subject: testSubject,
		Body:    "This is an end-to-end test message.",
	}

	err := smtpClient.Send(msg)
	require.NoError(t, err)

	// Small delay for message to be processed
	time.Sleep(100 * time.Millisecond)

	// Read via IMAP
	imapClient, err := imap.Connect(imap.Config{
		Host:     "localhost",
		Port:     3143,
		NoTLS:    true,
		Email:    "e2e@test.com",
		Password: "test",
	})
	require.NoError(t, err)
	defer imapClient.Close()

	// List messages
	messages, err := imapClient.ListMessages("INBOX", 10, false)
	require.NoError(t, err)
	require.NotEmpty(t, messages)

	// Find our message
	var found bool
	var foundUID uint32
	for _, m := range messages {
		if m.Subject == testSubject {
			found = true
			foundUID = m.UID
			break
		}
	}
	assert.True(t, found, "should find the sent message")

	if found {
		// Get the full message
		fullMsg, err := imapClient.GetMessage("INBOX", foundUID, false)
		require.NoError(t, err)
		assert.Equal(t, testSubject, fullMsg.Subject)

		// Test flag operations
		err = imapClient.SetFlag("INBOX", foundUID, "seen", true)
		require.NoError(t, err)

		err = imapClient.SetFlag("INBOX", foundUID, "flagged", true)
		require.NoError(t, err)

		err = imapClient.SetFlag("INBOX", foundUID, "flagged", false)
		require.NoError(t, err)

		// Test copy
		_ = imapClient.CreateFolder("E2ETest")
		err = imapClient.CopyMessage("INBOX", foundUID, "E2ETest")
		require.NoError(t, err)

		// Verify copy
		copied, err := imapClient.ListMessages("E2ETest", 10, false)
		require.NoError(t, err)
		assert.NotEmpty(t, copied)

		// Cleanup
		if len(copied) > 0 {
			_ = imapClient.DeleteMessage("E2ETest", copied[0].UID)
		}
		_ = imapClient.DeleteFolder("E2ETest")

		// Delete original
		err = imapClient.DeleteMessage("INBOX", foundUID)
		require.NoError(t, err)
	}
}

func TestSearchWithMessages(t *testing.T) {
	// Send a message first
	smtpClient := smtp.NewClient(smtp.Config{
		Host:     "localhost",
		Port:     3025,
		NoTLS:    true,
		Email:    "search@test.com",
		Password: "test",
	})

	err := smtpClient.Send(&smtp.Message{
		From:    "search@test.com",
		To:      []string{"search@test.com"},
		Subject: "Searchable Message",
		Body:    "This message is for search testing.",
	})
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Search via IMAP
	imapClient, err := imap.Connect(imap.Config{
		Host:     "localhost",
		Port:     3143,
		NoTLS:    true,
		Email:    "search@test.com",
		Password: "test",
	})
	require.NoError(t, err)
	defer imapClient.Close()

	// Search ALL
	messages, err := imapClient.SearchMessages("INBOX", "ALL", 10)
	require.NoError(t, err)
	assert.NotEmpty(t, messages)
}
