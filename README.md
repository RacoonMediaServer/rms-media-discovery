# rms-media-discovery

This is a backend-service, which provide ability of searching and downloading movies/TV series/music via torrent
trackers for free.

## Capabilities

* searching information about movies and TV series;
* searching torrent files for downloading media on different torrent trackers;
* service user management;
* external systems accounts management;
* Prometheus monitoring.

## Providers

### Search information about movies

* [IMDB](https://www.imdb.com/)
* [Кинопоиск](https://www.kinopoisk.ru/)

### Torrent trackers

* [RuTracker.org](https://rutracker.org/)
* [Rutor](http://www.rutor.info/)

## API

Service have RESTful JSON API, described as OpenAPI schema [here](api/discovery.yml). Auth via X-Token is supported.

## Build & Run

### Dependencies

* [MongoDb](https://www.mongodb.com/)
* [Prometheus](https://prometheus.io/) (**optional**)

### Admin Key

At the first run of service, admin key will be generated automatically. You can find out the key via logs or database.
### Command Line Arguments

```
./rms-media-discovery [-host host] [-port port] [-db db] [-verbose] [-help]
```

| Flag           | Default Value               | Description               |
|----------------|-----------------------------|---------------------------|
| `-db string`   | `mongodb://localhost:27017` | MongoDB connection string |
| `-host string` | `127.0.0.1`                 | Server IP address         |
| `-port int`    | `8080`                      | Server port               |
| `-verbose`     |                             | Verbose mode              |
| `-help`        |                             | Show help                 |

### Docker

```bash
docker run -p 8080:8080 -e RMS_DATABASE=mongodb://192.168.1.19:27017 rms-media-discovery 
```

