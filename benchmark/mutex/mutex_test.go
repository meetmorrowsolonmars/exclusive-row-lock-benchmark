package mutex

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

func Test_Benchmark(t *testing.T) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig("postgresql://benchmark:benchmark@localhost:6543/benchmark")
	if err != nil {
		t.Fatal(err)
	}

	config.MaxConnLifetime = 60
	config.MaxConnIdleTime = 30
	config.MaxConns = 10
	config.MinConns = 5

	connect, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}

	cases := []TestCase{
		{
			locker: &RepositoryWithMutex{
				db: connect,
				BaseRepository: BaseRepository{
					table: "mutex",
					db:    connect,
				}},
			name:             "Mutex",
			count:            20,
			amount:           100,
			multiplier:       10,
			numberOfLaunches: 12,
		},
		{
			locker: &RepositoryWithOptimisticLock{
				db: connect,
				BaseRepository: BaseRepository{
					table: "optimistic_lock",
					db:    connect,
				},
			},
			name:             "OptimisticLock",
			count:            20,
			amount:           100,
			multiplier:       10,
			numberOfLaunches: 12,
		},
	}

	for _, test := range cases {
		test := test
		// Вставляем балансы.
		ids, testErr := test.locker.InsertTestData(ctx, test.count, test.amount)
		if testErr != nil {
			t.Fatal(testErr)
		}

		t.Run(test.name, func(t *testing.T) {
			milliseconds := int64(0)
			start := time.Now()

			for i := 0; i < test.numberOfLaunches; i++ {
				wg := sync.WaitGroup{}
				wg.Add(len(ids) * test.multiplier)

				launchStart := time.Now()

				for j := 0; j < len(ids)*test.multiplier; j++ {
					go func(id int) {
						defer wg.Done()

						lockErr := test.locker.LockBalance(ctx, id, 1)

						switch e := lockErr.(type) {
						case SkipError:
						case error:
							t.Errorf("Lock balance error: %s", e)
							t.Fail()
						}
					}(ids[j%len(ids)])
				}

				wg.Wait()
				milliseconds += time.Since(launchStart).Milliseconds()
			}

			count, countErr := test.locker.GetCountNegativeBalances(ctx)

			assert.Nil(t, countErr)
			assert.Equal(t, 0, count)

			t.Logf("[%s] Total time: %s", test.name, time.Since(start))
			t.Logf("[%s] Average time: %dms", test.name, milliseconds/int64(test.numberOfLaunches))
		})
	}
}

type BalanceLocker interface {
	InsertTestData(ctx context.Context, count int, amount int) ([]int, error)
	LockBalance(ctx context.Context, id int, amount int) error
	GetCountNegativeBalances(ctx context.Context) (int, error)
}

type TestCase struct {
	locker           BalanceLocker
	name             string
	count            int
	amount           int
	multiplier       int
	numberOfLaunches int
}

type SkipError string

func (e SkipError) Error() string {
	return string(e)
}

type RepositoryWithMutex struct {
	BaseRepository
	db *pgxpool.Pool
	mu sync.Mutex
}

func (r *RepositoryWithMutex) LockBalance(ctx context.Context, id int, amount int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	selectSQL := "select amount from mutex where id = $1"

	var current int

	row := r.db.QueryRow(ctx, selectSQL, id)
	err := row.Scan(&current)
	if err != nil {
		return err
	}

	if current < amount {
		return SkipError("balance not enough")
	}

	updateSQL := "update mutex set amount = amount - $2 where id = $1"
	command, err := r.db.Exec(ctx, updateSQL, id, amount)
	if err != nil {
		return err
	}

	if command.RowsAffected() == 0 {
		return errors.New("row was not affected")
	}

	return nil
}

type RepositoryWithOptimisticLock struct {
	BaseRepository
	db *pgxpool.Pool
}

func (r *RepositoryWithOptimisticLock) LockBalance(ctx context.Context, id int, amount int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	// Обновляем баланс.
	updateSQL := "update optimistic_lock set amount = amount - $2 where id = $1"
	command, err := tx.Exec(ctx, updateSQL, id, amount)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if command.RowsAffected() == 0 {
		_ = tx.Rollback(ctx)
		return errors.New("row was not affected")
	}

	// Проверяем, что не вышли за пределы баланса.
	var current int
	selectSQL := "select amount from optimistic_lock where id = $1"

	row := tx.QueryRow(ctx, selectSQL, id)
	err = row.Scan(&current)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// Если вышли, то откатываем транзакцию.
	if current < 0 {
		_ = tx.Rollback(ctx)
		return SkipError("balance not enough")
	}

	return tx.Commit(ctx)
}

type BaseRepository struct {
	table string
	db    *pgxpool.Pool
}

func (r *BaseRepository) InsertTestData(ctx context.Context, count int, amount int) ([]int, error) {
	sql := "insert into " + r.table + " (amount) " +
		"select $2 as amount " +
		"from generate_series(1, $1) " +
		"returning id;"

	rows, err := r.db.Query(ctx, sql, count, amount)
	if err != nil {
		return nil, err
	}

	ids := make([]int, count)
	for i := 0; rows.Next(); i++ {
		if err = rows.Scan(&ids[i]); err != nil {
			return nil, err
		}
	}

	return ids, nil
}

func (r *BaseRepository) GetCountNegativeBalances(ctx context.Context) (int, error) {
	count := 0
	sql := "select count(1) from " + r.table + " where amount < 0;"

	row := r.db.QueryRow(ctx, sql)
	err := row.Scan(&count)

	return count, err
}
