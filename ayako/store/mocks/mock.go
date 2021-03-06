package mock_store

import (
	"github.com/deissh/osu-lazer/ayako/store"
	"github.com/golang/mock/gomock"
)

// Interface assertion
var _ store.Store = (*MockedStore)(nil)

type MockedStore struct {
	ctrl *gomock.Controller

	beatmap    *MockBeatmap
	beatmapSet *MockBeatmapSet
}

func InitStore(ctrl *gomock.Controller) MockedStore {
	return MockedStore{
		ctrl:       ctrl,
		beatmap:    NewMockBeatmap(ctrl),
		beatmapSet: NewMockBeatmapSet(ctrl),
	}
}

func (ss MockedStore) Beatmap() store.Beatmap       { return ss.beatmap }
func (ss MockedStore) BeatmapSet() store.BeatmapSet { return ss.beatmapSet }

func (ss MockedStore) BeatmapExpect() *MockBeatmapMockRecorder       { return ss.beatmap.EXPECT() }
func (ss MockedStore) BeatmapSetExpect() *MockBeatmapSetMockRecorder { return ss.beatmapSet.EXPECT() }
