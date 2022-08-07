CREATE or REPLACE PROCEDURE sp_set_task_state(
    p_time timestamptz 
  , p_state CHARACTER VARYING(8)
  , p_task_id uuid)
LANGUAGE 'plpgsql'
as $$
declare 
	cur record;
begin
    --аггрегация по состояниям задач
	for cur in 
        select t.last_state, t.task_id 
          from aggr_task_livetime t 
         where t.task_id = p_task_id
    loop 
	  --уменьшим кол-во задач в старом состоянии		 
    	update aggr_state 
    	   set cnt = cnt - 1 
    	 where state = cur.last_state;	
    end loop; 
	--увеличим кол-во задач в новом состоянии		 
	update aggr_state set cnt = cnt + 1 where state = p_state;	    
   
    --аггрегация по времени жизни задачи
	if (p_state = 'CREATED') THEN
		-- если таска пересоздана - начинаем отслеживание прогресса заново
		delete from task_uuid t where t.task_id = p_task_id;	
        -- уже накопленное время жизни не сбрасываем, только первичная вставка
	    INSERT into aggr_task_livetime (task_id, create_time, livetime, last_state) values (p_task_id, p_time, 0, p_state) ON CONFLICT DO NOTHING;
	end if;

	if (p_state = 'DECLINED') then    		
	    -- если пришел отказ по таске - считаем время жизни относительно последнего согласующего, которому отправили письмо
		for cur in (select COALESCE(max(t_uuid.create_time), p_time) calc_time from task_uuid t_uuid where t_uuid.task_id = p_task_id ) loop
			update aggr_task_livetime 
			   set livetime = livetime + eXTRACT(EPOCH FROM age(p_time, cur.calc_time))			     
			where task_id = p_task_id;
		end loop;
	
	end if;
 
	update aggr_task_livetime set last_state = p_state where task_id = p_task_id;
    -- сохраняем информацию о событии в лог
    insert into task_state_events(event_time, state, task_id) values (p_time, p_state, p_task_id);
 
    return;
end;  
$$;
