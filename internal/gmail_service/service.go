package gmail_service

import (
	"fmt"
	"io"
	"os"
	"planpilot/internal/config"
	"planpilot/internal/logger"
	"planpilot/internal/openai"

	"github.com/BrianLeishman/go-imap"
)

type GmailService struct {
	im *imap.Dialer
}

func New() (*GmailService, error) {
	config := config.New()
	imap.Verbose = false //  flag to logging all the http.Request
	imap.RetryCount = 3

	im, err := imap.New(config.EMAIL_USERNAME, config.EMAIL_PASSWORD, config.EMAIL_IMAP_ADDRESS, config.EMAIL_IMAP_ADDRESS_PORT)
	if err != nil {
		return &GmailService{}, err
	}
	return &GmailService{im: im}, err
}

func (s *GmailService) GetLastEmails(emailsCount int) ([]*imap.Email, error) {
	var emails []*imap.Email

	// examples of folder{ "INBOX", "INBOX/My Folder", "Sent Items", "Deleted",}
	folder := "INBOX"

	err := s.im.SelectFolder(folder)
	if err != nil {
		logger.Error("error while selecting folder INBOX: ", err)
		return emails, err
	}

	uids, err := s.im.GetUIDs("ALL")
	if err != nil {
		logger.Error("error while getting email UIDs: ", err)
		return emails, err
	}
	cntEmails := len(uids)
	uids = uids[cntEmails-emailsCount : cntEmails]

	emailsMap, err := s.im.GetEmails(uids...)
	if err != nil {
		logger.Error("error while fetching emails: ", err)
		return emails, err
	}

	for uid, email := range emailsMap {
		logger.Info("processed email with UID: ", uid)
		emails = append(emails, email)
	}

	return emails, nil
}

func (s *GmailService) MakeCompressedEmailsText(emailsCount int) (string, error) {
	emails, err := s.GetLastEmails(emailsCount)
	if err != nil {
		return "", err
	}

	emailsPrompt := getEmailsPromptText()

	allEmailsText := ""
	for _, email := range emails {
		compressedEmail := openai.MakeOpenAICall(string(emailsPrompt) + email.Text)
		allEmailsText += fmt.Sprintf(config.EMAIL_TEMPLATE_RESPONSE, email.From, compressedEmail)
	}

	logger.Info("get the emails from gmailService: ", len(emails), " text: ", allEmailsText)

	return allEmailsText, nil
}

func getEmailsPromptText() []byte {
	// TODO: change the way to work work with this file (do not open in every time)
	file, err := os.Open("../prompts/emails_determine")

	if err != nil {
		logger.Error("Error open emails prompt file: ", err)
		return []byte("")
	}

	prompt, err := io.ReadAll(file)
	if err != nil {
		logger.Error("Error reading emails prompt file: ", err)
		return []byte("")
	}
	return prompt

}

// func extractEmailCount(result string) (int, error) {
// 	emailsLine := strings.Split(result, "\n")[2]
// 	emailCountString := strings.Trim(strings.Split(emailsLine, ": ")[1], " \n\r")
// 	emailsCount, err := strconv.Atoi(emailCountString)
// 	return emailsCount, err
// }
