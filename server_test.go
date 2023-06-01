package pricing_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/karlskewes/pricing"
)

func makeGetPriceHTTPRequest(req pricing.GetPriceRequest) (pricing.GetPriceResponse, error) {
	var got pricing.GetPriceResponse

	// TODO: let tests find and define port to listen on
	url := fmt.Sprintf("http://localhost:8080/api/v1/prices?brand_id=%d&product_id=%d&date=%s&string_id=%s", req.BrandID, req.ProductID, req.Date.Format(time.RFC3339), req.StringID)

	resp, err := http.Get(url)
	if err != nil {
		return got, fmt.Errorf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return got, fmt.Errorf("unexpected http status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		return got, fmt.Errorf("unexpected error decoding json response: %v", err)
	}

	return got, nil
}

// TestRun is an almost complete end to end test. Whilst it doesn't call main()
// it does execute the application with defaults, including HTTP listen address
// and repository configuration. Tests will fail if port 8080 isn't available
// to bind to.
// An alternative is to listen on a random port decided by httptest.NewServer
// per the tests in ./api_test.go
func TestRun(t *testing.T) {
	// t.Parallel() // insignficant speed boost and makes it easier to review
	app, err := pricing.NewApp([]string{"pricing"})
	if err != nil {
		log.Fatalf("failed to create new App: %v", err)
	}

	// nolint:errcheck // http request will fail if server doesn't startup
	go app.Run(context.Background())

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost:8080", nil)
	if err != nil {
		t.Fatalf("failed to create a new request: %v", err)
	}

	// HTTP server `go main()` goroutine might not be scheduled yet.
	// Attempt GET request a few times with a delay between each request.
	client := &http.Client{}
	var resp *http.Response
	var doErr error

	for i := 0; i < 3; i++ {
		resp, doErr = client.Do(req)
		if doErr == nil {
			defer resp.Body.Close()

			break
		}

		// wait for server to startup
		time.Sleep(time.Duration(i) * time.Second)
	}

	if doErr != nil {
		t.Fatalf("failed to query HTTP server")
	}

	want := http.StatusNotFound
	if resp.StatusCode != want {
		t.Errorf("want: %d - got: %d", want, resp.StatusCode)
	}

	/*
	   - Test 1: request at 10:00 on day 14 of the product 35455 for brand 1 (EXAMPLE).
	   - Test 2: request at 16:00 on the 14th for product 35455 for brand 1 (EXAMPLE)
	   - Test 3: request at 21:00 of the 14th of the day of the product 35455 for brand 1 (EXAMPLE)
	   - Test 4: request at 10:00 on the 15th of the day of the product 35455 for brand 1 (EXAMPLE)
	   - Test 5: request at 21:00 on the 16th of the day for product 35455 for brand 1 (EXAMPLE).
	*/

	testCases := map[string]struct {
		input   pricing.GetPriceRequest
		want    pricing.GetPriceResponse
		wantErr bool
	}{
		"Test 1": {
			input:   pricing.GetPriceRequest{1, 35455, time.Date(2020, 06, 14, 10, 0, 0, 0, time.UTC), "test_1"},
			want:    pricing.GetPriceResponse{1, 35455, "35.50", "EUR", "2020-06-14 00:00:00 +0000 UTC", "2020-12-31 23:59:59 +0000 UTC", "test_1"},
			wantErr: false,
		},
		"Test 2": {
			input:   pricing.GetPriceRequest{1, 35455, time.Date(2020, 06, 14, 16, 0, 0, 0, time.UTC), "test_2"},
			want:    pricing.GetPriceResponse{1, 35455, "25.45", "EUR", "2020-06-14 15:00:00 +0000 UTC", "2020-06-14 18:30:00 +0000 UTC", "test_2"},
			wantErr: false,
		},
		"Test 3": {
			input:   pricing.GetPriceRequest{1, 35455, time.Date(2020, 06, 14, 21, 0, 0, 0, time.UTC), "test_3"},
			want:    pricing.GetPriceResponse{1, 35455, "35.50", "EUR", "2020-06-14 00:00:00 +0000 UTC", "2020-12-31 23:59:59 +0000 UTC", "test_3"},
			wantErr: false,
		},
		"Test 4": {
			input:   pricing.GetPriceRequest{1, 35455, time.Date(2020, 06, 15, 10, 0, 0, 0, time.UTC), "test_4"},
			want:    pricing.GetPriceResponse{1, 35455, "30.50", "EUR", "2020-06-15 00:00:00 +0000 UTC", "2020-06-15 11:00:00 +0000 UTC", "test_4"},
			wantErr: false,
		},
		"Test 5": {
			input:   pricing.GetPriceRequest{1, 35455, time.Date(2020, 06, 16, 21, 0, 0, 0, time.UTC), "test_5"},
			want:    pricing.GetPriceResponse{1, 35455, "38.95", "EUR", "2020-06-15 16:00:00 +0000 UTC", "2020-12-31 23:59:59 +0000 UTC", "test_5"},
			wantErr: false,
		},
	}

	for testName, tc := range testCases {
		tt := tc
		t.Run(testName, func(t *testing.T) {
			// t.Parallel() // insignificant speed boost and jumbles order which
			// makes it harder to compare.
			// No race conditions seem to occur when enabled.

			got, err := makeGetPriceHTTPRequest(tt.input)
			if err != nil {
				t.Errorf("unexpected error calling API endpoint: %v", err)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("db.GetPrice(...) mismatch (-want +got):\n%s", diff)
			}
			t.Logf("%s - %v", tt.want.StringID, got)
		})
	}
}
