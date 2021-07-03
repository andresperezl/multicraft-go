package client

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func TestClient(t *testing.T) {

	testClient := New("http://example.com", "user", "key")
	httpmock.ActivateNonDefault(testClient.Client.GetClient())
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, "http://example.com",
		func(r *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK,
				`{"success":true,"errors":[],"data":{"status":"online","onlinePlayers":"0","maxPlayers":"100"}}`,
			)
			// The Multicraft API sets this header in the response, which messes with
			// the resty client parsing of the result
			resp.Header.Add("Content-Type", "text/html; charset=UTF-8")
			return resp, nil
		},
	)
	mr, err := testClient.Do("someMethod", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if !mr.Success {
		t.Error("expected response success to be true, but got false")
	}

	if mr.IsError() {
		t.Error("expected response to not be an error, but it was")
	}

	httpmock.RegisterResponder(http.MethodPost, "http://example.com",
		func(r *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK,
				`{"success":false,"errors":["Some Error"],"data":[]}`,
			)
			// The Multicraft API sets this header in the response, which messes with
			// the resty client parsing of the result
			resp.Header.Add("Content-Type", "text/html; charset=UTF-8")
			return resp, nil
		},
	)

	mr, err = testClient.Do("someMethod", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if mr.Success {
		t.Error("expected response success to be false, but got true")
	}

	if !mr.IsError() {
		t.Error("expected response to be an error, but it wasn't")
	}

	if got := mr.Error(); got != "Some Error" {
		t.Errorf("expected error message to be %s, but got %s", "Some Error", got)
	}
}
