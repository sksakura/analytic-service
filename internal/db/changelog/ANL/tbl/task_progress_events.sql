CREATE TABLE task_progress_events (
	id serial4 NOT NULL,
	task_id uuid NOT NULL,
	user_id uuid NOT NULL,
	event_time timestamptz NULL,
    state VARCHAR(16) not null check (State in ('MAIL_SEND', 'CLICK_APPROVE')),
	CONSTRAINT task_progress_events_pkey PRIMARY KEY (id),
	CONSTRAINT task_progress_events_uk1 UNIQUE (event_time, task_id, user_id, state)
);
