package utils

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"sort"
)

type torrentList []model.Torrent

func (t torrentList) Len() int {
	return len(t)
}

func (t torrentList) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t torrentList) Less(i, j int) bool {
	return t[i].Seeders > t[j].Seeders
}

func SortTorrents(t []model.Torrent) {
	sort.Sort(torrentList(t))
}
