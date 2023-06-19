# Сервис бэкенд для системы аукциона

Основные параметры содержатся в файле properties.json (переименовать properties_example.json в properties.json, если не существует).

## База данных

Система расчитана на работу с Mysql/Mariadb. Для других СУБД необходимо дополнительное подключение драйвера и доп.перепись обращений к базе данных Использовался возврат engine...

Необходимые для работы таблицы автоматически создаются в базе данных, думать об этом не нужно. Данные автоматически не вносятся.

Базу данных необходимо поднять отдельно


Возможная конфигурация docker-compose для нее:
```yml
 mariadb:
    image: mariadb:latest
    environment:
      MYSQL_ROOT_PASSWORD: karim
      MYSQL_DATABASE: database
      MYSQL_USER: karim
      MYSQL_PASSWORD: karim
    logging:
      driver: syslog
      options:
        tag: "{{.DaemonName}}(image={{.ImageName}};name={{.Name}};id={{.ID}})"
    restart: always
    volumes:
     - ./data:/var/lib/mysql
```



## Сборка

Для линукс (сборка на системе windows)

```sh
 $Env:GOOS = "linux"; $Env:GOARCH = "amd64"; go build -ldflags="-s -w" -o go_aukt_backend
```

Для линукс (сборка на системе linux)

```sh
 GOOS="linux" GOARCH="amd64" go build -ldflags="-s -w" -o go_aukt_backend
```
