package util

import (
	"slices"

	"github.com/dweymouth/supersonic/backend/mediaprovider"
	"github.com/dweymouth/supersonic/sharedutil"
)

type TrackListModel struct {
	Item     mediaprovider.MediaItem
	Selected bool
}

// Returns the item as a *mediaprovider.Track, or panics if not a Track
// Use for tracklists that can only contain tracks (ie not PlayQueueList)
func (t TrackListModel) Track() *mediaprovider.Track {
	return t.Item.(*mediaprovider.Track)
}

func ToTrackListModels(trs []*mediaprovider.Track) []*TrackListModel {
	return sharedutil.MapSlice(trs, func(tr *mediaprovider.Track) *TrackListModel {
		return &TrackListModel{Item: tr, Selected: false}
	})
}

func SelectedTracks(items []*TrackListModel) []*mediaprovider.Track {
	return sharedutil.FilterMapSlice(items, func(tm *TrackListModel) (*mediaprovider.Track, bool) {
		return tm.Track(), tm.Selected
	})
}

func SelectedItemIDs(items []*TrackListModel) []string {
	return sharedutil.FilterMapSlice(items, func(tm *TrackListModel) (string, bool) {
		return tm.Item.Metadata().ID, tm.Selected
	})
}

func SelectItem(items []*TrackListModel, idx int) {
	if items[idx].Selected {
		return
	}
	UnselectAllItems(items)
	items[idx].Selected = true
}

func SelectAllItems(items []*TrackListModel) {
	for _, tm := range items {
		tm.Selected = true
	}
}

func UnselectAllItems(items []*TrackListModel) {
	for _, tm := range items {
		tm.Selected = false
	}
}

func SelectItemRange(items []*TrackListModel, idx int) {
	if items[idx].Selected {
		return
	}
	lastSelected := -1
	for i := len(items) - 1; i >= 0; i-- {
		if items[i].Selected {
			lastSelected = i
			break
		}
	}
	if lastSelected < 0 {
		items[idx].Selected = true
		return
	}
	from := min(idx, lastSelected)
	to := max(idx, lastSelected)
	for i := from; i <= to; i++ {
		items[i].Selected = true
	}
}

func FindItemByID(items []*TrackListModel, id string) (mediaprovider.MediaItem, int) {
	idx := slices.IndexFunc(items, func(tr *TrackListModel) bool {
		return tr.Item.Metadata().ID == id
	})
	if idx >= 0 {
		return items[idx].Item, idx
	}
	return nil, -1
}
