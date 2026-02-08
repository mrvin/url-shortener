### Описание методов API

#### Фронтенд
- Эндпоинт: GET /index.html

#### Swagger UI
- Эндпоинт: GET /swagger.html

#### Получение информации о приложении
- Эндпоинт - GET /api/info
- Статус ответа 200

##### Пример запроса
```bash
curl -i -X GET 'http://localhost:8080/api/info'
```
##### Пример ответа
```json
{
  "tag": "v0.0.2",
  "hash": "8f32719a5b9e74817ea7c62765545438e39377dd",
  "date": "2025-11-21.15:06:09"
}
```

#### Проверка работоспособности
- Эндпоинт - GET /api/health
- Статус ответа 200 если сервис работает исправно

##### Пример запроса
```bash
curl -i -X GET 'http://localhost:8080/api/health'
```
##### Пример ответа
```json
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

##### Пример запроса
```bash
curl -i -X POST 'http://localhost:8080/api/users' \
-H "Content-Type: application/json" \
-d '{
	"username":"Bob",
	"password":"qwerty"
}'
```
##### Пример ответа
```json
{
  "status": "OK"
}
```

#### Проверка учетных данных пользователя
- Эндпоинт - POST /api/users/login
- Параметры запроса:
	- JSON-объект в теле запроса с параметрами:
		- username – имя пользователя
		- password – пароль
- Статус ответа 200 если пользователь прошел проверку
- Статус ответа 401 если пользователь не прошел проверку

##### Пример запроса
```bash
curl -i -X POST 'http://localhost:8080/api/users/login' \
-H "Content-Type: application/json" \
-d '{
	"username":"Bob",
	"password":"qwerty"
}'
```
##### Пример ответа
```json
{
  "status": "OK"
}
```

#### Создание нового сокращенного URL-адреса
- Эндпоинт: POST /api/urls
- Параметры запроса:
	- JSON-объект в теле запроса с параметрами:
		- url – исходный, полный URL-адрес
		- alias - сокращенный путь
- Статус ответа 201 если новый URL-адреса создан успешно.

##### Пример запроса
```bash
curl --user Bob:qwerty -i -X POST 'http://localhost:8080/api/urls' \
-H "Content-Type: application/json" \
-d '{
	"url":"https://en.wikipedia.org/wiki/Systems_design",
	"alias":"zn9edcu"
}'
```
##### Пример ответа
```json
{
  "status": "OK"
}
```

#### Перенаправление URL-адреса
- Эндпоинт: GET /{alias}
- Статус ответа 302 (Перенаправление) если alias существует
- Статус ответа 404 если alias не найден

##### Пример запроса
```bash
curl -i -X GET 'http://localhost:8080/zn9edcu'
```
##### Пример ответа
```html
<a href="https://en.wikipedia.org/wiki/Systems_design">Found</a>.
```

#### Проверка доступности алиаса
- Эндпоинт: GET /api/urls/check/{alias}
- Статус ответа 200

##### Пример запроса
```bash
curl -i -X GET 'http://localhost:8080/api/urls/check/zn9edcu'
```
##### Пример ответа
```json
{
  "exists": true,
  "status": "OK"
}
```

#### Удаление сокращенного URL-адреса
- Эндпоинт: DELETE /api/urls/{alias}
- Статус ответа 200 если URL-адреса c 'alias' удален успешно

##### Пример запроса
```bash
curl --user Bob:qwerty -i -X DELETE 'http://localhost:8080/api/urls/zn9edcu'
```
##### Пример ответа
```json
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

##### Пример запроса
```bash
curl --user Bob:qwerty -i -X GET 'http://localhost:8080/api/urls?limit=10&offset=0'
```
##### Пример ответа
```json
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