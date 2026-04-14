CREATE TABLE IF NOT EXISTS tasks (
	id BIGSERIAL PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL,
	
	-- можно потом вынести в отдельную таблицу scheduleRule <- task
	-- срок выполнения
	scheduled_at TIMESTAMP WITH TIME ZONE NOT NULL, 
	-- тип частоты повторения. мб отдельная таблица
	type_of_repetition VARCHAR(30) NOT NULL, 
	-- переодичность в днях
	periodicity INT NOT NULL CHECK(periodicity > 0 AND periodicity < 31), 
	dates JSONB DEFAULT '[]' NOT NULL,

	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks (status);
