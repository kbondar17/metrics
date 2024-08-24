-- +goose Up
CREATE TABLE metric (
		id TEXT PRIMARY KEY,
		mtype TEXT NOT NULL,
		delta BIGINT,
		value DOUBLE PRECISION,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX idx_update_metrics_model_id ON metric (id);

-- +goose Down
DROP TABLE metric;
DROP INDEX idx_update_metrics_model_id;
