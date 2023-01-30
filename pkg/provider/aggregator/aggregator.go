package aggregator

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/provider"
)

// Policy определяет политику аггрегации результатов
type Policy int

const (
	// PriorityPolicy собираем наиболее релевантные результаты
	PriorityPolicy Policy = iota

	// FastPolicy выдаем как можно быстрее результаты
	FastPolicy
)

func NewTorrentProvider(policy Policy, providers []provider.TorrentsProvider) provider.TorrentsProvider {
	switch policy {
	case FastPolicy:
		return &torrentsFastAggregator{providers: providers}

	case PriorityPolicy:
		return &torrentsPriorityAggregator{providers: providers}
	}

	panic("policy not implemented")
}
