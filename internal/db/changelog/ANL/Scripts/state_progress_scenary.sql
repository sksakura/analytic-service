CALL anl.sp_set_task_state(now(),'CREATED','123e4567-e89b-12d3-a456-426614174001');
CALL anl.sp_set_task_progress(now(),'MAIL_SEND','123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174002');
CALL anl.sp_set_task_progress(now(),'CLICK_APPROVE','123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174002');
CALL anl.sp_set_task_progress(now(),'MAIL_SEND','123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174003');
CALL anl.sp_set_task_state(now(),'DECLINED','123e4567-e89b-12d3-a456-426614174001');

CALL anl.sp_set_task_state(now(),'CREATED','123e4567-e89b-12d3-a456-426614174001');
CALL anl.sp_set_task_progress(now(),'MAIL_SEND','123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174002');
CALL anl.sp_set_task_progress(now(),'CLICK_APPROVE','123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174002');
CALL anl.sp_set_task_progress(now(),'MAIL_SEND','123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174003');
CALL anl.sp_set_task_progress(now(),'CLICK_APPROVE','123e4567-e89b-12d3-a456-426614174001', '123e4567-e89b-12d3-a456-426614174003');
CALL anl.sp_set_task_state(now(),'APPROVED','123e4567-e89b-12d3-a456-426614174001');

select * from anl.aggr_task_livetime atl where atl.task_id = '123e4567-e89b-12d3-a456-426614174001';
select * from anl.task_progress_events tpe where tpe.task_id  = '123e4567-e89b-12d3-a456-426614174001';
select * from anl.task_uuid tu where tu.task_id = '123e4567-e89b-12d3-a456-426614174001';
select * from anl.aggr_state as2; 

