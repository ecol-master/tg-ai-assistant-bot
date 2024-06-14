package gmail_service

import (
	"testing"
)

func TestService(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Error("error in creating GmailService: ", err)
	}
}
