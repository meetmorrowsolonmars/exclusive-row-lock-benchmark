# Benchmark

# Запуск

```shell
docker-compose up -d

go install github.com/pressly/goose/v3/cmd/goose@latest
goose --dir migrations postgres "user=benchmark password=benchmark dbname=benchmark sslmode=disable host=localhost port=6543" up

go test -v ./benchmark/mutex
```

# Результаты

## Прогон 1

1. Количество строк в базе: 20
2. Изначальная сумма баланса: 100
3. Количество запусков: 12
4. Множитель запросов на запуск: 10

Ожидается, что будет выполняться 10 одновременных запросов на изменение каждой строки за запуск. В сумме 120 запросов
на каждую строку.

```
=== RUN   Test_Benchmark
=== RUN   Test_Benchmark/Mutex
    mutex_test.go:104: [Mutex] Total time: 57.628562666s
    mutex_test.go:105: [Mutex] Average time: 4800ms
=== RUN   Test_Benchmark/OptimisticLock
    mutex_test.go:104: [OptimisticLock] Total time: 10.676716208s
    mutex_test.go:105: [OptimisticLock] Average time: 887ms
--- PASS: Test_Benchmark (68.37s)
    --- PASS: Test_Benchmark/Mutex (57.63s)
    --- PASS: Test_Benchmark/OptimisticLock (10.68s)
PASS
ok      benchmark/benchmark/mutex       68.486s
```
