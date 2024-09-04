package customer

import "time"

type Entitlement struct {
	CustomerIdentifier string
	Dimension          string
	ExpirationDate     *time.Time
	ProductCode        string
	Value              string
}

type Data struct {
	Email string

	AWSAccountId string
	Identifier   string
	ProductCode  string

	Entitlements []Entitlement
}
