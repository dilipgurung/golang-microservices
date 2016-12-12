package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dilipgurung/golang-microservices/pb/rate"
	context "golang.org/x/net/context"

	"fmt"

	"google.golang.org/grpc"
)

const (
	rateUsdToGbp = 0.79435
	rateGbpToUsd = 1.25889
	usd          = "USD"
	gbp          = "GBP"
)

// response maps to the json response from the rates api service
type response struct {
	Rate float64 `json:"rate"`
}

func TestHealthCheck(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(healthCheckHandler))

	c := &http.Client{}
	res, err := c.Get(ts.URL + "/healthcheck")
	if err != nil {
		t.Fatalf("Expected error to be nil but got: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected response status to be: %d but got: %d", http.StatusOK, res.StatusCode)
	}

}

func TestGetCurrencies(t *testing.T) {
	curList = currencyList{}
	curList.Currencies = []string{"USD", "GBP"}

	ts := httptest.NewServer(http.HandlerFunc(getCurrenciesHandler))

	b, err := makeRequestTo(ts.URL + "/api/currencies")
	if err != nil {
		t.Fatal(err)
	}

	resCurList := &currencyList{}
	err = json.Unmarshal(b, resCurList)
	if err != nil {
		t.Fatalf("Expected error to be nil but got: %v", err)
	}

	if resCurList.Currencies[0] != usd {
		t.Fatalf("Expected result to be: %s but got: %s", usd, resCurList.Currencies[0])
	}
}

// mockRateClient is an implementation of rate.RateClient
type mockRateClient struct{}

func (c *mockRateClient) GetRates(ctx context.Context, in *rate.Request, opts ...grpc.CallOption) (*rate.Result, error) {
	r := &rate.Result{Rate: rateUsdToGbp}

	if in.GetSourceCurrency() == gbp {
		r.Rate = rateGbpToUsd
	}
	return r, nil
}

func TestGetRate(t *testing.T) {
	resp := new(response)
	rc := new(mockRateClient)
	ts := httptest.NewServer(getRateHandler(rc))

	for _, test := range []struct {
		sourceCurrency string
		targetCurrency string
		expectedRate   float64
	}{
		{usd, gbp, rateUsdToGbp},
		{gbp, usd, rateGbpToUsd},
	} {
		url := fmt.Sprintf(ts.URL+"/api/rates?calculationType=source&sourceCurrency=%s&targetCurrency=%s", test.sourceCurrency, test.targetCurrency)
		b, err := makeRequestTo(url)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(b, resp)
		if err != nil {
			t.Fatal(err)
		}

		if resp.Rate != test.expectedRate {
			t.Fatalf("Expected rate to be: %f but got: %f", test.expectedRate, resp.Rate)
		}

	}
}

// helper function to make a request to the given url
// and return the response body as []byte
func makeRequestTo(url string) ([]byte, error) {
	c := &http.Client{}
	res, err := c.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Expected error to be nil but got: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected response status to be: %d but got: %d", http.StatusOK, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Expected error to be nil but got: %v", err)
	}

	return b, nil
}
