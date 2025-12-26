## Сервис для сокращения URL-адресов

Сервис url-shortener для сокращения URL-адресов по примеру https://tinyurl.com/. Этот сервис позволяет пользователям сокращать длинные URL-адреса до более удобного формата, облегчая их использование и распространение.
Сервис реализовано в виде HTTP сервера, реализующего REST API.
[Описание методов API](./docs/API.md)

### Основные Функции
- Сокращение URL: пользователи могут преобразовывать длинные URL в короткие ссылки, которые легче обменивать и использовать.

- Перенаправление по коротким ссылкам: каждая сокращенная ссылка перенаправляет пользователя на оригинальный URL.

- Хранение и управление ссылками: сервис предоставляет интерфейс для управления сокращенными ссылками.

### Технологии
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![openapi initiative](https://img.shields.io/badge/openapiinitiative-%23000000.svg?style=for-the-badge&logo=openapiinitiative&logoColor=white)
![Git](https://img.shields.io/badge/git-%23F05033.svg?style=for-the-badge&logo=git&logoColor=white)
![GitHub](https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)
![Linux](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)


### Сборка и запуск приложения в Docker Compose
```bash
make run
```

### Todo list
- Добавить интеграционные тесты.
- Добавить более строгую валидацию логина и пароля при регистрации.
- Добавить в конфигурирование ttl для кэша.
- Уточнить ошибки валидации.
- Добавить структуру базы данных в описание.
- Подумать как лучще удалять url, что бы кэш остовался консистентным.
- Добавть удаление и обновление пользователя.
- Добавить обновление url. 

### Полезные ссылки
- [Пишем REST API сервис на Go - УЛЬТИМАТИВНЫЙ гайд](https://www.youtube.com/watch?v=rCJvW2xgnk0)
- [LRU cache](https://github.com/hashicorp/golang-lru)
- [bombardier](https://github.com/codesenberg/bombardier)