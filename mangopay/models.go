package mangopay

type PlatformClient struct {
	PlatformType            string                 `json:"PlatformType"`
	ClientId                string                 `json:"ClientId"`
	Name                    string                 `json:"Name"`
	RegisteredName          string                 `json:"RegisteredName"`
	TechEmails              []string               `json:"TechEmails"`
	AdminEmails             []string               `json:"AdminEmails"`
	BillingEmails           []string               `json:"BillingEmails"`
	FraudEmails             []string               `json:"FraudEmails"`
	HeadquartersAddress     Address                `json:"HeadquartersAddress"`
	HeadquartersPhoneNumber string                 `json:"HeadquartersPhoneNumber"`
	TaxNumber               string                 `json:"TaxNumber"`
	PlatformCategorization  PlatformCategorization `json:"PlatformCategorization"`
	PlatformURL             string                 `json:"PlatformURL"`
	PlatformDescription     string                 `json:"PlatformDescription"`
	CompanyReference        string                 `json:"CompanyReference"`
	PrimaryThemeColour      string                 `json:"PrimaryThemeColour"`
	PrimaryButtonColour     string                 `json:"PrimaryButtonColour"`
	Logo                    string                 `json:"Logo"`
	CompanyNumber           string                 `json:"CompanyNumber"`
	MCC                     string                 `json:"MCC"`
}

type Address struct {
	AddressLine1 string `json:"AddressLine1"`
	AddressLine2 string `json:"AddressLine2"`
	City         string `json:"City"`
	Region       string `json:"Region"`
	PostalCode   string `json:"PostalCode"`
	Country      string `json:"Country"`
}

type PlatformCategorization struct {
	BusinessType string `json:"BusinessType"`
	Sector       string `json:"Sector"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
