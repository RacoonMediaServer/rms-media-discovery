// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "API for Racoon Media Server Project",
    "title": "Media Discovery API",
    "version": "1.3.0"
  },
  "host": "136.244.108.126",
  "paths": {
    "/accounts": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "accounts"
        ],
        "summary": "Получить список список акканутов и токенов к внешним системам",
        "operationId": "getAccounts",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "results"
              ],
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/Account"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      },
      "post": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "accounts"
        ],
        "summary": "Создать новый аккаунт",
        "operationId": "createAccount",
        "parameters": [
          {
            "name": "account",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Account"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "id"
              ],
              "properties": {
                "id": {
                  "type": "string"
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/accounts/{id}": {
      "delete": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "accounts"
        ],
        "summary": "Удалить аккаунт",
        "operationId": "deleteAccount",
        "parameters": [
          {
            "type": "string",
            "description": "ID аккаунта",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          },
          "404": {
            "description": "Аккаунт не найден"
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/movies/search": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "description": "Поиск фильмов и сериалов по названию на различных платформах",
        "tags": [
          "movies"
        ],
        "summary": "Поиск фильмов и сериалов",
        "operationId": "searchMovies",
        "parameters": [
          {
            "maxLength": 128,
            "minLength": 2,
            "type": "string",
            "description": "Искомый запрос",
            "name": "q",
            "in": "query",
            "required": true
          },
          {
            "minimum": 1,
            "type": "integer",
            "description": "Ограничение на кол-во результатов",
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "results"
              ],
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/SearchMoviesResult"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/movies/{id}": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "description": "Получение информации о фильме или сериале",
        "tags": [
          "movies"
        ],
        "summary": "Получение информации о фильме или сериала",
        "operationId": "getMovieInfo",
        "parameters": [
          {
            "type": "string",
            "description": "ID фильма/сериала",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/SearchMoviesResult"
            }
          },
          "404": {
            "description": "Фильм не найден"
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/torrents/download": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "description": "Позволяет скачать торрент-файл, с помощью которого можно скачать контент",
        "produces": [
          "application/octet-stream"
        ],
        "tags": [
          "torrents"
        ],
        "summary": "Загрузка торрент-файла",
        "operationId": "downloadTorrent",
        "parameters": [
          {
            "type": "string",
            "description": "Хеш ссылки на результат поиска",
            "name": "link",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "string",
              "format": "binary"
            }
          },
          "404": {
            "description": "Неверный хеш ссылки"
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/torrents/search": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "description": "Поиск раздач на различных платформах",
        "tags": [
          "torrents"
        ],
        "summary": "Поиск контента на торрент-трекерах",
        "operationId": "searchTorrents",
        "parameters": [
          {
            "maxLength": 128,
            "minLength": 2,
            "type": "string",
            "description": "Искомый запрос",
            "name": "q",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "movies",
              "music",
              "books",
              "others"
            ],
            "type": "string",
            "description": "Подсказка, какого типа торренты искать",
            "name": "type",
            "in": "query"
          },
          {
            "minimum": 1,
            "type": "integer",
            "description": "Ограничение на кол-во результатов",
            "name": "limit",
            "in": "query"
          },
          {
            "minimum": 1900,
            "type": "integer",
            "description": "Год выхода (для фильмов и сериалов)",
            "name": "year",
            "in": "query"
          },
          {
            "minimum": 1,
            "type": "integer",
            "description": "Номер сезона (для сериалов)",
            "name": "season",
            "in": "query"
          },
          {
            "type": "boolean",
            "default": false,
            "description": "Строго отсеивать раздачи, эвристически определенное имя которых не соответствует строчке запроса",
            "name": "strong",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "results"
              ],
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/SearchTorrentsResult"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    }
  },
  "definitions": {
    "Account": {
      "type": "object",
      "required": [
        "service"
      ],
      "properties": {
        "id": {
          "type": "string"
        },
        "limit": {
          "type": "integer",
          "default": 0
        },
        "login": {
          "type": "string",
          "minLength": 1
        },
        "password": {
          "type": "string"
        },
        "service": {
          "type": "string"
        },
        "token": {
          "type": "string"
        }
      }
    },
    "SearchMoviesResult": {
      "type": "object",
      "required": [
        "id",
        "title"
      ],
      "properties": {
        "description": {
          "type": "string"
        },
        "genres": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "id": {
          "type": "string"
        },
        "poster": {
          "type": "string"
        },
        "preview": {
          "type": "string"
        },
        "rating": {
          "type": "number"
        },
        "seasons": {
          "type": "integer"
        },
        "title": {
          "type": "string"
        },
        "type": {
          "type": "string",
          "enum": [
            "tv-series",
            "film"
          ]
        },
        "year": {
          "type": "integer"
        }
      }
    },
    "SearchTorrentsResult": {
      "type": "object",
      "required": [
        "link",
        "title",
        "size",
        "seeders"
      ],
      "properties": {
        "format": {
          "description": "Формат",
          "type": "string"
        },
        "link": {
          "type": "string"
        },
        "quality": {
          "description": "Качество видео",
          "type": "string",
          "enum": [
            "",
            "480p",
            "720p",
            "1080p",
            "2160p"
          ]
        },
        "rip": {
          "description": "Rip для видео",
          "type": "string"
        },
        "seasons": {
          "description": "Количество сезонов в сериале (если это сериал)",
          "type": "array",
          "items": {
            "type": "integer",
            "minimum": 1
          }
        },
        "seeders": {
          "type": "integer"
        },
        "size": {
          "type": "integer"
        },
        "subtitles": {
          "description": "Коды языков, на которых предоставлены субтитры",
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "title": {
          "type": "string"
        },
        "voice": {
          "description": "Озвучка",
          "type": "string"
        }
      }
    },
    "principal": {
      "type": "object",
      "properties": {
        "canManageAccounts": {
          "type": "boolean"
        },
        "token": {
          "type": "string"
        }
      }
    }
  },
  "securityDefinitions": {
    "key": {
      "type": "apiKey",
      "name": "x-token",
      "in": "header"
    }
  },
  "tags": [
    {
      "description": "Фильмы/сериалы",
      "name": "movies"
    },
    {
      "description": "Торренты",
      "name": "torrents"
    },
    {
      "description": "Администрирование аккаунтов",
      "name": "accounts"
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "description": "API for Racoon Media Server Project",
    "title": "Media Discovery API",
    "version": "1.3.0"
  },
  "host": "136.244.108.126",
  "paths": {
    "/accounts": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "accounts"
        ],
        "summary": "Получить список список акканутов и токенов к внешним системам",
        "operationId": "getAccounts",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "results"
              ],
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/Account"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      },
      "post": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "accounts"
        ],
        "summary": "Создать новый аккаунт",
        "operationId": "createAccount",
        "parameters": [
          {
            "name": "account",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/Account"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "id"
              ],
              "properties": {
                "id": {
                  "type": "string"
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/accounts/{id}": {
      "delete": {
        "security": [
          {
            "key": []
          }
        ],
        "tags": [
          "accounts"
        ],
        "summary": "Удалить аккаунт",
        "operationId": "deleteAccount",
        "parameters": [
          {
            "type": "string",
            "description": "ID аккаунта",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK"
          },
          "404": {
            "description": "Аккаунт не найден"
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/movies/search": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "description": "Поиск фильмов и сериалов по названию на различных платформах",
        "tags": [
          "movies"
        ],
        "summary": "Поиск фильмов и сериалов",
        "operationId": "searchMovies",
        "parameters": [
          {
            "maxLength": 128,
            "minLength": 2,
            "type": "string",
            "description": "Искомый запрос",
            "name": "q",
            "in": "query",
            "required": true
          },
          {
            "minimum": 1,
            "type": "integer",
            "description": "Ограничение на кол-во результатов",
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "results"
              ],
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/SearchMoviesResult"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/movies/{id}": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "description": "Получение информации о фильме или сериале",
        "tags": [
          "movies"
        ],
        "summary": "Получение информации о фильме или сериала",
        "operationId": "getMovieInfo",
        "parameters": [
          {
            "type": "string",
            "description": "ID фильма/сериала",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/SearchMoviesResult"
            }
          },
          "404": {
            "description": "Фильм не найден"
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/torrents/download": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "description": "Позволяет скачать торрент-файл, с помощью которого можно скачать контент",
        "produces": [
          "application/octet-stream"
        ],
        "tags": [
          "torrents"
        ],
        "summary": "Загрузка торрент-файла",
        "operationId": "downloadTorrent",
        "parameters": [
          {
            "type": "string",
            "description": "Хеш ссылки на результат поиска",
            "name": "link",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "string",
              "format": "binary"
            }
          },
          "404": {
            "description": "Неверный хеш ссылки"
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    },
    "/torrents/search": {
      "get": {
        "security": [
          {
            "key": []
          }
        ],
        "description": "Поиск раздач на различных платформах",
        "tags": [
          "torrents"
        ],
        "summary": "Поиск контента на торрент-трекерах",
        "operationId": "searchTorrents",
        "parameters": [
          {
            "maxLength": 128,
            "minLength": 2,
            "type": "string",
            "description": "Искомый запрос",
            "name": "q",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "movies",
              "music",
              "books",
              "others"
            ],
            "type": "string",
            "description": "Подсказка, какого типа торренты искать",
            "name": "type",
            "in": "query"
          },
          {
            "minimum": 1,
            "type": "integer",
            "description": "Ограничение на кол-во результатов",
            "name": "limit",
            "in": "query"
          },
          {
            "minimum": 1900,
            "type": "integer",
            "description": "Год выхода (для фильмов и сериалов)",
            "name": "year",
            "in": "query"
          },
          {
            "minimum": 1,
            "type": "integer",
            "description": "Номер сезона (для сериалов)",
            "name": "season",
            "in": "query"
          },
          {
            "type": "boolean",
            "default": false,
            "description": "Строго отсеивать раздачи, эвристически определенное имя которых не соответствует строчке запроса",
            "name": "strong",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "type": "object",
              "required": [
                "results"
              ],
              "properties": {
                "results": {
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/SearchTorrentsResult"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Ошибка на стороне сервера"
          }
        }
      }
    }
  },
  "definitions": {
    "Account": {
      "type": "object",
      "required": [
        "service"
      ],
      "properties": {
        "id": {
          "type": "string"
        },
        "limit": {
          "type": "integer",
          "default": 0
        },
        "login": {
          "type": "string",
          "minLength": 1
        },
        "password": {
          "type": "string"
        },
        "service": {
          "type": "string"
        },
        "token": {
          "type": "string"
        }
      }
    },
    "SearchMoviesResult": {
      "type": "object",
      "required": [
        "id",
        "title"
      ],
      "properties": {
        "description": {
          "type": "string"
        },
        "genres": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "id": {
          "type": "string"
        },
        "poster": {
          "type": "string"
        },
        "preview": {
          "type": "string"
        },
        "rating": {
          "type": "number"
        },
        "seasons": {
          "type": "integer"
        },
        "title": {
          "type": "string"
        },
        "type": {
          "type": "string",
          "enum": [
            "tv-series",
            "film"
          ]
        },
        "year": {
          "type": "integer"
        }
      }
    },
    "SearchTorrentsResult": {
      "type": "object",
      "required": [
        "link",
        "title",
        "size",
        "seeders"
      ],
      "properties": {
        "format": {
          "description": "Формат",
          "type": "string"
        },
        "link": {
          "type": "string"
        },
        "quality": {
          "description": "Качество видео",
          "type": "string",
          "enum": [
            "",
            "480p",
            "720p",
            "1080p",
            "2160p"
          ]
        },
        "rip": {
          "description": "Rip для видео",
          "type": "string"
        },
        "seasons": {
          "description": "Количество сезонов в сериале (если это сериал)",
          "type": "array",
          "items": {
            "type": "integer",
            "minimum": 1
          }
        },
        "seeders": {
          "type": "integer",
          "minimum": 0
        },
        "size": {
          "type": "integer",
          "minimum": 0
        },
        "subtitles": {
          "description": "Коды языков, на которых предоставлены субтитры",
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "title": {
          "type": "string"
        },
        "voice": {
          "description": "Озвучка",
          "type": "string"
        }
      }
    },
    "principal": {
      "type": "object",
      "properties": {
        "canManageAccounts": {
          "type": "boolean"
        },
        "token": {
          "type": "string"
        }
      }
    }
  },
  "securityDefinitions": {
    "key": {
      "type": "apiKey",
      "name": "x-token",
      "in": "header"
    }
  },
  "tags": [
    {
      "description": "Фильмы/сериалы",
      "name": "movies"
    },
    {
      "description": "Торренты",
      "name": "torrents"
    },
    {
      "description": "Администрирование аккаунтов",
      "name": "accounts"
    }
  ]
}`))
}
