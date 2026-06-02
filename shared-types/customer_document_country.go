package types

import (
	"encoding/json"
	"fmt"
)

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
	if len(raw) == 4 && raw[1] == 'B' && raw[2] == 'R' {
		return CustomerDocumentCountryBR, nil
	}
	if len(raw) < 2 || raw[0] != '"' || raw[len(raw)-1] != '"' {
		return nil, fmt.Errorf("expected JSON string for CustomerDocumentCountry, got %s", raw)
	}
	unquoted := raw[1 : len(raw)-1]
	if interned, ok := customerDocumentCountryIntern[string(unquoted)]; ok {
		return interned, nil
	}
	s := string(unquoted)
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
