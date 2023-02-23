# Exclusive row lock benchmark

## Запуск

```shell
docker-compose up -d

go install github.com/pressly/goose/v3/cmd/goose@latest
goose --dir migrations postgres "user=benchmark password=benchmark dbname=benchmark sslmode=disable host=localhost port=6543" up

go test -v ./benchmark/mutex
```

## Результаты

### Прогон 1

* Количество строк в базе: 20
* Изначальная сумма баланса: 100
* Количество запусков: 12
* Множитель запросов на запуск: 10

Ожидается, что будет выполняться 10 одновременных запросов на изменение каждой строки за запуск. В сумме 120 запросов
на каждую строку.

```
=== RUN   Test_Benchmark
=== RUN   Test_Benchmark/Mutex
    mutex_test.go:105: [Mutex] Total time: 53.511216s
    mutex_test.go:106: [Mutex] Average time: 2241ms
=== RUN   Test_Benchmark/OptimisticLock
    mutex_test.go:105: [OptimisticLock] Total time: 9.803322s
    mutex_test.go:106: [OptimisticLock] Average time: 423ms
--- PASS: Test_Benchmark (63.41s)
    --- PASS: Test_Benchmark/Mutex (53.51s)
    --- PASS: Test_Benchmark/OptimisticLock (9.80s)
PASS
ok      benchmark/benchmark/mutex       63.664s
```

### Прогон 2

* Количество строк в базе: 20
* Изначальная сумма баланса: 100
* Количество запусков: 120
* Множитель запросов на запуск: 1

Ожидается, что будет выполняться 1 одновременный запрос к каждой строке. Всего 120 запросов на строку.

```
=== RUN   Test_Benchmark
=== RUN   Test_Benchmark/Mutex
    mutex_test.go:105: [Mutex] Total time: 1m0.295035458s
    mutex_test.go:106: [Mutex] Average time: 262ms
=== RUN   Test_Benchmark/OptimisticLock
    mutex_test.go:105: [OptimisticLock] Total time: 11.291624959s
    mutex_test.go:106: [OptimisticLock] Average time: 63ms
--- PASS: Test_Benchmark (71.68s)
    --- PASS: Test_Benchmark/Mutex (60.30s)
    --- PASS: Test_Benchmark/OptimisticLock (11.29s)
PASS
ok      benchmark/benchmark/mutex       71.928s
```

### Прогон 3

* Количество строк в базе: 20
* Изначальная сумма баланса: 200
* Количество запусков: 50
* Множитель запросов на запуск: 5

Ожидается, что будет выполняться 5 одновременных запросов к каждой строке. Всего 250 запросов на строку.

```
=== RUN   Test_Benchmark
=== RUN   Test_Benchmark/Mutex
    mutex_test.go:105: [Mutex] Total time: 1m55.548665s
    mutex_test.go:106: [Mutex] Average time: 1166ms
=== RUN   Test_Benchmark/OptimisticLock
    mutex_test.go:105: [OptimisticLock] Total time: 21.775219417s
    mutex_test.go:106: [OptimisticLock] Average time: 234ms
--- PASS: Test_Benchmark (137.42s)
    --- PASS: Test_Benchmark/Mutex (115.55s)
    --- PASS: Test_Benchmark/OptimisticLock (21.78s)
PASS
ok      benchmark/benchmark/mutex       137.757s
```

### Прогон 4

* Количество строк в базе: 20
* Изначальная сумма баланса: 100
* Количество запусков: 1
* Множитель запросов на запуск: 120

```
=== RUN   Test_Benchmark
=== RUN   Test_Benchmark/Mutex
    mutex_test.go:105: [Mutex] Total time: 55.978551167s
    mutex_test.go:106: [Mutex] Average time: 29829ms
=== RUN   Test_Benchmark/OptimisticLock
    mutex_test.go:105: [OptimisticLock] Total time: 9.991789958s
    mutex_test.go:106: [OptimisticLock] Average time: 5027ms
--- PASS: Test_Benchmark (66.03s)
    --- PASS: Test_Benchmark/Mutex (55.98s)
    --- PASS: Test_Benchmark/OptimisticLock (9.99s)
PASS
ok      benchmark/benchmark/mutex       66.281s
```
