package aws

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/marketplaceentitlementservice"
	"github.com/aws/aws-sdk-go-v2/service/marketplacemetering"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/giantswarm/fulfillment/customer"
)

const (
	REGION                         = "us-east-1" // Note: the AWS Marketplace Entitlement Service only has one endpoint, which is in us-east-1
	CUSTOMER_IDENTIFIER_FILTER_KEY = "CUSTOMER_IDENTIFIER"
)

type Service interface {
	FetchCustomerData(token string) (customer.Data, error)
}

func New(accessKeyId, secretAccessKey string, mock bool) (Service, error) {
	if mock {
		return NewMockService(), nil
	}

	return NewService(accessKeyId, secretAccessKey)
}

type Mock struct {
	Called bool
}

func NewMockService() *Mock {
	log.Printf("Using mock AWS service")

	return &Mock{}
}

func (m *Mock) FetchCustomerData(token string) (customer.Data, error) {
	m.Called = true

	now := time.Now()

	customerData := customer.Data{
		AWSAccountId: "123456789012",
		Identifier:   "example-customer-identifier",
		ProductCode:  "example-product-code",

		Entitlements: []customer.Entitlement{
			{
				CustomerIdentifier: "example-customer-identifier",
				Dimension:          "example-dimension",
				ExpirationDate:     &now,
				ProductCode:        "example-product-code",
				Value:              "example-value",
			},
		},
	}

	return customerData, nil
}

type AWS struct {
	meteringClient     *marketplacemetering.Client
	entitlementService *marketplaceentitlementservice.Client
}

func NewService(accessKeyId, secretAccessKey string) (*AWS, error) {
	if accessKeyId == "" {
		return nil, fmt.Errorf("AWS access key ID is required")
	}
	if secretAccessKey == "" {
		return nil, fmt.Errorf("AWS secret access key is required")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, "")),
		config.WithRegion(REGION),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	meteringClient := marketplacemetering.NewFromConfig(cfg)
	entitlementService := marketplaceentitlementservice.NewFromConfig(cfg)

	s := &AWS{
		meteringClient:     meteringClient,
		entitlementService: entitlementService,
	}

	return s, nil
}

func (c *AWS) FetchCustomerData(token string) (customer.Data, error) {
	customerData := customer.Data{}

	resolveCustomerResult, err := c.meteringClient.ResolveCustomer(
		context.TODO(),
		&marketplacemetering.ResolveCustomerInput{
			RegistrationToken: aws.String(token),
		},
	)
	if err != nil {
		return customer.Data{}, fmt.Errorf("failed to resolve customer: %w", err)
	}

	customerData.AWSAccountId = *resolveCustomerResult.CustomerAWSAccountId
	customerData.Identifier = *resolveCustomerResult.CustomerIdentifier
	customerData.ProductCode = *resolveCustomerResult.ProductCode

	getEntitlementsResults, err := c.entitlementService.GetEntitlements(
		context.TODO(),
		&marketplaceentitlementservice.GetEntitlementsInput{
			ProductCode: aws.String(customerData.ProductCode),

			Filter: map[string][]string{
				CUSTOMER_IDENTIFIER_FILTER_KEY: {
					customerData.Identifier,
				},
			},
		},
	)
	if err != nil {
		return customer.Data{}, fmt.Errorf("failed to get entitlements: %w", err)
	}

	for _, entitlement := range getEntitlementsResults.Entitlements {
		var value string

		if entitlement.Value.BooleanValue != nil {
			value = strconv.FormatBool(*entitlement.Value.BooleanValue)
		}
		if entitlement.Value.DoubleValue != nil {
			value = strconv.FormatFloat(*entitlement.Value.DoubleValue, 'f', -1, 64)
		}
		if entitlement.Value.IntegerValue != nil {
			value = strconv.Itoa(int(*entitlement.Value.IntegerValue))
		}
		if entitlement.Value.StringValue != nil {
			value = *entitlement.Value.StringValue
		}

		customerData.Entitlements = append(customerData.Entitlements, customer.Entitlement{
			CustomerIdentifier: *entitlement.CustomerIdentifier,
			Dimension:          *entitlement.Dimension,
			ExpirationDate:     entitlement.ExpirationDate,
			ProductCode:        *entitlement.ProductCode,
			Value:              value,
		})
	}

	return customerData, nil
}
