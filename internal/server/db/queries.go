package db

const (
	InsertOrUpdateCounter = `
		INSERT INTO metrics (id, m_type, delta)
			 VALUES ($1, 'counter', $2)
			 ON CONFLICT (id) DO UPDATE SET delta = metrics.delta + EXCLUDED.delta;
	`

	InsertOrUpdateGauge = `
		INSERT INTO metrics (id, m_type, value)
			 VALUES ($1, 'gauge', $2)
			 ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value;
	`
)
