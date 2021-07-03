package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

// MulticraftAPIClient holds the information, and underlying HTTP client to talk
// to a Multicraft API
type MulticraftAPIClient struct {
	Client *resty.Client
	user   string
	apiKey string
}

// New creates a new *MulticraftAPIClient, which uses the specified url, user,
// and apiKey
func New(url, user, apiKey string) *MulticraftAPIClient {
	return &MulticraftAPIClient{
		Client: resty.New().
			SetHostURL(url).
			SetHeaders(map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			}),
		user:   user,
		apiKey: apiKey,
	}
}

// Do performs a requests against the API client, attaching the necessary data
// and credentials in the expected format
func (mc *MulticraftAPIClient) Do(method string, params map[string]string) (
	*MulticraftResponse, error) {

	hmacSign := &bytes.Buffer{}
	content := &bytes.Buffer{}

	// Add the requests parameters first
	for k, v := range params {
		if err := addKeyValue(hmacSign, content, k, v); err != nil {
			return nil, err
		}
	}

	// Add the required method name and user
	if err := addKeyValue(hmacSign, content, "_MulticraftAPIMethod", method); err != nil {
		return nil, err
	}

	if err := addKeyValue(hmacSign, content, "_MulticraftAPIUser", mc.user); err != nil {
		return nil, err
	}

	// Create the APIKey hash with the previous content and add to the body
	signer := hmac.New(sha256.New, []byte(mc.apiKey))
	contentHMAC := hmacSign.Bytes()

	if _, err := signer.Write(contentHMAC); err != nil {
		return nil, err
	}

	signature := hex.EncodeToString(signer.Sum(nil))
	if _, err := content.WriteRune('&'); err != nil {
		return nil, err
	}

	if _, err := content.WriteString(fmt.Sprintf(
		"%s=%s", url.QueryEscape("_MulticraftAPIKey"), url.QueryEscape(signature))); err != nil {
		return nil, err
	}

	resp := &MulticraftResponse{}
	errResp := &MulticraftResponse{}
	restyResp, err := mc.Client.R().
		SetResult(resp).
		SetError(errResp).
		SetBody(content.String()).
		Post("")
	if err != nil {
		return nil, err
	}
	if restyResp.IsError() {
		return errResp, nil
	}
	return resp, nil
}

type MulticraftResponse struct {
	Status bool        `json:"status"`
	Errors []string    `json:"errors"`
	Data   interface{} `json:"data"`
}

func addKeyValue(hmacContent, bodyContent *bytes.Buffer, key, value string) error {
	if _, err := hmacContent.WriteString(key); err != nil {
		return err
	}

	if _, err := hmacContent.WriteString(value); err != nil {
		return err
	}
	if bodyContent.Len() > 0 {
		if _, err := bodyContent.WriteRune('&'); err != nil {
			return err
		}
	}
	if _, err := bodyContent.WriteString(fmt.Sprintf(
		"%s=%s", url.QueryEscape(key), url.QueryEscape(value))); err != nil {
		return err
	}
	return nil
}
