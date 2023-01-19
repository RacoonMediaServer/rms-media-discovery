package provider

import (
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/mocks"
	"git.rms.local/RacoonMediaServer/rms-media-discovery/internal/model"
	"github.com/golang/mock/gomock"
	"net/http"
	"testing"
)

func TestImdbProvider_SearchMovies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockAccessProvider(ctrl)
	p := NewImdbProvider(m)

	tester := newMovieInfoProviderTester(t, m, p)
	tester.testSearchMovieCanceled()
	tester.testSearchMovieCannotGetApiKey()
	tester.testSearchMovieMaxAttemptsReached()
	tester.testSearchMovieUnexpectedStatusCode()

	exchanges := []exchange{
		makeOkExchange(t, "https://imdb-api.com/ru/API/SearchMovie/key/StarGate", http.MethodGet, "content/SearchMovie_1.json"),
		makeOkExchange(t, "https://imdb-api.com/ru/API/Title/key/tt0111282", http.MethodGet, "content/SearchMovie_2.json"),
		makeOkExchange(t, "https://imdb-api.com/ru/API/Title/key/tt0118480", http.MethodGet, "content/SearchMovie_3.json"),
	}

	movies := []model.Movie{
		{
			ID:          "tt0111282",
			Description: "В 1928-м году в Египте археологи обнаружили гигантское каменное кольцо, покрытое загадочными иероглифами. Американские военные, на протяжении десятилетий бившиеся над разгадкой послания, привлекают для исследования египтолога доктора Дэниэла Джексона — сторонника теории внеземного происхождения пирамид. Ему удаётся прочесть надпись. Навстречу неизведанному сквозь Звёздные Врата вместе с Джексоном отправляется отряд специального назначения под командованием полковника Джонатана О'Нила. Звёздные Врата переносят их за миллионы световых лет от Земли, на странную планету, в место, напоминающее Древний Египет: бескрайние пустыни, диковинные животные и рабы, поклоняющиеся богу солнца Ра. Чтобы найти путь домой и спасти человечество, пришельцам предстоит сразиться с армией последнего представителя могущественной инопланетной расы, узнавшего, что дверь на Землю может быть открыта вновь, спустя тысячелетия.",
			Genres:      []string{"Действие", "Приключение", "Научная фантастика"},
			Poster:      "https://m.media-amazon.com/images/M/MV5BYWEyYTQzNzQtZTg5OS00NDhkLTg1NjYtMzA5Y2MyYjYzNWQ5L2ltYWdlL2ltYWdlXkEyXkFqcGdeQXVyNTAyODkwOQ@@._V1_Ratio0.6762_AL_.jpg",
			Preview:     "",
			Rating:      7,
			Seasons:     0,
			Title:       "Stargate",
			Type:        model.MovieType_Movie,
			Year:        1994,
		},
		{
			ID:          "tt0118480",
			Description: "«Звёздные врата» соединяют далекие пределы Млечного Пути и открывают двери для мгновенного межпланетного перемещения. Отряд ЗВ-1 с их помощью исследует нашу галактику и защищает Землю от нападок пришельцев.",
			Genres:      []string{"Действие", "Приключение", "Драма"},
			Poster:      "https://m.media-amazon.com/images/M/MV5BMTc3MjEwMTc5N15BMl5BanBnXkFtZTcwNzQ2NjQ4NA@@._V1_Ratio0.6762_AL_.jpg",
			Preview:     "",
			Rating:      8.4,
			Seasons:     10,
			Title:       "Stargate SG-1",
			Type:        model.MovieType_TvSeries,
			Year:        1997,
		},
	}

	tester.testSearchMovie("StarGate", movies, nil, exchanges)
}
