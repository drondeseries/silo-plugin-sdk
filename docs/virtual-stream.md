# Virtual Stream Provider Contract (`virtual_stream_provider.v1`)

`virtual_stream_provider.v1` defines an additive, provider-neutral result contract for just-in-time virtual playback stream resolution. Plugins implementing this capability act as virtual stream resolvers, returning structured stream candidates for playback requests.

## Overview

When the Silo host or an integration requests a playback stream, the virtual stream provider resolves matching candidate streams on demand. The response carries candidate ranking, temporary URIs, expiration timestamps, detailed stream tech specs (codecs, resolution, HDR/Dolby Vision, bitrate, container), multi-track audio and subtitle languages, along with availability and error metadata.

## Contract Architecture

The protobuf contract is defined in `proto/silo/plugin/v1/virtual_stream.proto` and generated under `pkg/pluginproto/silo/plugin/v1`.

### Key Messages

#### `VirtualStreamCandidate`

Represents a single stream candidate suitable for playback:

- `candidate_id` (`string`): Unique identifier for the stream candidate.
- `provider_id` (`string`): Identifier of the stream provider or source plugin (e.g. `plugin.debrid`, `plugin.usenet`).
- `temporary_uri` (`string`): Temporary playback URI generated for just-in-time stream access.
- `expires_at` (`google.protobuf.Timestamp`): Expiration timestamp for the temporary URI.
- `rank` (`int32`): Preference ranking of the candidate (e.g. `1` = primary/highest priority).
- `resolution` (`VirtualStreamResolution`): Video resolution specifications:
  - `width` (`int32`): Horizontal pixel count (e.g. 1920, 3840).
  - `height` (`int32`): Vertical pixel count (e.g. 1080, 2160).
  - `label` (`string`): Standard display label (e.g. `"1080p"`, `"2160p"`, `"4k"`).
  - `frame_rate` (`double`): Frame rate in frames per second (e.g. 23.976, 60.0).
- `video_codec` (`string`): Video codec identifier (e.g. `"hevc"`, `"h264"`, `"av1"`, `"vp9"`).
- `audio_codec` (`string`): Primary audio codec identifier (e.g. `"aac"`, `"ac3"`, `"eac3"`, `"truehd"`, `"dts"`, `"flac"`).
- `hdr` (`VirtualStreamHDR`): Dynamic range metadata:
  - `is_hdr` (`bool`): High Dynamic Range support flag.
  - `format` (`string`): Format string (e.g. `"HDR10"`, `"HDR10+"`, `"HLG"`, `"Dolby Vision"`).
  - `has_dolby_vision` (`bool`): Dolby Vision flag.
  - `dolby_vision_profile` (`string`): Profile identifier (e.g. `"Profile 5"`, `"Profile 8.1"`).
- `bitrate` (`int64`): Stream bitrate in bits per second (bps).
- `file_size_bytes` (`int64`): Total media payload file size in bytes.
- `container` (`string`): Container format string (e.g. `"mkv"`, `"mp4"`, `"ts"`, `"hls"`).
- `audio_languages` (`repeated string`): List of available audio track language codes (ISO 639 / BCP 47).
- `subtitle_languages` (`repeated string`): List of available subtitle track language codes (ISO 639 / BCP 47).
- `availability` (`VirtualStreamAvailability`): Candidate-specific availability state and progress.
- `error` (`VirtualStreamError`): Candidate-specific fault metadata if degraded.
- `metadata` (`google.protobuf.Struct`): Extensible key-value metadata for provider-specific attributes.

#### `VirtualStreamResult`

Top-level resolution result containing:

- `result_id` (`string`): Unique resolution session ID.
- `provider_id` (`string`): Primary provider identifier.
- `candidates` (`repeated VirtualStreamCandidate`): Multi-candidate list ordered by rank.
- `availability` (`VirtualStreamAvailability`): Aggregate availability state across candidates.
- `error` (`VirtualStreamError`): Aggregate error metadata if resolution failed or is degraded.
- `metadata` (`google.protobuf.Struct`): Arbitrary resolution session metadata.

#### `VirtualStreamAvailability`

- `state` (`VirtualStreamAvailabilityState`): Enum value:
  - `VIRTUAL_STREAM_AVAILABILITY_STATE_UNSPECIFIED` (0)
  - `VIRTUAL_STREAM_AVAILABILITY_STATE_AVAILABLE` (1)
  - `VIRTUAL_STREAM_AVAILABILITY_STATE_PENDING` (2)
  - `VIRTUAL_STREAM_AVAILABILITY_STATE_UNRESOLVED` (3)
  - `VIRTUAL_STREAM_AVAILABILITY_STATE_DEGRADED` (4)
  - `VIRTUAL_STREAM_AVAILABILITY_STATE_UNAVAILABLE` (5)
- `message` (`string`): Operator- or user-safe descriptive text.
- `estimated_ready_at` (`google.protobuf.Timestamp`): Estimated ready timestamp for pending or buffering streams.
- `progress_percent` (`int32`): Completion percentage (0-100) for asynchronous resolution or pre-buffering.

#### `VirtualStreamError`

- `code` (`VirtualStreamErrorCode`): Standardized error category:
  - `VIRTUAL_STREAM_ERROR_CODE_UNSPECIFIED` (0)
  - `VIRTUAL_STREAM_ERROR_CODE_NOT_FOUND` (1)
  - `VIRTUAL_STREAM_ERROR_CODE_TEMPORARY` (2)
  - `VIRTUAL_STREAM_ERROR_CODE_EXPIRED` (3)
  - `VIRTUAL_STREAM_ERROR_CODE_RATE_LIMITED` (4)
  - `VIRTUAL_STREAM_ERROR_CODE_UNAUTHORIZED` (5)
  - `VIRTUAL_STREAM_ERROR_CODE_PROVIDER_FAILURE` (6)
  - `VIRTUAL_STREAM_ERROR_CODE_PERMANENT` (7)
- `message` (`string`): Safe error detail message.
- `retryable` (`bool`): True if caller may retry stream resolution.
- `retry_after` (`google.protobuf.Duration`): Suggested delay before retry attempt.

## Manifest Declaration

Plugins declare support for virtual stream resolution in `manifest.json`:

```json
{
  "plugin_id": "com.example.debrid",
  "version": "1.0.0",
  "capabilities": [
    {
      "type": "virtual_stream_provider.v1",
      "id": "debrid_stream_resolver",
      "display_name": "Debrid Stream Provider",
      "virtual_stream_provider": {
        "supported_media_types": ["movie", "episode"],
        "supported_containers": ["mkv", "mp4", "hls"],
        "supports_just_in_time_resolution": true,
        "supports_multiple_candidates": true
      }
    }
  ]
}
```

## Backward Compatibility & Design Principles

1. **Additive Contract**: Modifies no existing RPCs or messages.
2. **Provider-Neutral**: Avoids vendor-specific fields in primary contracts, allowing any stream backend (Usenet, Torrent, Debrid, S3, HLS) to map into standard candidates.
3. **Structured & Extensible**: Core fields use structured types (`resolution`, `hdr`, `availability`, `error`), while `google.protobuf.Struct` metadata fields allow custom extensions without schema revisions.
