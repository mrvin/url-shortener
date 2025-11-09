## Сервис для сокращения URL-адресов

Необходимо разработать сервис url-shortener для сокращения URL-адресов по примеру https://tinyurl.com/.
Приложение должно быть реализовано в виде HTTP сервера, реализующее REST API. Сервер должен реализовывать
6 методов и их логику:

#### Проверка работоспособности
- Эндпоинт - GET /health
- Статус ответа 200 если сервис работает исправно

##### Пример
```bash
$ curl -i -X GET 'http://localhost:8081/health'

{
	"status":"OK"
}
```

#### Регистрация пользователя
- Эндпоинт - POST /users
- Параметры запроса:
	- JSON-объект в теле запроса с параметрами:
		- user_name – имя пользователя
		- password – пароль
- Статус ответа 201 если пользователь создан успешно

##### Пример
```bash
$ curl -i -X POST 'http://localhost:8081/users' \
-H "Content-Type: application/json" \
-d '{
	"user_name":"Bob",
	"password":"qwerty"
}'

{
  	"status": "OK"
}
```

#### Создание нового сокращенного URL-адреса
- Эндпоинт: POST /data/shorten
- Параметры запроса:
	- JSON-объект в теле запроса с параметрами:
		- url – исходный, полный URL-адрес
		- alias - сокращенный путь (необязательный параметр)
- Статус ответа 201 если новый URL-адреса создан успешно. Ответ должен содержать в теле JSON-объект:
	- alias – сокращенный путь

##### Пример
```bash
$ curl --user Bob:qwerty -i -X POST 'http://localhost:8081/data/shorten' \
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
- Эндпоинт: DELETE /{alias}
- Статус ответа 200 если URL-адреса c 'alias' удален успешно

##### Пример
```bash
$ curl --user Bob:qwerty -i -X DELETE 'http://localhost:8081/zn9edcu'

{
	"status":"OK"
}
```

#### Получение количества переходов по сокращенному URL-адресу
- Эндпоинт: GET /statistics/{alias}
- Статус ответа 200 если количества переходов получено успешно. Ответ должен содержать в теле JSON-объект. Объект содержит параметры:
	- count – количества переходов по сокращенному URL-адресу

##### Пример
```bash
$ curl --user Bob:qwerty -i -X GET 'http://localhost:8081/statistics/zn9edcu'

{
	"count":0,
	"status":"OK"
}
```

### Сборка и запуск приложения в Docker Compose
```shell script
$ make build
...............
$ make up
...............
```

### Полезные ссылки
- [Пишем REST API сервис на Go - УЛЬТИМАТИВНЫЙ гайд](https://www.youtube.com/watch?v=rCJvW2xgnk0)
- [LRU cache](https://github.com/hashicorp/golang-lru)
- [bombardier](https://github.com/codesenberg/bombardier)