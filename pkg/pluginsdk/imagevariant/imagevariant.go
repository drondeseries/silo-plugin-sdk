// Package imagevariant exposes the canonical image-variant hints Silo sends to
// plugins that resolve image paths into URLs.
//
// A variant is a semantic size hint, not an exact pixel contract: the sender
// says roughly how large the image will be rendered and the resolver picks the
// closest size its upstream offers. Variants travel as plain strings in
// ResolveImageURLRequest.variant, ResolveImageURLsRequest.variant, and
// ResolveCatalogImageURLsRequest.variant. An empty variant uses the receiver's
// default: the plugin default for the metadata-provider requests and the host
// default for the RuntimeHost request.
//
// Two contract rules govern the vocabulary:
//
//  1. The set is open and may grow additively. New variants are added between or
//     beyond the existing bands as Silo's clients gain new rendering surfaces;
//     adding one is not a breaking change, so a plugin must never assume the
//     constants in this package are exhaustive.
//  2. A plugin receiving an unknown variant MUST degrade gracefully to its
//     nearest supported size (or its default) and MUST NOT return an error. A
//     slightly wrong image size is always preferable to a failed resolve, which
//     surfaces to users as a missing poster.
//
// The idiomatic implementation is a switch with a default arm that returns the
// plugin's fallback size, so unrecognized values fall through instead of failing.
package imagevariant

const (
	// Card is the smallest band: dense grid and list thumbnails, search
	// results, and cast rows. Targets roughly 300px-wide posters.
	Card = "card"

	// Featured is the mid band: hero artwork, detail-page posters, and section
	// headers. Targets roughly 500px-wide posters and logos and 1280px-wide
	// backdrops.
	Featured = "featured"

	// Large sits between Featured and Full for high-density displays and
	// client-selected larger image sizes. Targets roughly 780px-wide posters
	// and stills and 1280px-wide logos and backdrops.
	Large = "large"

	// Full is the largest sized band, for full-screen and near-source
	// rendering. Targets the largest sized rendition the upstream offers,
	// falling back to the original asset when no sized rendition is large
	// enough.
	Full = "full"

	// Original requests the unresized source asset, with no host or plugin
	// downscaling. It is also the conventional fallback for an empty or
	// unrecognized variant.
	Original = "original"
)
