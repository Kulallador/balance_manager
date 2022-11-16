package dbmanager

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	_ "github.com/lib/pq"
)

var (
	userNotFound = fmt.Errorf("user not found")
)

type PostgresDB struct {
	conn *sql.DB
}

func CreateDB(conn *sql.DB) (PostgresDB, error) {
	db := PostgresDB{
		conn: conn,
	}

	err := db.CreateTables()
	if err != nil {
		return db, err
	}

	return db, nil
}

func (db *PostgresDB) GetAllBalance() ([]BalanceTransaction, error) {
	resp, err := db.conn.Query(`Select * FROM user_balance`)
	if err != nil {
		return nil, err
	}

	balances := make([]BalanceTransaction, 1)
	for resp.Next() {
		var balanceTransaction BalanceTransaction
		err := resp.Scan(&balanceTransaction.UserID, &balanceTransaction.Money)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balanceTransaction)

	}

	return balances, nil
}

func (db *PostgresDB) CreateTables() error {
	if _, err := db.conn.Exec(`
		CREATE TABLE IF NOT EXISTS user_balance (
			user_id INT PRIMARY KEY,
			balance NUMERIC CONSTRAINT positive_balance CHECK (balance >= 0) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS reserve (
			user_id INT REFERENCES user_balance (user_id),
			service_id INT,
			order_id INT,
			PRIMARY KEY (user_id, service_id, order_id),
			balance NUMERIC CONSTRAINT positive_balance CHECK (balance >= 0) NOT NULL
		);
		CREATE TABLE IF NOT EXISTS report (
			user_id INT REFERENCES user_balance (user_id) NOT NULL,
			service_id INT NOT NULL,
			order_id INT NOT NULL,
			balance NUMERIC CONSTRAINT positive_balance CHECK (balance >= 0) NOT NULL,
			date TIMESTAMP NOT NULL
		);`,
	); err != nil {
		return err
	}

	return nil
}

func (db *PostgresDB) IncrementMoney(ctx context.Context, transaction BalanceTransaction) error {
	tr, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begintx: %w", err)
	}
	defer tr.Rollback()

	res, err := tr.ExecContext(ctx, `
		UPDATE user_balance
		SET balance=balance+$1
		WHERE user_id=$2`,
		transaction.Money,
		transaction.UserID,
	)
	if err != nil {
		return fmt.Errorf("update user balance: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows updated: %w", err)
	}

	if n == 0 {
		if _, err := tr.ExecContext(ctx, `
			INSERT INTO user_balance
			VALUES ($1, $2)`,
			transaction.UserID,
			transaction.Money,
		); err != nil {
			return fmt.Errorf("create new user: %w", err)
		}
	}

	err = tr.Commit()
	if err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (db *PostgresDB) DecrementMoney(ctx context.Context, transaction BalanceTransaction) error {
	res, err := db.conn.ExecContext(ctx, `
		UPDATE user_balance
		SET balance=balance-$1
		WHERE user_id=$2`,
		transaction.Money,
		transaction.UserID,
	)
	if err != nil {
		return fmt.Errorf("decrement user balance: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows updated: %w", err)
	}
	if n == 0 {
		return userNotFound
	}

	return nil
}

func (db *PostgresDB) TranslationMoney(ctx context.Context, transaction TranslateTransaction) error {
	tr, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tr.Rollback()

	res, err := tr.ExecContext(ctx, `
		UPDATE user_balance 
		SET balance=balance-$1 
		WHERE user_id=$2`,
		transaction.Money,
		transaction.FromID,
	)
	if err != nil {
		return fmt.Errorf("decrement user balance: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get count from ids: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("get from_id: %w", userNotFound)
	}

	res, err = tr.ExecContext(ctx, `
		UPDATE user_balance
		SET balance=balance+$1
		WHERE user_id=$2`,
		transaction.Money,
		transaction.ToID,
	)
	if err != nil {
		return fmt.Errorf("increment user balance: %w", err)
	}

	n, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get count to ids: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("get to_id: %w", userNotFound)
	}

	if err := tr.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *PostgresDB) GetUserBalance(ctx context.Context, user_id int) (float64, error) {
	var balance float64
	if err := db.conn.QueryRowContext(ctx, `
		SELECT balance 
		FROM user_balance
		WHERE user_id=$1`,
		user_id,
	).Scan(&balance); err != nil {
		return balance, fmt.Errorf("get user balance: %w", err)
	}

	return balance, nil
}

func (db *PostgresDB) GetReserveBalance(ctx context.Context, transaction ReserveTransaction) (float64, error) {
	var balance float64
	err := db.conn.QueryRowContext(ctx, `
		SELECT balance 
		FROM reserve
		WHERE user_id=$1 AND service_id=$2 AND order_id=$3`,
		transaction.UserID,
		transaction.ServiceID,
		transaction.OrderID,
	).Scan(&balance)
	if err != nil {
		return balance, fmt.Errorf("get balance: %w", err)
	}

	return balance, nil
}

func (db *PostgresDB) ReserveMoney(ctx context.Context, transaction ReserveTransaction) error {
	tr, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tr.Rollback()

	res, err := tr.ExecContext(ctx, `
		UPDATE user_balance
		SET balance=balance-$1
		WHERE user_id=$2`,
		transaction.Money,
		transaction.UserID,
	)
	if err != nil {
		return fmt.Errorf("decrement user balance: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows updated: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("get user_id: %w", userNotFound)
	}

	res, err = tr.ExecContext(ctx, ` 
		UPDATE reserve
		SET balance=balance+$1
		WHERE user_id=$2 AND service_id=$3 AND order_id=$4`,
		transaction.Money,
		transaction.UserID,
		transaction.ServiceID,
		transaction.OrderID,
	)
	if err != nil {
		return fmt.Errorf("increment reserve balance: %w", err)
	}

	n, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows updated: %w", err)
	}
	if n == 0 {
		if _, err = tr.Exec(`
			INSERT INTO reserve
			VALUES ($1, $2, $3, $4)`,
			transaction.UserID,
			transaction.ServiceID,
			transaction.OrderID,
			transaction.Money,
		); err != nil {
			return fmt.Errorf("create reserve: %w", err)
		}
	}

	if err := tr.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

func (db *PostgresDB) DecReservedMoney(ctx context.Context, transaction ReserveTransaction) error {
	tr, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tr.Rollback()

	if _, err := tr.ExecContext(ctx, `
		UPDATE reserve
		SET balance=balance-$1
		WHERE user_id=$2 AND service_id=$3 AND order_id=$4`,
		transaction.Money,
		transaction.UserID,
		transaction.ServiceID,
		transaction.OrderID,
	); err != nil {
		return fmt.Errorf("decrement reserve balance: %w", err)
	}

	if _, err := tr.ExecContext(ctx, `
		INSERT INTO report
		VALUES ($1, $2, $3, $4, $5)`,
		transaction.UserID,
		transaction.ServiceID,
		transaction.OrderID,
		transaction.Money,
		time.Now(),
	); err != nil {
		return fmt.Errorf("create report row: %w", err)
	}

	if err := tr.Commit(); err != nil {
		return fmt.Errorf("context: %w", err)
	}

	return nil
}
