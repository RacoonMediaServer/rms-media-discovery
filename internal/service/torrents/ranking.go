package torrents

import (
	"container/list"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/utils"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/media"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/model"
	"github.com/antzucaro/matchr"
	"math"
	"sort"
)

type rankContext struct {
	torrents []model.Torrent
	q        *model.SearchQuery
	trash    *list.List
	results  *list.List
}

func newRankContext(torrents []model.Torrent, q *model.SearchQuery) *rankContext {
	ctx := &rankContext{
		torrents: torrents,
		q:        q,
		trash:    list.New(),
		results:  list.New(),
	}
	for i := 0; i < len(torrents); i++ {
		ctx.results.PushBack(i)
	}
	return ctx
}

func (ctx *rankContext) filter(f func(t *model.Torrent) bool) {
	var next *list.Element
	for cur := ctx.results.Front(); cur != nil; cur = next {
		next = cur.Next()
		val := cur.Value.(int)
		if f(&ctx.torrents[val]) {
			ctx.results.Remove(cur)
			ctx.trash.PushBack(val)
		}
	}
}

func (ctx *rankContext) makeResult() []model.Torrent {
	results := make([]model.Torrent, 0, ctx.results.Len())
	for cur := ctx.results.Front(); cur != nil; cur = cur.Next() {
		results = append(results, ctx.torrents[cur.Value.(int)])
	}
	if len(results) < int(ctx.q.Limit) {
		// разбавим раздачу не очень релевантными результатами
		canAppend := int(ctx.q.Limit) - len(results)
		for cur := ctx.trash.Front(); cur != nil && canAppend != 0; cur, canAppend = cur.Next(), canAppend-1 {
			results = append(results, ctx.torrents[cur.Value.(int)])
		}
	}

	results = utils.Bound(results, ctx.q.Limit)

	return results
}

func (ctx *rankContext) sortList(l *list.List, f func(lhs *model.Torrent, rhs *model.Torrent) bool) *list.List {
	// не оч хорошо, но зато O(N log N) суммарно - все равно лучше, если бы использовали слайсы (из-за удаления из середины в filter)
	tmp := make([]int, 0, l.Len())
	for cur := l.Front(); cur != nil; cur = cur.Next() {
		tmp = append(tmp, cur.Value.(int))
	}
	sort.SliceStable(tmp, func(i, j int) bool {
		return f(&ctx.torrents[tmp[i]], &ctx.torrents[tmp[j]])
	})
	res := list.New()
	for _, n := range tmp {
		res.PushBack(n)
	}
	return res
}

func (ctx *rankContext) sort(f func(lhs *model.Torrent, rhs *model.Torrent) bool) {
	ctx.results = ctx.sortList(ctx.results, f)
	ctx.trash = ctx.sortList(ctx.trash, f)
}

func getMinDistance(titles []string, q string) int {
	min := math.MaxInt
	for _, t := range titles {
		d := matchr.Levenshtein(t, q)
		if d < min {
			min = d
		}
	}
	return min
}

type sortCallback func(lhs *model.Torrent, rhs *model.Torrent) int
type sortChain struct {
	cb []sortCallback
}

func (c *sortChain) add(cb sortCallback) {
	c.cb = append(c.cb, cb)
}

func (c *sortChain) sortFunc() func(lhs *model.Torrent, rhs *model.Torrent) bool {
	return func(lhs *model.Torrent, rhs *model.Torrent) bool {
		for _, cb := range c.cb {
			res := cb(lhs, rhs)
			if res < 0 {
				return true
			}
			if res > 0 {
				return false
			}
		}
		return true
	}
}

func rank(torrents []model.Torrent, q model.SearchQuery) []model.Torrent {
	ctx := newRankContext(torrents, &q)
	sort := sortChain{}

	// в любом случае задаем сортировку результатов по степени удаленности от поискового запроса в первую очередь
	sort.add(func(lhs *model.Torrent, rhs *model.Torrent) int {
		return getMinDistance(lhs.Info.Titles, q.Query) - getMinDistance(rhs.Info.Titles, q.Query)
	})

	// фильтруем раздачи, у которых нету сидов - они все равно бесполезны
	ctx.filter(func(t *model.Torrent) bool {
		return t.Seeders == 0
	})

	if q.Type != media.Other {
		// фильтруем раздачи, которые по эвристически определенному типу не подходят
		ctx.filter(func(t *model.Torrent) bool {
			return t.Info.Type != q.Type
		})
	}

	if q.Type == media.Movies {
		if q.Season != nil {
			// фильтруем раздачи, которые не содержат искомого сезона
			ctx.filter(func(t *model.Torrent) bool {
				found := false
				for _, s := range t.Info.Seasons {
					if s == *q.Season {
						found = true
					}
				}
				return !found
			})
			// сортируем, если в раздаче сезонов больше чем запросили
			sort.add(func(lhs *model.Torrent, rhs *model.Torrent) int {
				return len(lhs.Info.Seasons) - len(rhs.Info.Seasons)
			})
		}
		// сортируем по качеству
		sort.add(func(lhs *model.Torrent, rhs *model.Torrent) int {
			return int(rhs.Info.Quality - lhs.Info.Quality)
		})
	}

	// после всех сортировок - сортируем по количеству сидов
	sort.add(func(lhs *model.Torrent, rhs *model.Torrent) int {
		return int(rhs.Seeders) - int(lhs.Seeders)
	})

	// сортируем по всем критериям
	ctx.sort(sort.sortFunc())

	return ctx.makeResult()
}
