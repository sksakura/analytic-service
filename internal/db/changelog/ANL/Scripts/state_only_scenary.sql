CALL sp_set_task_state(now(),'CREATED','123e4567-e89b-12d3-a456-426614174000');
CALL sp_set_task_state(now(),'APPROVED','123e4567-e89b-12d3-a456-426614174000');
CALL sp_set_task_state(now(),'DECLINED','123e4567-e89b-12d3-a456-426614174000');


select * from task_uuid tu ;
select * from task_state_events tse ;
select * from aggr_task_livetime atl;
select * from aggr_state as2; 

delete from task_progress_events;
delete from aggr_task_livetime t ;
update aggr_state set cnt = 0;

