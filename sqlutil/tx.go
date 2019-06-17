package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"time"

	"github.com/puppetlabs/horsehead/lifecycle"
)

type txContextKey uintptr

type txContextValue struct {
	tx *sql.Tx
	c  uint64
}

func (key txContextKey) Get(ctx context.Context) (v txContextValue, ok bool) {
	v, ok = ctx.Value(key).(txContextValue)
	return
}

func (key txContextKey) Set(ctx context.Context, v txContextValue) context.Context {
	return context.WithValue(ctx, key, v)
}

type txDelegate interface {
	Rollback() error
	Commit() error
}

type savepointDelegate struct {
	ctx context.Context
	v   txContextValue
}

func (sd *savepointDelegate) Rollback() error {
	return lifecycle.NewCloserBuilder().
		RequireContext(func(ctx context.Context) error {
			_, err := sd.v.tx.ExecContext(sd.ctx, fmt.Sprintf("ROLLBACK TO SAVEPOINT tx_%d", sd.v.c))
			return err
		}).
		Timeout(500 * time.Millisecond).
		Build().
		Do(sd.ctx)
}

func (sd *savepointDelegate) Commit() error {
	_, err := sd.v.tx.ExecContext(sd.ctx, fmt.Sprintf("RELEASE SAVEPOINT tx_%d", sd.v.c))
	if err != nil {
		sd.Rollback()
		return err
	}

	return nil
}

// WithTx executes the given function within a database transaction.
//
// This function is reentrant. It will create nested transactions using
// SAVEPOINTs if needed. However, transactions are not thread-safe, so care must
// be used when creating goroutines inside transactions: ensure a separate
// context is used that does not carry the transaction with it.
func WithTx(ctx context.Context, db *sql.DB, fn func(ctx context.Context, tx *sql.Tx) error) (err error) {
	key := txContextKey(reflect.ValueOf(db).Pointer())

	var d txDelegate

	v, ok := key.Get(ctx)
	if ok {
		// Use savepoints.
		v.c++

		if _, err := v.tx.ExecContext(ctx, fmt.Sprintf("SAVEPOINT tx_%d", v.c)); err != nil {
			return err
		}

		ctx = key.Set(ctx, v)
		d = &savepointDelegate{ctx: ctx, v: v}
	} else {
		// Use BeginTx().
		tx, err := db.BeginTx(ctx, &sql.TxOptions{})
		if err != nil {
			return err
		}

		v = txContextValue{tx: tx}
		ctx = key.Set(ctx, v)
		d = tx
	}

	defer func() {
		if p := recover(); p != nil {
			d.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx, v.tx); err != nil {
		d.Rollback()
		return err
	}

	return d.Commit()
}
