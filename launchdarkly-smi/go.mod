module example.com/intheclouddan/launchdarkly-smi/v2

go 1.14
replace github.com/intheclouddan/launchdarkly-smi => /Users/danielobrien/Projects/ld-rate-limiter/launchdarkly-smi

require (
	github.com/intheclouddan/launchdarkly-smi v1.0.0
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/envoyproxy/go-control-plane v0.9.5
	github.com/google/uuid v1.1.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/launchdarkly/eventsource v1.4.1 // indirect
	github.com/lyft/gostats v0.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.3
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd
	google.golang.org/grpc v1.29.0
	gopkg.in/launchdarkly/go-sdk-common.v1 v1.0.0-20200401173443-991b2f427a01 // indirect
	gopkg.in/launchdarkly/go-server-sdk.v4 v4.0.0-20200416175003-6f2d5c743567
)
