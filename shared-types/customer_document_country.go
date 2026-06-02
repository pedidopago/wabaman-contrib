package types

import "encoding/json"

type CustomerDocumentCountry *string

var (
	customerDocumentCountryBR = "BR"
	customerDocumentCountryUS = "US"
	customerDocumentCountryAR = "AR"
	customerDocumentCountryCL = "CL"
	customerDocumentCountryCO = "CO"
	customerDocumentCountryMX = "MX"
	customerDocumentCountryPY = "PY"
	customerDocumentCountryUY = "UY"

	CustomerDocumentCountryBR  CustomerDocumentCountry = &customerDocumentCountryBR
	CustomerDocumentCountryUS  CustomerDocumentCountry = &customerDocumentCountryUS
	CustomerDocumentCountryAR  CustomerDocumentCountry = &customerDocumentCountryAR
	CustomerDocumentCountryCL  CustomerDocumentCountry = &customerDocumentCountryCL
	CustomerDocumentCountryCO  CustomerDocumentCountry = &customerDocumentCountryCO
	CustomerDocumentCountryMX  CustomerDocumentCountry = &customerDocumentCountryMX
	CustomerDocumentCountryPY  CustomerDocumentCountry = &customerDocumentCountryPY
	CustomerDocumentCountryUY  CustomerDocumentCountry = &customerDocumentCountryUY
	CustomerDocumentCountryDel CustomerDocumentCountry = &delStr
)

var customerDocumentCountryIntern = map[string]CustomerDocumentCountry{
	"BR":   CustomerDocumentCountryBR,
	"US":   CustomerDocumentCountryUS,
	"AR":   CustomerDocumentCountryAR,
	"CL":   CustomerDocumentCountryCL,
	"CO":   CustomerDocumentCountryCO,
	"MX":   CustomerDocumentCountryMX,
	"PY":   CustomerDocumentCountryPY,
	"UY":   CustomerDocumentCountryUY,
	"$del": CustomerDocumentCountryDel,
}

func internCustomerDocumentCountry(raw json.RawMessage) (CustomerDocumentCountry, error) {
	if len(raw) == 4 && raw[0] == '"' && raw[1] == 'B' && raw[2] == 'R' && raw[3] == '"' {
		return CustomerDocumentCountryBR, nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	if interned, ok := customerDocumentCountryIntern[s]; ok {
		return interned, nil
	}
	return &s, nil
}

func EqualCustomerDocumentCountry(a, b CustomerDocumentCountry) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
