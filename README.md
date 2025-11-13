## Сервис для сокращения URL-адресов

Необходимо разработать сервис url-shortener для сокращения URL-адресов по примеру https://tinyurl.com/.
Приложение должно быть реализовано в виде HTTP сервера, реализующее REST API. Сервер должен реализовывать
7 методов и их логику:

#### Проверка работоспособности
- Эндпоинт - GET /api/health
- Статус ответа 200 если сервис работает исправно

##### Пример
```bash
$ curl -i -X GET 'http://localhost:8081/api/health'

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
$ curl -i -X POST 'http://localhost:8081/api/users' \
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
$ curl -i -X POST 'http://localhost:8081/api/login' \
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
		- alias - сокращенный путь (необязательный параметр)
- Статус ответа 201 если новый URL-адреса создан успешно. Ответ должен содержать в теле JSON-объект:
	- alias – сокращенный путь

##### Пример
```bash
$ curl --user Bob:qwerty -i -X POST 'http://localhost:8081/api/data/shorten' \
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
$ curl -i -X GET 'http://localhost:8081/zn9edcu'

<a href="https://en.wikipedia.org/wiki/Systems_design">Found</a>.
```

#### Удаление сокращенного URL-адреса
- Эндпоинт: DELETE /api/{alias}
- Статус ответа 200 если URL-адреса c 'alias' удален успешно

##### Пример
```bash
$ curl --user Bob:qwerty -i -X DELETE 'http://localhost:8081/api/zn9edcu'

{
	"status":"OK"
}
```

#### Получение списка всех сокращенных URL-адресов пользователя
- Эндпоинт: GET /api/urls
- Ответ должен содержать общее количество сокращенных URL-адресов пользователя (total) и в теле массив JSON-объектов с информацией о сокращенных URL-адресах пользователя. Каждый объект содержит параметры:
	- url - исходный, полный URL-адрес
	- alias - сокращенный путь
	- count - количества переходов по сокращенному URL-адресу
	- created_at - дата и время создания сокращенного URL-адреса
- Статус ответа 200 если список получен успешно.

##### Пример
```bash
$ curl --user Bob:qwerty -i -X GET 'http://localhost:8081/api/urls'

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
$ make build
...............
$ make up
...............
```

### Todo list
- Добавить метод "Проверка учетных данных пользователя".
- Добавить метод "Проверка доступности алиаса".
- Вынести encode url в отдельный пакет и добавить модульные тесты.
- Добавить модульные тесты для всех методов (обработчиков).
- Добавить интеграционные тесты.
- url -> urls
- 8081 -> 8080

### Полезные ссылки
- [Пишем REST API сервис на Go - УЛЬТИМАТИВНЫЙ гайд](https://www.youtube.com/watch?v=rCJvW2xgnk0)
- [LRU cache](https://github.com/hashicorp/golang-lru)
- [bombardier](https://github.com/codesenberg/bombardier)