package sqlutil_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/puppetlabs/horsehead/v2/sqlutil"
	"github.com/stretchr/testify/require"
)

func TestTxSimple(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"<value>"}).AddRow(1))
	mock.ExpectCommit()

	require.NoError(t, sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.QueryContext(ctx, "SELECT 1")
		return err
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxNested(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE t SET a = 1").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("SAVEPOINT tx_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE t SET a = 2").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("RELEASE SAVEPOINT tx_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SAVEPOINT tx_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE t SET a = 3").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("SAVEPOINT tx_2").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE t SET a = 4").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("RELEASE SAVEPOINT tx_2").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("RELEASE SAVEPOINT tx_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	require.NoError(t, sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, "UPDATE t SET a = 1"); err != nil {
			return err
		}

		err := sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
			_, err = tx.ExecContext(ctx, "UPDATE t SET a = 2")
			return err
		})
		if err != nil {
			return err
		}

		return sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
			if _, err := tx.ExecContext(ctx, "UPDATE t SET a = 3"); err != nil {
				return err
			}

			return sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
				_, err = tx.ExecContext(ctx, "UPDATE t SET a = 4")
				return err
			})
		})
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxSimpleRollback(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT 1").WillReturnError(fmt.Errorf("in test"))
	mock.ExpectRollback()

	require.NotNil(t, sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.QueryContext(ctx, "SELECT 1")
		return err
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTxNestedRollback(t *testing.T) {
	ctx := context.Background()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE t SET a = 1").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("SAVEPOINT tx_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE t SET a = 2").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("RELEASE SAVEPOINT tx_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SAVEPOINT tx_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("SAVEPOINT tx_2").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE t SET a = 3").WillReturnError(fmt.Errorf("in test"))
	mock.ExpectExec("ROLLBACK TO SAVEPOINT tx_2").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("UPDATE t SET a = 4").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("RELEASE SAVEPOINT tx_1").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	require.NoError(t, sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, "UPDATE t SET a = 1"); err != nil {
			return err
		}

		err := sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
			_, err = tx.ExecContext(ctx, "UPDATE t SET a = 2")
			return err
		})
		if err != nil {
			return err
		}

		return sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
			err := sqlutil.WithTx(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
				_, err = tx.ExecContext(ctx, "UPDATE t SET a = 3")
				return err
			})
			require.NotNil(t, err)

			_, err = tx.ExecContext(ctx, "UPDATE t SET a = 4")
			return err
		})
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}
