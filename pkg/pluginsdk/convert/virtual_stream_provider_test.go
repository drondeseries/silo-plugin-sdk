package convert_test

import (
	"testing"

	pluginv1 "github.com/Silo-Server/silo-plugin-sdk/pkg/pluginproto/silo/plugin/v1"
	"github.com/Silo-Server/silo-plugin-sdk/pkg/pluginsdk/capability"
	"github.com/Silo-Server/silo-plugin-sdk/pkg/pluginsdk/convert"
)

func TestVirtualStreamProviderConvertRoundtrip(t *testing.T) {
	manifest := &pluginv1.PluginManifest{
		Capabilities: []*pluginv1.CapabilityDescriptor{
			{
				Type:        capability.VirtualStreamProvider,
				Id:          "stream_convert_1",
				DisplayName: "Stream Provider Convert Test",
				VirtualStreamProvider: &pluginv1.VirtualStreamProviderDescriptor{
					SupportedMediaTypes:          []string{"movie", "episode"},
					SupportedContainers:          []string{"mkv", "mp4"},
					SupportsJustInTimeResolution: true,
					SupportsMultipleCandidates:  true,
				},
			},
		},
	}

	records, err := convert.CapabilityRecordsFromManifest(manifest)
	if err != nil {
		t.Fatalf("CapabilityRecordsFromManifest error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("records count = %d, want 1", len(records))
	}

	decoded, err := convert.DecodeCapability(records[0])
	if err != nil {
		t.Fatalf("DecodeCapability error: %v", err)
	}
	if decoded.GetType() != capability.VirtualStreamProvider {
		t.Errorf("decoded type = %q, want %q", decoded.GetType(), capability.VirtualStreamProvider)
	}
	vdesc := decoded.GetVirtualStreamProvider()
	if vdesc == nil {
		t.Fatalf("expected non-nil VirtualStreamProvider descriptor")
	}
	if !vdesc.GetSupportsJustInTimeResolution() || !vdesc.GetSupportsMultipleCandidates() {
		t.Errorf("vdesc flags = %v", vdesc)
	}
	if len(vdesc.GetSupportedContainers()) != 2 || vdesc.GetSupportedContainers()[0] != "mkv" {
		t.Errorf("supported containers = %v", vdesc.GetSupportedContainers())
	}
}
