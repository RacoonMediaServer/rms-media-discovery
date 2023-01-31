# rms-media-discovery

This is a backend-service, which provide ability of searching and downloading movies/TV series/music via torrent
trackers for free. Is uses a set of providers (external systems), which contents searched
information. `rms-media-service` could be described as web-crawler application.

## Capabilities

* searching information about movies and TV series;
* searching torrent files for downloading media on different torrent trackers;
* service user management;
* external systems accounts management;
* Prometheus monitoring.

## Packages

Some packages of source code are importable. Some useful of this:

* `client` - Swagger-generated client to the service API;
* `provider` - crawlers and API-clients for various external systems;
* `navigator` - comfortable wrapper for headless browser;
* etc.

## Providers

### Search information about movies

* [IMDB](https://www.imdb.com/)
* [Кинопоиск](https://www.kinopoisk.ru/)

### Torrent trackers

* [RuTracker.org](https://rutracker.org/)
* [Rutor](http://www.rutor.info/)
* [The Pirate Bay](https://thepiratebay.org/)

## API

Service have RESTful JSON API, described as OpenAPI schema [here](api/discovery.yml). Auth via X-Token is supported.

## Build & Run

### Dependencies

* [MongoDb](https://www.mongodb.com/)
* [Prometheus](https://prometheus.io/) (**optional**)
* [Chromium](https://www.chromium.org/chromium-projects/)

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
docker build -t rms-media-discovery .
docker run -p 8080:8080 -e RMS_DATABASE=mongodb://192.168.1.19:27017 rms-media-discovery 
```

## Accounts management

Some external systems need user accounts or limited API keys for providing functionality. You can register and append
they to database manually or via [API](api/discovery.yml). Each account links to external system id. They can be:

* `imdb` - [IMDB](https://www.imdb.com/) API key. You can get it [here](https://imdb-api.com/Identity/Account/Register).
* `kinopoisk` - [Кинопоиск](https://www.kinopoisk.ru/) API key. You can get i [here](https://kinopoisk.dev/)
* `2captcha` - [2Captcha](https://2captcha.com/) API Key. Needs to resolve captchas for some external systems login
* `rutracker` - [RuTracker.org](https://rutracker.org/) user account.

Best practice is use a few accounts of each external system (except 2Captcha) for avoiding ban or limits exceeds.