package pluginv1_test

import (
	"testing"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pluginv1 "github.com/Silo-Server/silo-plugin-sdk/pkg/pluginproto/silo/plugin/v1"
)

func TestVirtualStreamContract(t *testing.T) {
	now := time.Now().Truncate(time.Second).UTC()
	expiresAt := now.Add(2 * time.Hour)
	readyAt := now.Add(5 * time.Second)

	extraMeta, err := structpb.NewStruct(map[string]any{
		"cdn_node": "us-east-1",
		"tier":     "premium",
	})
	if err != nil {
		t.Fatalf("failed to create structpb: %v", err)
	}

	cand1 := &pluginv1.VirtualStreamCandidate{
		CandidateId:  "cand-1080p-hevc",
		ProviderId:   "plugin.usenet.provider",
		TemporaryUri: "https://stream.internal/v1/play/stream-1080p.mkv?token=xyz123",
		ExpiresAt:    timestamppb.New(expiresAt),
		Rank:         1,
		Resolution: &pluginv1.VirtualStreamResolution{
			Width:     1920,
			Height:    1080,
			Label:     "1080p",
			FrameRate: 23.976,
		},
		VideoCodec: "hevc",
		AudioCodec: "eac3",
		Hdr: &pluginv1.VirtualStreamHDR{
			IsHdr:               true,
			Format:              "HDR10",
			HasDolbyVision:      false,
			DolbyVisionProfile: "",
		},
		Bitrate:        8000000,
		FileSizeBytes:  10737418240, // 10 GB
		Container:      "mkv",
		AudioLanguages:    []string{"eng", "fre"},
		SubtitleLanguages: []string{"eng", "spa", "fre"},
		Availability: &pluginv1.VirtualStreamAvailability{
			State:            pluginv1.VirtualStreamAvailabilityState_VIRTUAL_STREAM_AVAILABILITY_STATE_AVAILABLE,
			Message:          "Stream ready for playback",
			EstimatedReadyAt: timestamppb.New(readyAt),
			ProgressPercent:  100,
		},
		Error: nil,
		Metadata: extraMeta,
	}

	cand2 := &pluginv1.VirtualStreamCandidate{
		CandidateId:  "cand-4k-dv",
		ProviderId:   "plugin.usenet.provider",
		TemporaryUri: "https://stream.internal/v1/play/stream-4k.mkv?token=abc456",
		ExpiresAt:    timestamppb.New(expiresAt),
		Rank:         2,
		Resolution: &pluginv1.VirtualStreamResolution{
			Width:     3840,
			Height:    2160,
			Label:     "2160p",
			FrameRate: 23.976,
		},
		VideoCodec: "hevc",
		AudioCodec: "truehd",
		Hdr: &pluginv1.VirtualStreamHDR{
			IsHdr:               true,
			Format:              "Dolby Vision",
			HasDolbyVision:      true,
			DolbyVisionProfile: "Profile 8.1",
		},
		Bitrate:        25000000,
		FileSizeBytes:  32212254720, // 30 GB
		Container:      "mkv",
		AudioLanguages:    []string{"eng"},
		SubtitleLanguages: []string{"eng"},
		Availability: &pluginv1.VirtualStreamAvailability{
			State:            pluginv1.VirtualStreamAvailabilityState_VIRTUAL_STREAM_AVAILABILITY_STATE_AVAILABLE,
			Message:          "Stream ready for playback",
			EstimatedReadyAt: timestamppb.New(readyAt),
			ProgressPercent:  100,
		},
	}

	result := &pluginv1.VirtualStreamResult{
		ResultId:   "res-9999",
		ProviderId: "plugin.usenet.provider",
		Candidates: []*pluginv1.VirtualStreamCandidate{cand1, cand2},
		Availability: &pluginv1.VirtualStreamAvailability{
			State:           pluginv1.VirtualStreamAvailabilityState_VIRTUAL_STREAM_AVAILABILITY_STATE_AVAILABLE,
			Message:         "2 stream candidates resolved",
			ProgressPercent: 100,
		},
		Error: nil,
	}

	// Verify fields on generated types
	if result.GetResultId() != "res-9999" {
		t.Errorf("GetResultId() = %q, want %q", result.GetResultId(), "res-9999")
	}
	if result.GetProviderId() != "plugin.usenet.provider" {
		t.Errorf("GetProviderId() = %q, want %q", result.GetProviderId(), "plugin.usenet.provider")
	}
	if len(result.GetCandidates()) != 2 {
		t.Fatalf("len(Candidates) = %d, want 2", len(result.GetCandidates()))
	}

	c1 := result.GetCandidates()[0]
	if c1.GetCandidateId() != "cand-1080p-hevc" {
		t.Errorf("c1 candidate_id = %q, want %q", c1.GetCandidateId(), "cand-1080p-hevc")
	}
	if c1.GetTemporaryUri() != "https://stream.internal/v1/play/stream-1080p.mkv?token=xyz123" {
		t.Errorf("c1 temporary_uri = %q", c1.GetTemporaryUri())
	}
	if c1.GetRank() != 1 {
		t.Errorf("c1 rank = %d, want 1", c1.GetRank())
	}
	if c1.GetResolution().GetHeight() != 1080 || c1.GetResolution().GetLabel() != "1080p" {
		t.Errorf("c1 resolution = %v", c1.GetResolution())
	}
	if c1.GetVideoCodec() != "hevc" || c1.GetAudioCodec() != "eac3" {
		t.Errorf("c1 codecs = %s/%s", c1.GetVideoCodec(), c1.GetAudioCodec())
	}
	if !c1.GetHdr().GetIsHdr() || c1.GetHdr().GetFormat() != "HDR10" || c1.GetHdr().GetHasDolbyVision() {
		t.Errorf("c1 hdr = %v", c1.GetHdr())
	}
	if c1.GetBitrate() != 8000000 || c1.GetFileSizeBytes() != 10737418240 {
		t.Errorf("c1 bitrate/size = %d / %d", c1.GetBitrate(), c1.GetFileSizeBytes())
	}
	if c1.GetContainer() != "mkv" {
		t.Errorf("c1 container = %q", c1.GetContainer())
	}
	if len(c1.GetAudioLanguages()) != 2 || len(c1.GetSubtitleLanguages()) != 3 {
		t.Errorf("c1 languages = %v / %v", c1.GetAudioLanguages(), c1.GetSubtitleLanguages())
	}

	c2 := result.GetCandidates()[1]
	if !c2.GetHdr().GetHasDolbyVision() || c2.GetHdr().GetDolbyVisionProfile() != "Profile 8.1" {
		t.Errorf("c2 dolby vision = %v", c2.GetHdr())
	}

	// Test proto binary roundtrip
	binaryData, err := proto.Marshal(result)
	if err != nil {
		t.Fatalf("proto.Marshal error: %v", err)
	}
	var unmarshaledResult pluginv1.VirtualStreamResult
	if err := proto.Unmarshal(binaryData, &unmarshaledResult); err != nil {
		t.Fatalf("proto.Unmarshal error: %v", err)
	}
	if len(unmarshaledResult.GetCandidates()) != 2 {
		t.Fatalf("unmarshaled len(Candidates) = %d, want 2", len(unmarshaledResult.GetCandidates()))
	}
	if unmarshaledResult.GetCandidates()[0].GetCandidateId() != "cand-1080p-hevc" {
		t.Errorf("unmarshaled c1 candidate_id = %q", unmarshaledResult.GetCandidates()[0].GetCandidateId())
	}

	// Test protojson roundtrip
	jsonData, err := protojson.Marshal(result)
	if err != nil {
		t.Fatalf("protojson.Marshal error: %v", err)
	}
	var jsonResult pluginv1.VirtualStreamResult
	if err := protojson.Unmarshal(jsonData, &jsonResult); err != nil {
		t.Fatalf("protojson.Unmarshal error: %v", err)
	}
	if jsonResult.GetResultId() != "res-9999" {
		t.Errorf("protojson result_id = %q, want %q", jsonResult.GetResultId(), "res-9999")
	}
}

func TestVirtualStreamErrorAndAvailability(t *testing.T) {
	errMeta := &pluginv1.VirtualStreamError{
		Code:       pluginv1.VirtualStreamErrorCode_VIRTUAL_STREAM_ERROR_CODE_TEMPORARY,
		Message:    "Upstream provider connection rate limited",
		Retryable:  true,
		RetryAfter: durationpb.New(30 * time.Second),
	}

	availMeta := &pluginv1.VirtualStreamAvailability{
		State:           pluginv1.VirtualStreamAvailabilityState_VIRTUAL_STREAM_AVAILABILITY_STATE_DEGRADED,
		Message:         "High latency stream path",
		ProgressPercent: 50,
	}

	result := &pluginv1.VirtualStreamResult{
		ResultId:     "res-error-test",
		ProviderId:   "plugin.torrent.provider",
		Candidates:   nil,
		Availability: availMeta,
		Error:        errMeta,
	}

	if result.GetError().GetCode() != pluginv1.VirtualStreamErrorCode_VIRTUAL_STREAM_ERROR_CODE_TEMPORARY {
		t.Errorf("error code = %v, want TEMPORARY", result.GetError().GetCode())
	}
	if !result.GetError().GetRetryable() {
		t.Errorf("error retryable = false, want true")
	}
	if result.GetError().GetRetryAfter().AsDuration() != 30*time.Second {
		t.Errorf("retry after = %v, want 30s", result.GetError().GetRetryAfter().AsDuration())
	}
	if result.GetAvailability().GetState() != pluginv1.VirtualStreamAvailabilityState_VIRTUAL_STREAM_AVAILABILITY_STATE_DEGRADED {
		t.Errorf("availability state = %v, want DEGRADED", result.GetAvailability().GetState())
	}
}

func TestResolveVirtualStreamRPCMessages(t *testing.T) {
	req := &pluginv1.ResolveVirtualStreamRequest{
		CapabilityId: "virtual-stream-1",
		MediaType:    "movie",
		Title:        "Inception",
		Year:         2010,
		ExternalIds: map[string]string{
			"tmdb": "27205",
			"imdb": "tt1375666",
		},
	}

	resp := &pluginv1.ResolveVirtualStreamResponse{
		Result: &pluginv1.VirtualStreamResult{
			ResultId:   "res-101",
			ProviderId: "plugin.virtual.provider",
			Candidates: []*pluginv1.VirtualStreamCandidate{
				{
					CandidateId:  "cand-1",
					TemporaryUri: "https://example.com/stream.mkv",
					Rank:         1,
				},
			},
		},
	}

	if req.GetTitle() != "Inception" || req.GetYear() != 2010 {
		t.Errorf("req = %v", req)
	}
	if len(resp.GetResult().GetCandidates()) != 1 {
		t.Errorf("resp candidates len = %d, want 1", len(resp.GetResult().GetCandidates()))
	}
}
