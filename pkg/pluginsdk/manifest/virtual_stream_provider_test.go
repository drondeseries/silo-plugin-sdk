package manifest_test

import (
	"testing"

	pluginv1 "github.com/Silo-Server/silo-plugin-sdk/pkg/pluginproto/silo/plugin/v1"
	publicmanifest "github.com/Silo-Server/silo-plugin-sdk/pkg/pluginsdk/manifest"
)

func TestLoadVirtualStreamProvider(t *testing.T) {
	raw := []byte(`{
	  "plugin_id":"silo.virtual.stream", "version":"1.0.0", "silo_api_version":"v1",
	  "capabilities":[{
	    "type":"virtual_stream_provider.v1", "id":"stream_provider_1", "display_name":"Virtual Stream Provider",
	    "virtual_stream_provider":{
	      "supported_media_types":["movie", "episode"],
	      "supported_containers":["mkv", "mp4"],
	      "supports_just_in_time_resolution":true,
	      "supports_multiple_candidates":true
	    }
	  }]
	}`)
	manifest, err := publicmanifest.Load(raw)
	if err != nil {
		t.Fatalf("failed to load virtual stream manifest: %v", err)
	}
	caps := manifest.GetCapabilities()
	if len(caps) != 1 {
		t.Fatalf("capabilities count = %d, want 1", len(caps))
	}
	desc := caps[0].GetVirtualStreamProvider()
	if desc == nil {
		t.Fatalf("expected non-nil VirtualStreamProvider descriptor")
	}
	if !desc.GetSupportsJustInTimeResolution() || !desc.GetSupportsMultipleCandidates() {
		t.Errorf("descriptor flags = %v", desc)
	}
	if len(desc.GetSupportedMediaTypes()) != 2 || desc.GetSupportedMediaTypes()[0] != "movie" {
		t.Errorf("supported media types = %v", desc.GetSupportedMediaTypes())
	}
}

func TestValidateVirtualStreamProviderRejectsDescriptorOnOtherCapability(t *testing.T) {
	manifest := &pluginv1.PluginManifest{
		PluginId: "silo.invalid", Version: "1.0.0",
		Capabilities: []*pluginv1.CapabilityDescriptor{{
			Type: "event_consumer.v1", Id: "events",
			VirtualStreamProvider: &pluginv1.VirtualStreamProviderDescriptor{
				SupportsJustInTimeResolution: true,
			},
		}},
	}
	if err := publicmanifest.Validate(manifest); err == nil {
		t.Fatal("expected misplaced virtual stream descriptor to fail validation")
	}
}

func TestValidateVirtualStreamProviderRequiresTypedDescriptor(t *testing.T) {
	manifest := &pluginv1.PluginManifest{
		PluginId: "silo.invalid", Version: "1.0.0",
		Capabilities: []*pluginv1.CapabilityDescriptor{{
			Type: "virtual_stream_provider.v1", Id: "streams",
		}},
	}
	if err := publicmanifest.Validate(manifest); err == nil {
		t.Fatal("expected missing virtual stream descriptor to fail validation")
	}
}

func TestValidateVirtualStreamProviderRequiresKnownMediaTypes(t *testing.T) {
	manifest := &pluginv1.PluginManifest{
		PluginId: "silo.invalid", Version: "1.0.0",
		Capabilities: []*pluginv1.CapabilityDescriptor{{
			Type: "virtual_stream_provider.v1", Id: "streams",
			VirtualStreamProvider: &pluginv1.VirtualStreamProviderDescriptor{},
		}},
	}
	if err := publicmanifest.Validate(manifest); err == nil {
		t.Fatal("expected empty virtual stream media types to fail validation")
	}
	manifest.Capabilities[0].VirtualStreamProvider.SupportedMediaTypes = []string{"movie", "banana"}
	if err := publicmanifest.Validate(manifest); err == nil {
		t.Fatal("expected unknown virtual stream media type to fail validation")
	}
}
