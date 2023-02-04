package main

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const exampleScheme = "etcd"

// ExampleResolverBuilder builds a custom resolver.
type ExampleResolverBuilder struct{}

// Build returns a custom resolver.
func (b *ExampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return &exampleResolver{target: target, cc: cc}, nil
}

// Scheme returns the scheme supported by this resolver builder.
func (b *ExampleResolverBuilder) Scheme() string {
	return exampleScheme
}

// exampleResolver is a custom resolver.
type exampleResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

// ResolveNow is called by gRPC to resolve the target.
func (r *exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {
	// Implement your custom resolve logic here.
	addrs := []resolver.Address{{Addr: "localhost:1234"}}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}

func (r *exampleResolver) Close() {
	// Clean up any resources you created in Build().
}

func main() {
	// Register the custom resolver builder.
	resolver.Register(&ExampleResolverBuilder{})

	// Use the custom resolver when dialing a gRPC server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:///localhost:1234", exampleScheme),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithResolvers(&ExampleResolverBuilder{}))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// ... use the connection to call gRPC services ...
}
