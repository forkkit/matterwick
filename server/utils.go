package server

import (
	"fmt"

	"github.com/mattermost/mattermost-server/v5/mlog"
	"github.com/mattermost/matterwick/model"
)

func (s *Server) logErrorToMattermost(msg string, args ...interface{}) {
	if s.Config.MattermostWebhookURL == "" {
		mlog.Warn("No Mattermost webhook URL set: unable to send message")
		return
	}

	webhookMessage := fmt.Sprintf(msg, args...)
	mlog.Debug("Sending Mattermost message", mlog.String("message", webhookMessage))

	if s.Config.MattermostWebhookFooter != "" {
		webhookMessage += "\n---\n" + s.Config.MattermostWebhookFooter
	}

	webhookRequest := &WebhookRequest{Username: "MatterWick", Text: webhookMessage}

	if err := s.sendToWebhook(webhookRequest); err != nil {
		mlog.Error("Unable to post to Mattermost webhook", mlog.Err(err))
	}
}

func (s *Server) logPrettyErrorToMattermost(msg string, pr *model.PullRequest, err error, additionalFields map[string]string) {
	if s.Config.MattermostWebhookURL == "" {
		mlog.Warn("No Mattermost webhook URL set: unable to send message")
		return
	}

	mlog.Debug("Sending Mattermost message", mlog.String("message", msg))

	fullMessage := fmt.Sprintf("%s\n---\nError: %s\nRepository: %s/%s\nPull Request: %d [ status=%s ]\nURL: %s\n",
		msg,
		err,
		pr.RepoOwner, pr.RepoName,
		pr.Number, pr.State,
		pr.URL,
	)
	for key, value := range additionalFields {
		fullMessage = fullMessage + fmt.Sprintf("%s: %s\n", key, value)
	}
	fullMessage = fullMessage + s.Config.MattermostWebhookFooter

	webhookRequest := &WebhookRequest{Username: "MatterWick", Text: fullMessage}

	if err := s.sendToWebhook(webhookRequest); err != nil {
		mlog.Error("Unable to post to Mattermost webhook", mlog.Err(err))
	}
}

// NewBool return a bool pointer
func NewBool(b bool) *bool { return &b }

// NewInt return an int pointer
func NewInt(n int) *int { return &n }

// NewInt64 return an int64 pointer
func NewInt64(n int64) *int64 { return &n }

// NewInt32 return an int32 pointer
func NewInt32(n int32) *int32 { return &n }

// NewString return a string pointer
func NewString(s string) *string { return &s }
