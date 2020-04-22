/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"net"
	"time"

	rls "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v2"
	stats "github.com/lyft/gostats"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	ld "gopkg.in/launchdarkly/go-server-sdk.v4"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

var port string
var ldApiKey string

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "port to listen on")
	viper.BindPFlag("port", serverCmd.PersistentFlags().Lookup("port"))
	viper.SetDefault("port", ":50052")

	serverCmd.PersistentFlags().StringVarP(&ldApiKey, "apiKey", "", "", "Launchdarkly API Key")
	viper.BindPFlag("apiKey", serverCmd.PersistentFlags().Lookup("apiKey"))
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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

func startServer() {
	listenAddr := viper.GetString("port")
	ldApiKey := viper.GetString("apiKey")
	// create a TCP listener on port 50052
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on %s", lis.Addr())

	// create a gRPC server and register the RateLimitService server
	ldClient, _ := ld.MakeClient(ldApiKey, 5*time.Second)
	s := grpc.NewServer()
	// setup stats
	store := stats.NewDefaultStore()
	rls.RegisterRateLimitServiceServer(s, &server{limit: false, ldClient: ldClient, store: store, scope: store.Scope("rate_limit")})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
