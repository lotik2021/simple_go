package notification

import (
	"bitbucket.movista.ru/maas/maasapi/internal/common"
	"bitbucket.movista.ru/maas/maasapi/internal/config"
	"bitbucket.movista.ru/maas/maasapi/internal/models"
)

type Notification struct {
	Reason  string `json:"reason"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type UNotification struct {
	Message      Message `json:"message"`
	ProviderName string  `json:"providerName"`
}

type Message struct {
	To          []string    `json:"to"`
	Cc          []string    `json:"cc, omitempty"`
	Bcc         []string    `json:"bcc, omitempty"`
	From        string      `json:"from"`
	Subject     string      `json:"subject"`
	IsBodyHTML  bool        `json:"isBodyHtml"`
	Attachments Attachments `json:"attachments, omitempty"`
	UserID      int         `json:"userId"`
	Text        string      `json:"text"`
}

type Attachments []struct {
	Content string `json:"content"`
	Type    string `json:"type"`
	Name    string `json:"name"`
}

func SendFeedback(ctx common.Context, req Notification) (*models.RawResponse, error) {
	unotification := UNotification{
		Message: Message{
			To:      []string{config.C.Notification.SupportEmail},
			From:    req.Email,
			Subject: req.Reason,
			Text:    req.Message,
			UserID:  1,
		},
		ProviderName: config.C.Notification.Provider,
	}

	return common.UapiAuthorizedPost(ctx, pushClient, unotification, config.C.Notification.Urls.Feedback)
}
