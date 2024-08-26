package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	db "metrics/internal/database"
	appErrors "metrics/internal/errors"
	"metrics/internal/models"
	"net"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	Conn       *sql.DB
	CanRetrier appErrors.RetryableError
	logger     *zap.SugaredLogger
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func initDB(conn *sql.DB, logger *zap.SugaredLogger) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("unable to set postgres as dialect in goose: %w", err)
	}

	if err := goose.Up(conn, "migrations"); err != nil {
		return fmt.Errorf("unable to run migrations: %w", err)
	}
	logger.Infof("Table metric created")
	return nil
}

func NewPostgresStorage(dns string, logger *zap.SugaredLogger) (PostgresStorage, error) {
	conn, err := sql.Open("pgx", dns)

	if err != nil {
		logger.Errorf("unable to open postgres conn: %v", err)
		return PostgresStorage{}, err
	}

	err = initDB(conn, logger)

	if err != nil {
		return PostgresStorage{}, err
	}
	retrErr := appErrors.NewRetryableError()
	return PostgresStorage{Conn: conn, CanRetrier: *retrErr, logger: logger}, nil
}

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
	if p.Conn == nil {
		return errors.New("database connection is not initialized")
	}

	inner := func() error {
		stmt := "select 1"

		row := p.Conn.QueryRow(stmt)
		var result int
		err := row.Scan(&result)

		if err != nil {
			return fmt.Errorf("unable to ping db: %w", err)
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
		sql := `INSERT INTO metric (id, mtype, delta, value) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (id) DO UPDATE 
		SET 
			delta = CASE 
				WHEN metric.mtype = 'gauge' THEN EXCLUDED.delta
				WHEN metric.mtype = 'counter' THEN metric.delta + EXCLUDED.delta
			END,
			value = $4, 
			updated_at = now()`

		for _, metric := range metrics {
			_, err := tx.Exec(sql, metric.ID, metric.MType, metric.Delta, metric.Value)
			if err != nil {
				fmt.Println("err while updating:::", err)
				return err
			} else {
				if metric.MType == string(models.CounterType) {
					// fmt.Println("udpated ok :: ", metric.ID, metric.MType, *metric.Delta)
				} else {
					// fmt.Println("udpated ok :: ", metric.ID, metric.MType, *metric.Value)
				}
			}
		}
		return nil
	}

	transaction := func() error {
		return p.transactionWrapper(query)
	}

	return appErrors.RetryWrapper(transaction, errIsRetriable, p.CanRetrier)

}

// CheckIfMetricExists is a stub for backward compatibility
func (p PostgresStorage) CheckIfMetricExists(name string, mType models.MetricType) (bool, error) {
	return true, nil
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

	if res, ok := rawRes.(float64); ok {
		return res, nil
	}
	return 0, fmt.Errorf("unable to convert result to float64 %s", rawRes)

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

	if res, ok := rawRes.(int64); ok {
		return res, nil
	}
	return 0, fmt.Errorf("unable to convert result to int64 %s", rawRes)

}

// TODO: должен ли падать с ошишбкой?
func (p PostgresStorage) Create(metricName string, metricType models.MetricType) error {
	stmt := func(tx *sql.Tx) error {

		// FIXME: тут проблема потому создает с пустым значением
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
		ON CONFLICT (id) DO UPDATE 
		SET 
			delta = CASE 
				WHEN metric.mtype = 'gauge' THEN EXCLUDED.delta
				WHEN metric.mtype = 'counter' THEN metric.delta + EXCLUDED.delta
			END,
			value = $4, 
			updated_at = now()`

		_, err := tx.Exec(sql, metric.ID, metric.MType, metric.Delta, metric.Value)
		return err

	}
	transaction := func() error {
		return p.transactionWrapper(query)
	}

	err := appErrors.RetryWrapper(transaction, errIsRetriable, p.CanRetrier)

	if err != nil {
		return fmt.Errorf("unable to update metric: %w", err)
	}
	if syncStorage {
		p.logger.Infoln("saving metric to file: ", metric)
		err := db.SaveMetric(storagePath, metric)
		if err != nil {
			return fmt.Errorf("unable to save metric to file: %w", err)
		}
	}
	return nil

}

func (p PostgresStorage) GetAllMetrics() ([]models.UpdateMetricsModel, error) {
	if p.Conn == nil {
		return nil, errors.New("database connection is not initialized")
	}

	query := "SELECT id, mtype, value, delta FROM metric;"

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
