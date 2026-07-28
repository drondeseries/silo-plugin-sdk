package runtime_test

import (
	"testing"

	pluginv1 "github.com/Silo-Server/silo-plugin-sdk/pkg/pluginproto/silo/plugin/v1"
	runtime "github.com/Silo-Server/silo-plugin-sdk/pkg/pluginsdk/runtime"
	"google.golang.org/grpc"
)

type stubVirtualStreamProvider struct {
	pluginv1.UnimplementedVirtualStreamProviderServer
}

func TestGRPCServerRegistersVirtualStreamProvider(t *testing.T) {
	p := &runtime.GRPCPlugin{Servers: runtime.CapabilityServers{
		Runtime:               stubRuntime{},
		VirtualStreamProvider: stubVirtualStreamProvider{},
	}}
	srv := grpc.NewServer()
	if err := p.GRPCServer(nil, srv); err != nil {
		t.Fatalf("GRPCServer with VirtualStreamProvider = %v, want nil", err)
	}
	if _, ok := srv.GetServiceInfo()["silo.plugin.v1.VirtualStreamProvider"]; !ok {
		t.Fatalf("VirtualStreamProvider service not registered; got %v", srv.GetServiceInfo())
	}
}

func TestClientVirtualStreamProvider(t *testing.T) {
	c := runtime.NewClient(nil)
	if client := c.VirtualStreamProvider(); client == nil {
		t.Fatal("VirtualStreamProvider() returned nil client")
	}
}
