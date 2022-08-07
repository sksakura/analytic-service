CREATE TABLE task_uuid(  
    user_id uuid NOT NULL PRIMARY KEY,
    task_id uuid NOT NULL,
    create_time timestamptz NOT NULL,
    approve_time timestamptz,
    CONSTRAINT task_uuid_UK1 UNIQUE (task_id, create_time)
);