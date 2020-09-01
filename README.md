# Game Keys Store
Итоговое задание летней школы Xsolla backend 2020

## Содержание
* [Обзор](#Обзор)
* [Требования](#Требования)
* [Что сделано](#Что-сделано)
* [Установка](#Установка)
* [Как улучшить](#Как-улучшить)

## Обзор
API для покупки и продажи игровых ключей. Продавец может добавить игру и ключи, которые затем будут доступны для покупки. 
Покупатель получает ключ на электронную почту, а на сервер продавца отправляется запрос, в котором содержится информация о покупателе.

## Требования
- Go v1.15
- golang-migrate
- Docker

## Что сделано
- Добавление игр и ключей
- Создание и оплата платежных сеcсий
- Удаление неоплаченых сессий
- Отправка ключа по электронной почте ([link](https://github.com/rdsalakhov/game-keys-store/blob/master/internal/services/paymentService.go#L111))
- Отправка запроса с оповещением продавцу ([link](https://github.com/rdsalakhov/game-keys-store/blob/master/internal/services/notificationService.go#L27))
- В заголовке передается контрольная сумма: md5 хэш тела запроса с солью ([link](https://github.com/rdsalakhov/game-keys-store/blob/master/internal/services/notificationService.go#L68))
- Изменение параметров платформы в config файле ([link](https://github.com/rdsalakhov/game-keys-store/blob/master/configs/config.yml))
- JWT аутентификация по Access и Refresh токенам ([link](https://github.com/rdsalakhov/game-keys-store/blob/master/internal/server/authenti%D1%81ation.go))
- Хранение токенов пользователей в Redis
- Middleware для логирования запросов ([link](https://github.com/rdsalakhov/game-keys-store/blob/master/internal/server/middleware.go))
- docker-compose файл ([link](https://github.com/rdsalakhov/game-keys-store/blob/master/docker-compose.yml))
- Управление версиями базы данных с помощью миграций
- Swagger-спецификация ([link](https://app.swaggerhub.com/apis/rs-org/game-keys-store/1.0.0#/free))
- Деплой на Heroku ([link](https://game-keys-store.herokuapp.com))

## Установка
1. ```docker-compose up```
2. ```migrate -database "mysql://root:root@tcp(localhost:3306)/game_keys_db_mysql" -path ./migrations/ up```

## Как улучшить
- Добавить модульные тесты
- Добавить CI/CD pipeline
- Добавить сохранение логов в базу данных

