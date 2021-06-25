# Тестовое задание Advertising

<!-- ToC start -->
# Содержание

1. [Описание задачи](#Описание-задачи)
1. [Реализация](#Реализация)
1. [Архитектура](#Архитектура)
1. [Endpoints](#Endpoints)
1. [Запуск](#Запуск)
1. [Тестирование](#Тестирование)
1. [Документация](#Документация)
1. [Нагрузочное тестирование](#Нагрузочное-тестирование)
1. [Примеры](#Примеры)
<!-- ToC end -->

# Описание задачи
Цель задания – разработать чат-сервер, предоставляющий HTTP API для работы с чатами и сообщениями пользователя.

## Основные сущности

Ниже перечислены основные сущности, которыми должен оперировать сервер.

### User

Пользователь приложения. Имеет следующие свойства:

* **id** - уникальный идентификатор пользователя
* **username** - уникальное имя пользователя
* **created_at** - время создания пользователя

### Chat

Отдельный чат. Имеет следующие свойства:

* **id** - уникальный идентификатор чата
* **name** - уникальное имя чата
* **users** - список пользователей в чате, отношение многие-ко-многим
* **created_at** - время создания

### Message

Сообщение в чате. Имеет следующие свойства:

* **id** - уникальный идентификатор сообщения
* **chat** - ссылка на идентификатор чата, в который было отправлено сообщение
* **author** - ссылка на идентификатор отправителя сообщения, отношение многие-к-одному
* **text** - текст отправленного сообщения
* **created_at** - время создания

## Основные API методы

Методы обрабатывают HTTP POST запросы c телом, содержащим все необходимые параметры в JSON.

### Добавить нового пользователя

Принимает в качестве полей: имя пользователя.

Ответ: `id` созданного пользователя или HTTP-код ошибки.

### Создать новый чат между пользователями

Принимает в качестве полей: имя чата и массив id пользователей в чате.

Ответ: `id` созданного чата или HTTP-код ошибки.

Количество пользователей не ограничено.

### Отправить сообщение в чат от лица пользователя

Принимает в качестве полей: id чата, id автора сообщения и текст сообщения.

Ответ: `id` созданного сообщения или HTTP-код ошибки.

### Получить список чатов конкретного пользователя

Принимает в качестве полей: id пользователя.

Ответ: cписок всех чатов со всеми полями, отсортированный по времени создания последнего сообщения в чате (от позднего к раннему). Или HTTP-код ошибки.

### Получить список сообщений в конкретном чате

Принимает в качестве полей: id чата.

Ответ: список всех сообщений чата со всеми полями, отсортированный по времени создания сообщения (от раннего к позднему). Или HTTP-код ошибки.

# Реализация

- Следование дизайну REST API.
- Подход "Чистой Архитектуры" и техника внедрения зависимости.
- Работа с фреймворком [echo](https://echo.labstack.com/).
- Работа с БД Postgres с использованием библиотеки [sqlx](https://github.com/jmoiron/sqlx) и написанием SQL запросов.
- Конфигурация приложения - библиотека [viper](https://github.com/spf13/viper).
- Реализация Graceful Shutdown.
- Запуск из Docker.
- Юнит-тестирование бизнес-логики и взаимодействия с БД классическим способом и с помощью моков - библиотеки [testify](https://github.com/stretchr/testify), [mock](https://github.com/golang/mock).
- Сквозное (E2E) тестирование - BDD фреймворк [goconvey](https://github.com/smartystreets/goconvey).
- Проверка кода на соответствие стандартам с помощью линтера - утилита [golangci-lint](https://github.com/golangci/golangci-lint)
- Автоматическое создание документации с помощью Swagger 2.0 - библиотека [echo-swagger](https://github.com/swaggo/echo-swagger).
- Непрерывная интеграция - сборка приложения, проверка линтером и запуск тестов в Github action.

**Структура проекта:**
```
.
├── pkg
│   ├── error_message  // сообщения об ошибках
│   ├── model          // основные структуры
│   ├── delivery       // обработчики запросов
│   ├── service        // бизнес-логика
│   └── repository     // взаимодействие с БД
├── cmd                // точка входа в приложение
├── migrations         // SQL файлы с миграциями
├── tet_scripts        // SQL файлы с тестовыми данными
├── configs            // файлы конфигурации
├── test               // инициализация тестовой БД
└── e2e_test.go        // сквозной тест
```

# Архитектура
![Схема](https://github.com/VolkovEgor/sellerx-task/blob/develop/docs/img/architecture.jpg)

Приложение имеет 3 основных слоя, реализованных в отдельных пакетах.

- Repository - слой взаимодействия с БД. Методы этого слоя принимают данные от Service и выполняют запросы к БД.
- Service - слой бизнес-логики. Методы этого слоя принимают данные от Handler и применяют к ним бизнес-правила для достижения цели варианта использования.
- Delivery - слой обработчиков запросов. Содержит методы-обработчики для endpoints.

Пакет Model содержит структуры сущностей, используемых остальными слоями.

# Запуск

```
make build
make run
```

Если приложение запускается впервые, необходимо применить миграции к базе данных:

```
make migrate_up
```

Для миграций используется [golang-migrate/migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation).

# Тестирование

Локальный запуск тестов:
```
make run_test
```

Для локального запуска тестов необходимо создать тестовую БД. Это можно сделать следующей командой (необходима утилита psql):
```
make create_test_db
```
# Документация

Для просмотра документации Swagger необходимо запустить приложение и перейти по ссылке [http://127.0.0.1:9000/swagger/index.html](http://127.0.0.1:9000/swagger/index.html) 

# Нагрузочное тестирование


# Примеры

### Создание пользователя

**Запрос:**
```
$ curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"username": "user_1"}' \
  http://localhost:9000/users/add
```
**Тело ответа:**
```
[
    {
        "id": 1,
        "title": "Advert 1",
        "main_photo": "link1",
        "price": 10000
    },
    {
        "id": 2,
        "title": "Advert 2",
        "main_photo": "link1",
        "price": 60000
    },
    {
        "id": 3,
        "title": "Advert 3",
        "main_photo": "link1",
        "price": 30000
    }
]
```

### Создание чата

**Запрос:**
```
$ curl GET localhost:9000/api/adverts?page=1&sort=price_asc
```
**Тело ответа:**
```
{
    "id":"74e6f204-4a47-40b9-aad1-8ea40887867f"
}
```

### 3. GET для _page=1_ и _sort_=date_desc

**Запрос:**
```
$ curl GET localhost:9000/api/adverts?page=1&sort=date_desc
```
**Тело ответа:**
```
[
    {
        "id": 3,
        "title": "Advert 3",
        "main_photo": "link1",
        "price": 30000
    },
    {
        "id": 2,
        "title": "Advert 2",
        "main_photo": "link1",
        "price": 60000
    },
    {
        "id": 1,
        "title": "Advert 1",
        "main_photo": "link1",
        "price": 10000
    }
]
```

## Получение конкретного объявления

### 1. GET для _id=1_ 

**Запрос:**
```
$ curl GET localhost:9000/api/adverts/1
```
**Тело ответа:**
```
{
    "id": 1,
    "title": "Advert 1",
    "photos": [
        "link1"
    ],
    "price": 10000
}
```
### 2. GET для _id=1_ и _fields_=true

**Запрос:**
```
$ curl GET localhost:9000/api/adverts/1?fields=true
```
**Тело ответа:**
```
{
    "id": 1,
    "title": "Advert 1",
    "description": "Description 1",
    "photos": [
        "link1",
        "link2",
        "link3"
    ],
    "price": 10000
}
```

## Создание объявления

**Запрос:**
```
$ curl --location --request POST 'localhost:9000/api/adverts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "New advert",
    "description": "Description of new advert",
    "photos": ["link1", "link2"],
    "price": 400000
}'
```
**Тело ответа:**
```
{
    "advert_id": 4
}
```