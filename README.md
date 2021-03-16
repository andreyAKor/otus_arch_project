# Проектная работа по курсу "Архитектура и шаблоны проектирования"

### Запуск сервиса
Поднять docker-контейнер с Postgress:
```shell script
docker run -d \
--name pg-otus-arch \
-e POSTGRES_USER=bid \
-e POSTGRES_PASSWORD=bid123 \
-e POSTGRES_DB=bid \
-e PGDATA=/var/lib/postgresql/data/pgdata \
-v '/home/andreyakor/Документы/OTUS/Архитектура и шаблоны проектирования/psqldata':/var/lib/postgresql/data \
-p 5432:5432 \
postgres
```

Выполняем сборку сервисов:
```shell script
$ make build
```

Открываем три окна терминала и запускаем сервисы:
```shell script
./bin/gateway --config=./configs/gateway.yml
```
```shell script
./bin/bid --config=./configs/bid.yml
```
```shell script
./bin/pedding --config=./configs/pedding.yml
```

Сервис поднимается на локальном хосте на порту 6080.

Отдельно настраиваем Apache и Nginx для работы сайты из директории `./website`. Примеры файлов конфигов для Apache и Nginx находятся в директории `./website/configs/`.

---

### Прочие операции с make
- `make install` - устанавливает все необходимые модули через go mod
- `make generate` - go-генерация небходимых для проекта пакетов
- `make lint` - прогонка проекта линтером
- `make build` - сборка сервиса
- `make run` - запуск сервсиа в docker-контейнере через docker-compose
- `make test` - запуст юнит-тестов
