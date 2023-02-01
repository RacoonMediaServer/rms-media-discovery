package model

import (
	"sort"
)

// OrderByFunc is a func for sort results
type OrderByFunc func(a, b *Torrent) bool

func OrderBySeeders(a, b *Torrent) bool {
	return a.Seeders > b.Seeders
}

func OrderBySize(a, b *Torrent) bool {
	return a.SizeMB > b.SizeMB
}

type torrentList struct {
	t []Torrent
	f OrderByFunc
}

func (t torrentList) Len() int {
	return len(t.t)
}

func (t torrentList) Swap(i, j int) {
	t.t[i], t.t[j] = t.t[j], t.t[i]
}

func (t torrentList) Less(i, j int) bool {
	return t.f(&t.t[i], &t.t[j])
}

func SortTorrents(t []Torrent, f OrderByFunc) {
	sort.Sort(&torrentList{t: t, f: f})
}
