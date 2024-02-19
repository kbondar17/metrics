package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	db "metrics/internal/database"
	appErrors "metrics/internal/errors"
	"metrics/internal/models"
	"net"
	"reflect"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

// type PostgresStorage struct {
// 	Conn       *pgx.Conn
// 	CanRetrier appErrors.RetryableError
// 	// Dns string
// }

type PostgresStorage struct {
	Conn       *sql.DB
	CanRetrier appErrors.RetryableError
	// Dns string
}

// func initDB(conn *pgx.Conn) error {
// 	stmt := `CREATE TABLE metric (
// 		id TEXT PRIMARY KEY,
// 		mtype TEXT NOT NULL,
// 		delta BIGINT,
// 		value DOUBLE PRECISION,
// 		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
// 		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
// 	);
// 	CREATE INDEX idx_update_metrics_model_id ON metric (id);
// 	`

// 	_, err := conn.Exec(stmt)

// 	if err != nil {
// 		var pgErr pgx.PgError

// 		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.DuplicateTable {
// 			log.Println("Table metric already exists")
// 			return nil
// 		}
// 		return err
// 	}
// 	log.Println("Table metric created")
// 	return nil
// }

func initDB(conn *sql.DB) error {
	stmt := `CREATE TABLE metric (
		id TEXT PRIMARY KEY,
		mtype TEXT NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX idx_update_metrics_model_id ON metric (id);
	`

	_, err := conn.Exec(stmt)

	if err != nil {
		fmt.Println("err:::", err, "type:::", reflect.TypeOf(err))
		var pgErr *pgx.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.DuplicateTable {
			log.Println("Table metric already exists")
			return nil
		}
		var pgErr2 *pgconn.PgError

		if errors.Is(err, pgErr2) {
			log.Println("Table metric already exists")
			return nil
		}

		return err
	}
	log.Println("Table metric created")
	return nil
}

func NewPostgresStorage(dns string, logger *zap.SugaredLogger) (PostgresStorage, error) {

	conn, err := sql.Open("pgx", dns)

	if err != nil {
		logger.Errorf("unable to open postgres conn: %v", err)
		return PostgresStorage{}, err
	}

	// err = initDB(conn)

	if err != nil {
		return PostgresStorage{}, err
	}
	retrErr := appErrors.NewRetryableError()
	return PostgresStorage{Conn: conn, CanRetrier: *retrErr}, nil
}

// func NewPostgresStorage(dns string, logger *zap.SugaredLogger) (PostgresStorage, error) {

// 	config, err := pgx.ParseDSN(dns)
// 	if err != nil {
// 		logger.Errorf("unable to parse dns: %v", err)
// 		return PostgresStorage{}, err
// 	}

// 		logger.Infoln("connecting to db with config: ", config)

// 	conn, err := pgx.Connect(config)
// 	if err != nil {
// 		logger.Errorf("unable to connect to db: %v, dns %v", err, dns)
// 		return PostgresStorage{}, err
// 	}
// 	err = initDB(conn)

// 	if err != nil {
// 		return PostgresStorage{}, err
// 	}
// 	retrErr := appErrors.NewRetryableError()
// 	return PostgresStorage{Conn: conn, CanRetrier: *retrErr}, nil

// }

func errIsRetriable(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return true
	}

	if errors.Is(err, pgx.ErrDeadConn) {
		return true
	}
	var e *net.OpError
	return errors.As(err, &e)

}

func (p PostgresStorage) Ping() error {
	inner := func() error {
		stmt := "select 1"

		row := p.Conn.QueryRow(stmt)
		var result int
		err := row.Scan(&result)

		if err != nil {
			fmt.Println("err type:::", reflect.TypeOf(err))
			return err
		} else {
			fmt.Println("Pong")
		}
		return nil
	}
	// TODO: move it to class
	retrErr := appErrors.NewRetryableError()
	return appErrors.RetryWrapper(inner, errIsRetriable, *retrErr)
}

// func (p PostgresStorage) transactionWrapper(f func(tx *pgx.Tx) error) error {
func (p PostgresStorage) transactionWrapper(f func(tx *sql.Tx) error) error {
	tx, err := p.Conn.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if pnc := recover(); pnc != nil {
			tx.Rollback()
			panic(pnc) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit() // return this error if commit fails
		}
	}()

	err = f(tx)
	return err

}

func (p PostgresStorage) UpdateMultipleMetric(metrics []models.UpdateMetricsModel) error {
	query := func(tx *sql.Tx) error {
		sql := "INSERT INTO metric (id, mtype, delta, value) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET delta = $3, value = $4, updated_at = now() ;"
		for _, metric := range metrics {
			_, err := tx.Exec(sql, metric.ID, metric.MType, metric.Delta, metric.Value)
			if err != nil {
				fmt.Println("err while updating:::", err)
				return err
			} else {
				fmt.Println("udpated ok :: ", metric.ID, metric.MType, metric.Delta, metric.Value)
			}
		}
		return nil
	}

	transaction := func() error {
		return p.transactionWrapper(query)
	}

	return appErrors.RetryWrapper(transaction, errIsRetriable, p.CanRetrier)

}

func (p PostgresStorage) CheckIfMetricExists(name string, mType models.MetricType) (bool, error) {
	return true, nil
	// stmt := func() (interface{}, error) {
	// 	sql := "SELECT EXISTS(SELECT 1 FROM metric WHERE id = $1 and mtype = $2);"
	// 	row := p.Conn.QueryRow(sql, name, mType)
	// 	var result bool
	// 	err := row.Scan(&result)
	// 	if err != nil {
	// 		return false, err
	// 	}
	// 	return result, nil
	// }

	// rawRes, err := appErrors.RetryWrapperWithResult(stmt, errIsRetriable, p.CanRetrier)

	// if err != nil {
	// 	return false, fmt.Errorf("unable to check if metric exists: %w", err)
	// }
	// return rawRes.(bool), nil

}

func (p PostgresStorage) GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error) {
	stmt := func() (interface{}, error) {
		sql := "SELECT value from metric where id = $1 and mtype = $2;"
		row := p.Conn.QueryRow(sql, name, mType)
		var result float64
		err := row.Scan(&result)
		if err != nil {
			return 0, err
		}
		return result, nil
	}

	rawRes, err := appErrors.RetryWrapperWithResult(stmt, errIsRetriable, p.CanRetrier)

	if err != nil {
		return 0, fmt.Errorf("unable to get metric value: %w", err)
	}
	fmt.Println("rawRes:::", rawRes)
	return rawRes.(float64), nil
}

func (p PostgresStorage) GetCountMetricValueByName(name string) (int64, error) {
	stmt := func() (interface{}, error) {
		sql := "SELECT delta from metric where id = $1 and mtype = $2;"
		row := p.Conn.QueryRow(sql, name, models.CounterType)
		var result int64
		err := row.Scan(&result)
		if err != nil {
			return 0, err
		}
		return result, nil
	}

	rawRes, err := appErrors.RetryWrapperWithResult(stmt, errIsRetriable, p.CanRetrier)

	if err != nil {
		return 0, fmt.Errorf("unable to get metric value: %w", err)
	}
	return rawRes.(int64), nil

}

// TODO: должен ли падать с ошишбкой?
func (p PostgresStorage) Create(metricName string, metricType models.MetricType) error {
	stmt := func(tx *sql.Tx) error {

		// тут проблема потому создает с пустым значением
		sql := "INSERT INTO metric (id, mtype) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING;"
		_, err := tx.Exec(sql, metricName, metricType)
		return err
	}
	transaction := func() error {
		return p.transactionWrapper(stmt)
	}

	return appErrors.RetryWrapper(transaction, errIsRetriable, p.CanRetrier)
}

func (p PostgresStorage) UpdateMetric(name string, metricType models.MetricType, value interface{}, syncStorage bool, storagePath string) error {
	var metric models.UpdateMetricsModel

	if metricType == models.CounterType {
		val, ok := value.(int64)
		if !ok {
			return fmt.Errorf("value %s is not int64", value)
		}
		metric = models.UpdateMetricsModel{ID: name, MType: string(models.CounterType), Delta: &val}
	} else {
		val, ok := value.(float64)
		if !ok {
			return fmt.Errorf("value %s is not int64", value)
		}
		metric = models.UpdateMetricsModel{ID: name, MType: string(models.GaugeType), Value: &val}
	}

	query := func(tx *sql.Tx) error {
		sql := `INSERT INTO metric (id, mtype, delta, value)
	        VALUES ($1, $2, $3, $4)
	        ON CONFLICT (id) DO UPDATE SET delta = $3, value = $4, updated_at = now()`

		// _, err := tx.Exec(sql, metric.ID, metric.MType, metric.Delta, metric.Value)
		_, err := tx.Exec(sql, metric.ID, metric.MType, 11, 22.22)

		return err

	}
	// transaction := func() error {
	// 	return p.transactionWrapper(query)
	// }
	// err := appErrors.RetryWrapper(transaction, errIsRetriable, p.CanRetrier)
	tx, err := p.Conn.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	err = query(tx)
	// err := p.transactionWrapper(query)

	if err != nil {
		return fmt.Errorf("unable to update metric: %w", err)
	}
	if syncStorage {
		log.Println("saving metric to file: ", metric)
		err := db.SaveMetric(storagePath, metric)
		if err != nil {
			return fmt.Errorf("unable to save metric to file: %w", err)
		}
	}
	return nil

}

func (p PostgresStorage) GetAllMetrics() ([]models.UpdateMetricsModel, error) {
	query := "SELECT id, mtype, value, delta FROM metric;"

	// if err := p.Ping(); err != nil {
	// 	return nil, fmt.Errorf("unable to ping db: %w", err)
	// }
	rows, err := p.Conn.Query(query)

	if err != nil {
		return nil, fmt.Errorf("unable to execute query: %w", err)
	}
	defer rows.Close()

	var result []models.UpdateMetricsModel

	for rows.Next() {
		var metric models.UpdateMetricsModel
		err = rows.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return nil, fmt.Errorf("unable to scan row: %w", err)
		}
		result = append(result, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return result, nil
}
