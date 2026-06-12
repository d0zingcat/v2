// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package storage // import "miniflux.app/v2/internal/storage"

func attentionPriorityCondition() string {
	return `(
		e.starred IS TRUE OR
		lower(c.title) IN ('priority', 'important', 'must read', 'must-read', 'vip') OR
		lower(f.title) IN ('priority', 'important', 'must read', 'must-read', 'vip') OR
		lower(array_to_string(e.tags, ' ')) SIMILAR TO '%(priority|important|must-read|must read|vip)%'
	)`
}

func attentionFastNewsCondition() string {
	return `(
		(
			SELECT count(*)
			FROM entries recent
			WHERE recent.user_id = e.user_id
				AND recent.feed_id = e.feed_id
				AND recent.created_at >= now() - interval '7 days'
		) >= 10 OR
		e.reading_time <= 3 OR
		lower(c.title) SIMILAR TO '%(news|fast|hacker news|hn|reddit)%' OR
		lower(f.title) SIMILAR TO '%(news|hacker news|hn|reddit)%' OR
		lower(array_to_string(e.tags, ' ')) SIMILAR TO '%(news|fast|hn)%'
	)`
}

func attentionSlowReadsCondition() string {
	return `(
		(
			SELECT count(*)
			FROM entries recent
			WHERE recent.user_id = e.user_id
				AND recent.feed_id = e.feed_id
				AND recent.created_at >= now() - interval '7 days'
		) <= 2 OR
		e.reading_time >= 8 OR
		lower(c.title) SIMILAR TO '%(slow|long|essay|read later|deep)%' OR
		lower(f.title) SIMILAR TO '%(longform|essay|deep)%' OR
		lower(array_to_string(e.tags, ' ')) SIMILAR TO '%(slow|long|essay|deep)%'
	)`
}

func attentionTechRadarCondition() string {
	return `(
		lower(c.title) SIMILAR TO '%(tech|engineering|developer|programming|radar)%' OR
		lower(f.title) SIMILAR TO '%(github|engineering|developer|programming|security|release|changelog|database|postgres|golang|rust|kubernetes|ai|llm)%' OR
		lower(array_to_string(e.tags, ' ')) SIMILAR TO '%(tech|engineering|developer|programming|security|release|database|ai|llm|kubernetes)%' OR
		lower(e.title) SIMILAR TO '%(release|changelog|security|vulnerability|database|postgres|golang|go |rust|kubernetes|ai|llm|openai|performance|frontend|backend|api)%'
	)`
}
