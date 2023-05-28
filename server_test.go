package pricing_test

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/karlskewes/pricing"
)

func TestRun(t *testing.T) {
	app, err := pricing.NewApp(nil)
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
}
