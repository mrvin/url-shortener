## Сервис для сокращения URL-адресов

Необходимо разработать сервис url-shortener для сокращения URL-адресов по примеру https://tinyurl.com/.
Приложение должно быть реализовано в виде HTTP сервера, реализующее REST API. Сервер должен реализовывать 9 методов и их логику:

#### Фронтенд
- Эндпоинт: GET /static/index.html

#### Получение информации о приложении
- Эндпоинт - GET /api/info
- Статус ответа 200

##### Пример
```bash
$ curl -i -X GET 'http://localhost:8080/api/info'

{
  "tag": "v0.0.2",
  "hash": "8f32719a5b9e74817ea7c62765545438e39377dd",
  "date": "2025-11-21.15:06:09"
}
```

#### Проверка работоспособности
- Эндпоинт - GET /api/health
- Статус ответа 200 если сервис работает исправно

##### Пример
```bash
$ curl -i -X GET 'http://localhost:8080/api/health'

{
	"status":"OK"
}
```

#### Регистрация пользователя
- Эндпоинт - POST /api/users
- Параметры запроса:
	- JSON-объект в теле запроса с параметрами:
		- username – имя пользователя
		- password – пароль
- Статус ответа 201 если пользователь создан успешно

##### Пример
```bash
$ curl -i -X POST 'http://localhost:8080/api/users' \
-H "Content-Type: application/json" \
-d '{
	"username":"Bob",
	"password":"qwerty"
}'

{
  	"status": "OK"
}
```

#### Проверка учетных данных пользователя
- Эндпоинт - POST /api/login
- Параметры запроса:
	- JSON-объект в теле запроса с параметрами:
		- username – имя пользователя
		- password – пароль
- Статус ответа 200 если пользователь прошел проверку
- Статус ответа 401 если пользователь не прошел проверку

##### Пример
```bash
$ curl -i -X POST 'http://localhost:8080/api/login' \
-H "Content-Type: application/json" \
-d '{
	"username":"Bob",
	"password":"qwerty"
}'

{
  	"status": "OK"
}
```

#### Создание нового сокращенного URL-адреса
- Эндпоинт: POST /api/data/shorten
- Параметры запроса:
	- JSON-объект в теле запроса с параметрами:
		- url – исходный, полный URL-адрес
		- alias - сокращенный путь
- Статус ответа 201 если новый URL-адреса создан успешно. Ответ должен содержать в теле JSON-объект:
	- alias – сокращенный путь

##### Пример
```bash
$ curl --user Bob:qwerty -i -X POST 'http://localhost:8080/api/data/shorten' \
-H "Content-Type: application/json" \
-d '{
	"url":"https://en.wikipedia.org/wiki/Systems_design",
	"alias":"zn9edcu"
}'

{
	"alias":"zn9edcu",
	"status":"OK"
}
```

#### Перенаправление URL-адреса
- Эндпоинт: GET /{alias}
- Статус ответа 302 (Перенаправление) если alias существует
- Статус ответа 404 если alias не найден

##### Пример
```bash
$ curl -i -X GET 'http://localhost:8080/zn9edcu'

<a href="https://en.wikipedia.org/wiki/Systems_design">Found</a>.
```

#### Проверка доступности алиаса
- Эндпоинт: GET /api/check/{alias}
- Статус ответа 200

##### Пример
```bash
$ curl -i -X GET 'http://localhost:8080/api/check/zn9edcu'

{
  "exists": true,
  "status": "OK"
}
```

#### Удаление сокращенного URL-адреса
- Эндпоинт: DELETE /api/{alias}
- Статус ответа 200 если URL-адреса c 'alias' удален успешно

##### Пример
```bash
$ curl --user Bob:qwerty -i -X DELETE 'http://localhost:8080/api/zn9edcu'

{
	"status":"OK"
}
```

#### Получение списка всех сокращенных URL-адресов пользователя
- Эндпоинт: GET /api/urls
- Параметры запроса:
	- limit – количество url-адресов в ответе (по умолчанию 100)
	- offset - смищение от начала (по умолчанию 0)
- Ответ должен содержать общее количество сокращенных URL-адресов пользователя (total) и в теле массив JSON-объектов с информацией о сокращенных URL-адресах пользователя. Каждый объект содержит параметры:
	- url - исходный, полный URL-адрес
	- alias - сокращенный путь
	- count - количества переходов по сокращенному URL-адресу
	- created_at - дата и время создания сокращенного URL-адреса
- Статус ответа 200 если список получен успешно.

##### Пример
```bash
$ curl --user Bob:qwerty -i -X GET 'http://localhost:8080/api/urls?limit=10&offset=0'

{
  "urls": [
    {
      "url": "https://en.wikipedia.org/wiki/Systems_design",
      "alias": "zn9edcu",
      "count": "24812",
      "created_at": "2025-09-25T16:18:38.384975Z"
    }
  ],
  "total": 1,
  "status": "OK"
}
```

### Сборка и запуск приложения в Docker Compose
```shell script
$ make run-compose
...............
```

### Todo list
- Добавить модульные тесты для всех методов/обработчиков(CheckAlias, Login).
- Добавить интеграционные тесты.
- Добавить документацию OpenAPI.
- Добавить более строгую валидацию логина и пароля при регистрации.
- Добавить в конфигурирование ttl для кэша.

### Полезные ссылки
- [Пишем REST API сервис на Go - УЛЬТИМАТИВНЫЙ гайд](https://www.youtube.com/watch?v=rCJvW2xgnk0)
- [LRU cache](https://github.com/hashicorp/golang-lru)
- [bombardier](https://github.com/codesenberg/bombardier)