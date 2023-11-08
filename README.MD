# Пример проекта и кода для обращения к API GitHub

## Цель
- Получить информацию об участниках каждого репозитория в организации `microsoft`

## Файлы в проекте
- `api/orgs/repos/get.go` - описание структуры Репозитория и подготовка запроса к API
- `api/repos/contributors/get.go` - описание структуры Участника
- `api/users/get.go` - описание структуры Пользователя  
&nbsp;
- `api/init.go` - инициализация пакета для доступа к API
- `api/option.go` - определение структуры паттерна функциональных опций  
&nbsp;
- `api/repos.go` - определение метода для доступа к Репозиторию по указанной организации
- `api/contributors.go` - определение метода для доступа к Участникам по указанному Репозиторию
- `api/user.go` - определение метода для доступа к Пользователю  
&nbsp;
- `cmd/repos/main.go` - отправная точка для запуска проекта  
&nbsp;
- `internal/repos/main.go` - внутренняя логика проекта
- `internal/env/env.go` - определение переменных к проекту из файла


## Хранение токена
Для того чтобы предоставить токен проекту необходимо создать файл `.git.env` в папке `secret`  
Со следующим содержимым:
```toml
GIT_TOKEN="<token>"
```

## Запуск проекта
```bash
go run cmd/repos/main.go
```

## Сборка проекта и запуск собранного проекта
> Почитать про указанные переменные перед командой `go build`
```bash
GOGC=off GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C "cmd/repos" -ldflags "-s -w" -o ../../repos
./repos
```

## Прогон тестов
```bash
go -C internal/repos/ test -bench=. -benchmem
```
Вывод результатов
```bash
RunSimple elapsed time: 15.316222158s
goos: linux
goarch: amd64
pkg: gitmic/internal/repos
cpu: Intel(R) Core(TM) i5-8400 CPU @ 2.80GHz
BenchmarkRunSimple-6                   1        15316264658 ns/op        8737640 B/op      52356 allocs/op
RunConcurrency elapsed time: 1.114834146s
BenchmarkRunConcurrency-6              1        1114883946 ns/op          505480 B/op       1284 allocs/op
PASS
ok      gitmic/internal/repos   16.438s
```

## 
```
go run cmd/repos/main.go -cpuprofile cpu.prof -memprofile mem.prof
go tool pprof mem.prof
top10 -cum
list gitmic/api.doRequest
```