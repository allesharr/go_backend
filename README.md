# Сервис бэкенд для системы СКУД

Сервис бэкенд для системы СКУД. Основные параметры в файле properties.json (переименовать properties_example.json в properties.json, если не существует).

## Сбор данных для БД

На сервере SKUD в папке `/home/admin/paradox_reader_dotnet` лежит приложение (`paradox-csv`) импорта 2 файлов (события `/home/admin/pLogData.DB` и пользователи `/home/admin/pList.DB`) формата Paradox-DB в CSV. Запускается приложение-конвертер как

```sh
LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib ./paradox-csv ../pLogData.db out

LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib ./paradox-csv ../pList.DB out
```

после выполняется конвертация из CP1251 в UTF-8

После этого управление передается в `/home/admin/paradox_reader_dotnet/paradox_reader_dotnet`, который заносит и CSV данные в БД: таблицу пользователей очищает и заносит заново, а в таблицу событий добавляются только изменения.

Все это организовано в виде одного скрипта `/home/admin/paradox_reader_dotnet/run_paradox_reader.sh`, который запускается по системному крону.

Таблица в свою очередь уже используется сервисом бэкендом.

## Сборка

Для линукс (сборка на системе windows)

```sh
 $Env:GOOS = "linux"; $Env:GOARCH = "amd64"; go build -ldflags="-s -w" -o go_skud_backend
```

Для линукс (сборка на системе linux)

```sh
 GOOS="linux" GOARCH="amd64" go build -ldflags="-s -w" -o go_skud_backend
```
