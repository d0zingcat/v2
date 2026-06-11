// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package model // import "miniflux.app/v2/internal/model"

const (
	AttentionLanePriority   = "priority"
	AttentionLaneFastNews   = "fast-news"
	AttentionLaneSlowReads  = "slow-reads"
	AttentionLaneTechRadar  = "tech-radar"
	AttentionLaneEverything = "everything"
)

// AttentionLane describes a deterministic reading lane derived from existing
// entry, feed, tag, and category metadata.
type AttentionLane struct {
	Slug        string
	Title       string
	Description string
}

func AttentionLanes() []AttentionLane {
	return []AttentionLane{
		{
			Slug:        AttentionLanePriority,
			Title:       "Priority",
			Description: "Starred items and sources or tags named priority, important, must-read, or VIP.",
		},
		{
			Slug:        AttentionLaneFastNews,
			Title:       "Fast News",
			Description: "Short reads and feeds, categories, or tags that look like high-volume news streams.",
		},
		{
			Slug:        AttentionLaneSlowReads,
			Title:       "Slow Reads",
			Description: "Longer items and feeds, categories, or tags intended for focused reading.",
		},
		{
			Slug:        AttentionLaneTechRadar,
			Title:       "Tech Radar",
			Description: "Technical signals such as releases, security, infrastructure, AI, and developer tooling.",
		},
		{
			Slug:        AttentionLaneEverything,
			Title:       "Everything",
			Description: "All visible unread entries in the normal reading queue.",
		},
	}
}

func FindAttentionLane(slug string) (AttentionLane, bool) {
	for _, lane := range AttentionLanes() {
		if lane.Slug == slug {
			return lane, true
		}
	}

	return AttentionLane{}, false
}
