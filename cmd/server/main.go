package main

import (
	"log"
	"net"
	"time"

	rls "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v2"
	stats "github.com/lyft/gostats"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	ld "gopkg.in/launchdarkly/go-server-sdk.v4"
)

// server is used to implement rls.RateLimitService
type server struct {
	// limit specifies if the next request is to be rate limited
	limit    bool
	ldClient *ld.LDClient
	store    stats.Store
	scope    stats.Scope
}

func (s *server) ShouldRateLimit(ctx context.Context,
	request *rls.RateLimitRequest) (*rls.RateLimitResponse, error) {
	log.Printf("request: %v\n", request)

	// logic to rate limit every second request
	var overallCode rls.RateLimitResponse_Code
	if s.limit {
		overallCode = rls.RateLimitResponse_OVER_LIMIT
		s.limit = false
	} else {
		overallCode = rls.RateLimitResponse_OK
		s.limit = true
	}

	response := &rls.RateLimitResponse{OverallCode: overallCode}
	log.Printf("response: %v\n", response)
	return response, nil
}

func main() {
	// create a TCP listener on port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s", lis.Addr())

	// create a gRPC server and register the RateLimitService server
	ldClient, _ := ld.MakeClient("sdk-3617fc9c-5e8b-47e3-9a07-ca04d23fa4fa", 5*time.Second)
	s := grpc.NewServer()
	// setup stats
	store := stats.NewDefaultStore()
	rls.RegisterRateLimitServiceServer(s, &server{limit: false, ldClient: ldClient, store: store, scope: store.Scope("rate_limit")})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
