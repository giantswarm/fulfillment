package slack

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"

	"github.com/giantswarm/fulfillment/customer"
)

type Service interface {
	PostCustomerData(customerData customer.Data) error
}

func New(token string, mock bool) (Service, error) {
	if mock {
		return NewMockService(), nil
	}

	return NewService(token)
}

type Mock struct {
	Called bool
}

func NewMockService() *Mock {
	log.Printf("Using mock Slack service")

	return &Mock{}
}

func (m *Mock) PostCustomerData(customerData customer.Data) error {
	m.Called = true

	return nil
}

type Slack struct {
	client *slack.Client
}

func NewService(token string) (*Slack, error) {
	if token == "" {
		return nil, fmt.Errorf("Slack token is required")
	}

	client := slack.New(token)

	s := &Slack{
		client: client,
	}

	return s, nil
}

func (s *Slack) PostCustomerData(customerData customer.Data) error {
	message := fmt.Sprintf("A new customer has signed up with the email address %s, and the following data:", customerData.Email)

	dataAttachment := slack.Attachment{
		Title: "Customer Data",
		Fields: []slack.AttachmentField{
			{
				Title: "AWS Account ID",
				Value: customerData.AWSAccountId,
			},
			{
				Title: "Identifier",
				Value: customerData.Identifier,
			},
			{
				Title: "Product Code",
				Value: customerData.ProductCode,
			},
		},
	}

	entitlementAttachments := []slack.Attachment{}
	for i, entitlement := range customerData.Entitlements {
		entitlementAttachment := slack.Attachment{
			Title: fmt.Sprintf("Entitlement %v", i+1),
			Fields: []slack.AttachmentField{
				{
					Title: "Customer Identifier",
					Value: entitlement.CustomerIdentifier,
				},
				{
					Title: "Dimension",
					Value: entitlement.Dimension,
				},
				{
					Title: "Expiration Date",
					Value: entitlement.ExpirationDate.String(),
				},
				{
					Title: "Product Code",
					Value: entitlement.ProductCode,
				},
				{
					Title: "Value",
					Value: entitlement.Value,
				},
			},
		}

		entitlementAttachments = append(entitlementAttachments, entitlementAttachment)
	}

	footerAttachment := slack.Attachment{
		Footer: "Posted by https://github.com/giantswarm/fulfillment",
	}

	attachments := []slack.Attachment{}
	attachments = append(attachments, dataAttachment)
	attachments = append(attachments, entitlementAttachments...)
	attachments = append(attachments, footerAttachment)

	_, _, err := s.client.PostMessage(
		"#noise-fulfillment",
		slack.MsgOptionAsUser(true),
		slack.MsgOptionText(message, false),
		slack.MsgOptionAttachments(attachments...),
	)

	return err
}
