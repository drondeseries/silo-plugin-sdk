package convert

import (
	"testing"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"

	pluginv1 "github.com/Silo-Server/silo-plugin-sdk/pkg/pluginproto/silo/plugin/v1"
)

// TestOldVirtualDescriptorWireFieldIsTolerated pins the fork's 11->13
// renumber of virtual_stream_provider: plugins built before the move emit
// the descriptor at field 11, which the new host must drop as unknown
// instead of rejecting the capability or misreading the request_router
// descriptor at field 12.
func TestOldVirtualDescriptorWireFieldIsTolerated(t *testing.T) {
	oldVirtual, err := proto.Marshal(&pluginv1.VirtualStreamProviderDescriptor{})
	if err != nil {
		t.Fatalf("marshal old virtual descriptor: %v", err)
	}
	router, err := proto.Marshal(&pluginv1.RequestRouterDescriptor{SupportsSeasons: true})
	if err != nil {
		t.Fatalf("marshal request router descriptor: %v", err)
	}
	var raw []byte
	raw = protowire.AppendTag(raw, 11, protowire.BytesType)
	raw = protowire.AppendBytes(raw, oldVirtual)
	raw = protowire.AppendTag(raw, 12, protowire.BytesType)
	raw = protowire.AppendBytes(raw, router)

	var got pluginv1.CapabilityDescriptor
	if err := proto.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal capability with old field 11: %v", err)
	}
	if !got.GetRequestRouter().GetSupportsSeasons() {
		t.Fatalf("request_router.supports_seasons lost: %+v", got.GetRequestRouter())
	}
	if got.GetVirtualStreamProvider() != nil {
		t.Fatalf("old field-11 bytes must decode as unknown, got %+v", got.GetVirtualStreamProvider())
	}
}
