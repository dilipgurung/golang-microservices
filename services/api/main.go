package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	context "golang.org/x/net/context"

	"github.com/dilipgurung/golang-microservices/pb/rate"
	"google.golang.org/grpc"
)

var (
	port        = flag.Int("port", 9000, "The server port")
	rateSvcAddr = flag.String("rate", "rate:8080", "The Rate service address")
	curList     currencyList
)

func main() {
	flag.Parse()

	// load the currency list to the memory
	curList = getCurList()

	// make a connection to the rate server
	// And create a rate client
	conn, err := grpc.Dial(*rateSvcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	rc := rate.NewRateClient(conn)

	http.HandleFunc("/api/currencies", getCurrenciesHandler)
	http.HandleFunc("/healthcheck", healthCheckHandler)
	http.Handle("/api/rates", getRateHandler(rc))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func healthCheckHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

func getCurrenciesHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")

	b, err := json.Marshal(curList)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(b)
}

func getRateHandler(rc rate.RateClient) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		req := new(rate.Request)

		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Content-Type", "application/json")

		cType, sourceCur, targetCur := r.URL.Query().Get("calculationType"), r.URL.Query().Get("sourceCurrency"), r.URL.Query().Get("targetCurrency")
		if cType == "" || sourceCur == "" || targetCur == "" {
			http.Error(rw, "Please specify calculationType/sourceCurrency/targetCurrency params", http.StatusBadRequest)
			return
		}

		switch cType {
		case "source":
			req.SourceCurrency = sourceCur
			req.TargetCurrency = targetCur
		case "target":
			req.SourceCurrency = targetCur
			req.TargetCurrency = sourceCur
		default:
			http.Error(rw, "Unknown calculationType", http.StatusBadRequest)
			return
		}

		rates, err := rc.GetRates(context.Background(), req)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(rates)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(b)
	})
}
