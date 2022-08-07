CREATE TABLE aggr_state(   
    State VARCHAR(8) not null PRIMARY KEY check (State in ('CREATED', 'DECLINED', 'APPROVED')),
    CNT int not null not null check (cnt >= 0) 
);
INSERT into aggr_state (CNT, State) VALUES (0, 'CREATED');
INSERT into aggr_state (CNT, State) VALUES (0, 'APPROVED');
INSERT into aggr_state (CNT, State) VALUES (0, 'DECLINED');

CREATE UNIQUE INDEX aggr_state_idx ON aggr_state (State);