CREATE or REPLACE PROCEDURE sp_set_task_progress(
    p_time timestamptz 
  , p_event CHARACTER VARYING(16)
  , p_task_id uuid
  , p_user_id uuid)
LANGUAGE 'plpgsql' 
as $$
declare 
	cur record;
begin    
	if (p_event = 'MAIL_SEND') then
		insert into task_uuid (user_id, task_id, create_time, approve_time) values (p_user_id, p_task_id, p_time, null);
	end if;

	if (p_event = 'CLICK_APPROVE') then 
		update task_uuid set approve_time = p_time where task_id  = p_task_id and user_id = p_user_id;
		for cur in select tu.create_time , tu.approve_time from task_uuid tu where tu.task_id  = p_task_id and tu.user_id = p_user_id loop
	    	--Добавим время от создания до клика к общему времени жизни задачи
			update aggr_task_livetime 
			   set livetime = livetime + eXTRACT(EPOCH FROM age(cur.approve_time, cur.create_time))			     
			 where task_id = p_task_id;			
		end loop;		
	end if;

	insert into task_progress_events (task_id, user_id, event_time, state) values (p_task_id, p_user_id, p_time, p_event);
	
end;
$$;