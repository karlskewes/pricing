package pricing_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/karlskewes/pricing"
)

func TestAPIGetBrand(t *testing.T) {
	t.Parallel()
	// TODO table tests with non match
	repo := pricing.NewMockRepository()
	svc := pricing.NewService(repo)
	h, err := pricing.NewHandler(svc)
	if err != nil {
		t.Fatalf("unable to create handler with mock repository: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(h.GetBrand))

	t.Cleanup(func() {
		ts.Close()
	})

	// per MockRepository in ./repository.go
	want := pricing.Brand{
		ID:   1234,
		Name: "EXAMPLE",
	}

	url := fmt.Sprintf("%s/api/v1/brands?name=EXAMPLE", ts.URL)

	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected http status: %d", resp.StatusCode)
	}

	var got pricing.Brand

	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Errorf("unexpected error decoding json response: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("db.GetPrice(...) mismatch (-want +got):\n%s", diff)
	}

	resp.Body.Close()
}

func TestAPIGetPrice(t *testing.T) {
	t.Parallel()
	// TODO table tests with non match
	repo := pricing.NewMockRepository()
	svc := pricing.NewService(repo)
	h, err := pricing.NewHandler(svc)
	if err != nil {
		t.Fatalf("unable to create handler with mock repository: %v", err)
	}

	ts := httptest.NewServer(http.HandlerFunc(h.GetPrice))

	t.Cleanup(func() {
		ts.Close()
	})

	input := pricing.GetPriceRequest{
		BrandID:   1,
		ProductID: 1234,
		Date:      time.Date(2020, 06, 14, 10, 0, 0, 0, time.UTC),
		StringID:  "test_1",
	}
	// per MockRepository in ./repository.go
	want := pricing.GetPriceResponse{
		BrandID:   input.BrandID,
		ProductID: input.ProductID,
		StartDate: "2020-06-14 10:00:00 +0000 UTC",
		EndDate:   "2020-06-15 10:00:00 +0000 UTC",
		Price:     "1.00",
		Curr:      "USD",
		StringID:  input.StringID,
	}

	url := fmt.Sprintf("%s/api/v1/prices?brand_id=%d&product_id=%d&date=%s&string_id=%s", ts.URL, input.BrandID, input.ProductID, input.Date.Format(time.RFC3339), input.StringID)

	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected http status: %d", resp.StatusCode)
	}

	var got pricing.GetPriceResponse

	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Errorf("unexpected error decoding json response: %v", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("db.GetPrice(...) mismatch (-want +got):\n%s", diff)
	}

	resp.Body.Close()
}
