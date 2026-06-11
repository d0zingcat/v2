// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package ui // import "miniflux.app/v2/internal/ui"

import (
	"net/http"

	"miniflux.app/v2/internal/http/request"
	"miniflux.app/v2/internal/http/response"
	"miniflux.app/v2/internal/model"
	"miniflux.app/v2/internal/ui/view"
)

func (h *handler) showDefaultAttentionLanePage(w http.ResponseWriter, r *http.Request) {
	response.HTMLRedirect(w, r, h.routePath("/lanes/%s", model.AttentionLanePriority))
}

func (h *handler) showAttentionLanePage(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	laneSlug := request.RouteStringParam(r, "lane")
	lane, ok := model.FindAttentionLane(laneSlug)
	if !ok {
		response.HTMLNotFound(w, r)
		return
	}

	offset := request.QueryIntParam(r, "offset", 0)

	entries, total, err := h.store.NewEntryQueryBuilder(user.ID).
		WithStatuses(model.EntryStatusUnread).
		WithAttentionLane(lane.Slug).
		WithSorting(user.EntryOrder, user.EntryDirection).
		WithSorting("id", user.EntryDirection).
		WithOffset(offset).
		WithLimit(user.EntriesPerPage).
		WithGloballyVisible().
		WithoutContent().
		GetEntriesWithCount()
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	if offset >= total && total > 0 {
		offset = 0

		entries, total, err = h.store.NewEntryQueryBuilder(user.ID).
			WithStatuses(model.EntryStatusUnread).
			WithAttentionLane(lane.Slug).
			WithSorting(user.EntryOrder, user.EntryDirection).
			WithSorting("id", user.EntryDirection).
			WithLimit(user.EntriesPerPage).
			WithGloballyVisible().
			WithoutContent().
			GetEntriesWithCount()
		if err != nil {
			response.HTMLServerError(w, r, err)
			return
		}
	}

	view := view.New(h.tpl, r)
	view.Set("entries", entries)
	view.Set("total", total)
	view.Set("pagination", getPagination(h.routePath("/lanes/%s", lane.Slug), total, offset, user.EntriesPerPage))
	view.Set("menu", "lanes")
	view.Set("user", user)
	view.Set("lane", lane)
	view.Set("lanes", model.AttentionLanes())
	navMetadata, _ := h.store.GetNavMetadata(user.ID)
	view.Set("countUnread", navMetadata.CountUnread)
	view.Set("countErrorFeeds", navMetadata.CountErrorFeeds)
	view.Set("hasSaveEntry", navMetadata.HasSaveEntry)

	response.HTML(w, r, view.Render("attention_lane"))
}

func (h *handler) showAttentionLaneEntryPage(w http.ResponseWriter, r *http.Request) {
	user, err := h.store.UserByID(request.UserID(r))
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	laneSlug := request.RouteStringParam(r, "lane")
	lane, ok := model.FindAttentionLane(laneSlug)
	if !ok {
		response.HTMLNotFound(w, r)
		return
	}

	entryID := request.RouteInt64Param(r, "entryID")

	entry, err := h.store.NewEntryQueryBuilder(user.ID).
		WithEntryIDs(entryID).
		WithAttentionLane(lane.Slug).
		WithGloballyVisible().
		GetEntry()
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	if entry == nil {
		response.HTMLRedirect(w, r, h.routePath("/lanes/%s", lane.Slug))
		return
	}

	if entry.Status == model.EntryStatusRead {
		err = h.store.SetEntriesStatus(user.ID, []int64{entry.ID}, model.EntryStatusUnread)
		if err != nil {
			response.HTMLServerError(w, r, err)
			return
		}
	}

	prevEntry, nextEntry, err := h.store.NewEntryPaginationBuilder(user.ID, entry.ID, user.EntryOrder, user.EntryDirection).
		WithStatus(model.EntryStatusUnread).
		WithAttentionLane(lane.Slug).
		WithGloballyVisible().
		Entries()
	if err != nil {
		response.HTMLServerError(w, r, err)
		return
	}

	nextEntryRoute := ""
	if nextEntry != nil {
		nextEntryRoute = h.routePath("/lanes/%s/entry/%d", lane.Slug, nextEntry.ID)
	}

	prevEntryRoute := ""
	if prevEntry != nil {
		prevEntryRoute = h.routePath("/lanes/%s/entry/%d", lane.Slug, prevEntry.ID)
	}

	if entry.ShouldMarkAsReadOnView(user) {
		entry.Status = model.EntryStatusRead
	}

	if entry.Status == model.EntryStatusRead {
		err = h.store.SetEntriesStatus(user.ID, []int64{entry.ID}, model.EntryStatusRead)
		if err != nil {
			response.HTMLServerError(w, r, err)
			return
		}
	}

	if user.AlwaysOpenExternalLinks {
		response.HTMLRedirect(w, r, entry.URL)
		return
	}

	view := view.New(h.tpl, r)
	view.Set("entry", entry)
	view.Set("prevEntry", prevEntry)
	view.Set("nextEntry", nextEntry)
	view.Set("nextEntryRoute", nextEntryRoute)
	view.Set("prevEntryRoute", prevEntryRoute)
	view.Set("menu", "lanes")
	view.Set("user", user)
	navMetadata, _ := h.store.GetNavMetadata(user.ID)
	view.Set("countUnread", navMetadata.CountUnread)
	view.Set("countErrorFeeds", navMetadata.CountErrorFeeds)
	view.Set("hasSaveEntry", navMetadata.HasSaveEntry)

	response.HTML(w, r, view.Render("entry"))
}
