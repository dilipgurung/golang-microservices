package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	context "golang.org/x/net/context"

	"github.com/dilipgurung/golang-microservices/pb/rate"
	"google.golang.org/grpc"
)

const baseCurrency = "USD"

type rateServer struct {
	rateTable map[string]float64
}

func (s *rateServer) GetRates(ctx context.Context, r *rate.Request) (*rate.Result, error) {
	result := new(rate.Result)

	curSource, ok := s.rateTable[r.SourceCurrency]
	if !ok {
		return result, fmt.Errorf("The given currency %s was not found", r.SourceCurrency)
	}

	curTarget, ok := s.rateTable[r.TargetCurrency]
	if !ok {
		return result, fmt.Errorf("The given currency %s was not found", r.TargetCurrency)
	}

	// if the source currency is USD then we can return the
	// base conversion rate for the target currency straight away
	if r.SourceCurrency == baseCurrency {
		result.Rate = curTarget
		return result, nil
	}

	result.Rate = curTarget / curSource

	return result, nil
}

func (s *rateServer) loadRates(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &s.rateTable)
}

var port = flag.Int("port", 8080, "The server port")

func main() {
	flag.Parse()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// crete the grpc server
	srv := grpc.NewServer()
	rs := new(rateServer)
	err = rs.loadRates("./rates.json")
	if err != nil {
		log.Fatalf("Failed to load rates data: %v", err)
	}
	rate.RegisterRateServer(srv, rs)
	srv.Serve(listen)
}
