CREATE TABLE aggr_task_livetime(  
    task_id uuid NOT NULL PRIMARY KEY,
    create_time timestamptz  NOT NULL,
    livetime int NOT null,
    last_state varchar(8) not null
);