// Package melissa is a simple wrapper around Melissa Data's GlobalAddress service.
package melissa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const globalAddressURL = "https://address.melissadata.net/v3/WEB/GlobalAddress/doGlobalAddress"

var (
	// Transmission code mappings
	TransmissionCodes = map[string]string{
		"SE01": http.StatusText(http.StatusInternalServerError),
		"GE01": "empty request structure",
		"GE02": "empty request record structure",
		"GE03": "records per request exceeded",
		"GE04": "empty CustomerID",
		"GE05": "invalid CustomerID",
		"GE06": "disabled CustomerID",
		"GE07": http.StatusText(http.StatusBadRequest),
		"GE08": "invalid CustomerID for product",
	}
	// Result code mappings
	ResultCodes = map[string]string{
		"AE01": "No Verification",
		"AE02": "Unknown Street",
		"AE03": "Component Error",
		"AE05": "Multiple Matches",
		"AE08": "SubPremises Number Invalid",
		"AE09": "SubPremises Number Missing",
		"AE10": "Premises Number Invalid",
		"AE11": "Premises Number Missing",
		"AE12": "PO Box Number Invlalid",
		"AE13": "PO Box Number Missing",
		"AE14": "Private Mail Box Missing",
		"AE17": "SubPremises Not Required",

		"AC01": "PostalCode",
		"AC02": "Administrative Area",
		"AC03": "Locality",
		"AC09": "Dependent Locality",
		"AC10": "Thoroughfare Name",
		"AC11": "Thoroughfare Type",
		"AC12": "Thoroughfare Direction",
		"AC13": "SubPremises Type",
		"AC14": "SubPremises Number",
		"AC15": "DoubleDependent Locality",
		"AC16": "SubAdministrative Area",
		"AC17": "SubNational Area",
	}
	// Geocode mappings
	GeoCodes = map[string]string{
		"GS01": "Geocoded to ZIP+4 (U.S.) or 6-digit Postal Code (Canada) Centroid",
		"GS02": "Geocoded to ZIP+2 Centroid",
		"GS03": "Geocoded to 5-digit (U.S.) or 3-digit (Canada) ZIP Code Centroid",
		"GS05": "Geocoded to 11-digit Rooftop level",
		"GS06": "Geocoded to 11-digit Interpolated Rooftop level",
		"GE01": "Invalid ZIP Code entered",
		"GE02": "Zip Code not found",
	}
	// Address code mappings (United States)
	AddressCodesUS = map[string]string{
		"A": "Alias",
		"F": "Firm or Company",
		"G": "General Delivery",
		"H": "Highrise or Business Complex",
		"P": "PO Box",
		"R": "Rural Route",
		"S": "Street of Residential",
	}
	// Address code mappings (Canada)
	AddressCodesCA = map[string]string{
		"1": "Street",
		"2": "Street Served by Route and GD",
		"3": "Lock Box",
		"4": "Route Service",
		"5": "General Delivery",
		"B": "LVR Street",
		"C": "Government Street",
		"D": "LVR Lock Box",
		"E": "Government Lock Box",
		"L": "LVR General Delivery",
		"K": "Building",
	}
)

// Client used to communicated with Melissa Data's GlobalAddress service.
type Client struct {
	client http.Client
	urlStr string
	key    string
}

// Melissa Data response type mapping
type Response struct {
	Records               []Record
	TotalRecords          string
	TransmissionReference string
	TransmissionResults   string
	Version               string
}

// Melissa Data record type mapping
type Record struct {
	AddressKey                         string
	AddressLine1                       string
	AddressLine2                       string
	AddressLine3                       string
	AddressLine4                       string
	AddressLine5                       string
	AddressLine6                       string
	AddressLine7                       string
	AddressLine8                       string
	AddressType                        string
	AdministrativeArea                 string
	Building                           string
	CountryISO3166_1_Alpha2            string
	CountryISO3166_1_Alpha3            string
	CountryISO3166_1_Numeric           string
	CountryName                        string
	DependentLocality                  string
	DependentThoroughfare              string
	DependentThoroughfareLeadingType   string
	DependentThoroughfareName          string
	DependentThoroughfarePostDirection string
	DependentThoroughfarePreDirection  string
	DependentThoroughfareTrailingType  string
	DoubleDependentLocality            string
	FormattedAddress                   string
	Latitude                           string
	Locality                           string
	Longitude                          string
	Organization                       string
	PostBox                            string
	PostalCode                         string
	PremisesNumber                     string
	PremisesType                       string
	RecordID                           string
	Results                            string
	SubAdministrativeArea              string
	SubNationalArea                    string
	SubPremises                        string
	SubPremisesNumber                  string
	SubPremisesType                    string
	Thoroughfare                       string
	ThoroughfareLeadingType            string
	ThoroughfareName                   string
	ThoroughfarePostDirection          string
	ThoroughfarePreDirection           string
	ThoroughfareTrailingType           string
}

// Ping simply hits the base URL for the GlobalAddress endpoint to ensure there is connectivity.
func (c Client) Ping() error {
	resp, err := http.Get(globalAddressURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response code, %d, received for ping", resp.StatusCode)
	}
	return nil
}

// Query invokes a JSON request to Melissa data using the given `qs` url.Values
// as the query params. A populated Response object is returned only when there are no errors.
func (c Client) Query(qs url.Values) (Response, error) {
	var r Response
	// Gets the query-string, excluding empty values from the address.
	qs.Add("id", c.key)
	urlStr := fmt.Sprintf("%s?%s", c.urlStr, qs.Encode())
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return r, err
	}

	// Invoke a JSON request.
	req.Header.Add("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return r, err
	}
	defer resp.Body.Close()

	// TODO check response status code for 200
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}

	// Read and transform data.
	err = json.Unmarshal(data, &r)
	if err != nil {
		return r, err
	}
	return r, err
}

// NewClient returns a new client using the given `apiKey` as the private key.
func NewClient(apiKey string) Client {
	client := http.Client{}
	return Client{
		client: client,
		urlStr: globalAddressURL,
		key:    apiKey,
	}
}
