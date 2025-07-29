module api_gateway

go 1.24

require (
	github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common v0.0.0
	github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils v0.0.0
	github.com/gorilla/mux v1.8.1
	github.com/sony/gobreaker/v2 v2.1.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redsync/redsync/v4 v4.13.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/redis/go-redis/v9 v9.11.0 // indirect
)

replace github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/utils => ../utils

replace github.com/MarcoFontana48/AUSL-Romagna-CCE-Microservices-Project-Proposal/common => ../common
