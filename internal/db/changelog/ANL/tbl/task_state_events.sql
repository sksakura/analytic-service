CREATE TABLE task_state_events(  
    id SERIAL NOT NULL PRIMARY KEY,
    event_time timestamptz not null,
    task_id uuid not null,
    state VARCHAR(8) not null check (State in ('CREATED', 'DECLINED', 'APPROVED')),
    CONSTRAINT task_state_events_UK1 UNIQUE (event_time, task_id, state)
);