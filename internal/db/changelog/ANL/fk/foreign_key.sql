alter table task_state_events add constraint fk1_state_events foreign key (State) REFERENCES aggr_state (state);
alter table task_state_events add constraint fk2_state_events foreign key (task_id) REFERENCES aggr_task_livetime (task_id) ON DELETE CASCADE;
alter table task_uuid add constraint fk1_task_uuid foreign key (task_id) references aggr_task_livetime(task_id)  ON DELETE CASCADE;
alter table aggr_task_livetime  add constraint fk1_task_livetime foreign key (last_state) REFERENCES aggr_state (state);
ALTER TABLE task_progress_events ADD CONSTRAINT fk_task_progress_events_task_id FOREIGN KEY (task_id) REFERENCES aggr_task_livetime (task_id);
