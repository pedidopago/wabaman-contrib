package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/pedidopago/go-common/util"
	"github.com/pedidopago/wabaman-contrib/rest"
)

const (
	DefaultBaseURL = "https://api.first.v2.pedidopago.com.br/wabaman"
)

type Client struct {
	JWT     string
	BaseURL string
}

func (c *Client) NewMessage(ctx context.Context, req *rest.NewMessageRequest) (*rest.NewMessageResponse, error) {
	output := &rest.NewMessageResponse{}
	if err := c.post(ctx, "/api/v1/message", req, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *Client) UpdateContact(ctx context.Context, contactID uint64, req *rest.UpdateContactRequest) (*rest.UpdateContactResponse, error) {
	resp := &rest.UpdateContactResponse{}
	if err := c.put(ctx, fmt.Sprintf("/api/v1/contact/%d", contactID), req, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetContacts(ctx context.Context, req *rest.GetContactsRequest) (*rest.GetContactsResponse, error) {
	if c == nil {
		return nil, fmt.Errorf("nil client")
	}
	q := req.BuildQuery()
	qenc := q.Encode()
	resp := &rest.GetContactsResponse{}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/contacts?%s", qenc), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetContactsV2(ctx context.Context, req rest.GetContactsV2Request) (rest.GetContactsV2Response, error) {
	uresp := rest.GetContactsV2Response{}
	if c == nil {
		return uresp, fmt.Errorf("nil client")
	}
	urlv, err := query.Values(req)
	if err != nil {
		return uresp, fmt.Errorf("query values: %w", err)
	}
	if err := c.get(ctx, fmt.Sprintf("/api/v2/contacts?%s", urlv.Encode()), &uresp); err != nil {
		return uresp, err
	}
	return uresp, nil
}

func (c *Client) CheckIntegration(ctx context.Context, req *rest.CheckIntegrationRequest) (*rest.CheckIntegrationResponse, error) {
	q := make(url.Values)
	if req.BranchID != "" {
		q.Set("branch_id", req.BranchID)
	}
	if req.ContactPhoneNumber != "" {
		q.Set("contact_phone_number", req.ContactPhoneNumber)
	}
	qenc := q.Encode()
	resp := &rest.CheckIntegrationResponse{}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/check-integration?%s", qenc), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetBusinesses(ctx context.Context, req *rest.GetBusinessesRequest) (*rest.GetBusinessesResponse, error) {
	q := make(url.Values)
	if req.StoreID != "" {
		q.Set("store_id", req.StoreID)
	}
	resp := &rest.GetBusinessesResponse{}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/business?%s", q.Encode()), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetPhones(ctx context.Context, req *rest.GetPhonesRequest) (*rest.GetPhonesResponse, error) {
	q := make(url.Values)
	if req.ID != 0 {
		q.Set("id", fmt.Sprint(req.ID))
	}
	if req.BranchID != "" {
		q.Set("branch_id", req.BranchID)
	}
	if req.BusinessID != 0 {
		q.Set("business_id", fmt.Sprint(req.BusinessID))
	}
	if req.PhoneNumber != "" {
		q.Set("phone_number", req.PhoneNumber)
	}
	if req.WithStatistics {
		q.Set("with_statistics", "true")
	}
	resp := &rest.GetPhonesResponse{}
	if err := c.get(ctx, fmt.Sprintf("/api/v1/phones?%s", q.Encode()), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) urlPrefix() string {
	return util.Default(c.BaseURL, DefaultBaseURL)
}

func (c *Client) get(ctx context.Context, suffix string, output any) error {
	return c.doRequest(ctx, http.MethodGet, suffix, nil, output)
}

func (c *Client) put(ctx context.Context, suffix string, input, output any) error {
	return c.doRequest(ctx, http.MethodPut, suffix, input, output)
}

func (c *Client) post(ctx context.Context, suffix string, input, output any) error {
	return c.doRequest(ctx, http.MethodPost, suffix, input, output)
}

func (c *Client) doRequest(ctx context.Context, method, suffix string, input, output any) error {
	var rdr io.Reader
	if input != nil {
		d, err := json.Marshal(input)
		if err != nil {
			return fmt.Errorf("marshal input: %w", err)
		}
		rdr = bytes.NewReader(d)
	}
	if c == nil {
		return fmt.Errorf("nil client inside doRequest")
	}
	// fmt.Println("will create request", method, c.urlPrefix()+suffix, rdr)
	req, err := http.NewRequest(method, c.urlPrefix()+suffix, rdr)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	if input != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.JWT != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.JWT))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return c.errorFromResponse(ctx, resp)
	}
	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (c *Client) errorFromResponse(ctx context.Context, resp *http.Response) error {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return fmt.Errorf("copy response body: %w", err)
	}
	jerr := &rest.ErrorResponse{}
	if err := json.Unmarshal(buf.Bytes(), jerr); err != nil {
		// not a valid json response!
		return &rest.ErrorResponse{
			StatusCode: resp.StatusCode,
			Raw:        buf.String(),
		}
	}
	jerr.StatusCode = resp.StatusCode
	return jerr
}
