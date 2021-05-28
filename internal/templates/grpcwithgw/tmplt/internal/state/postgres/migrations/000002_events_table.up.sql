CREATE TABLE events 
(
    id bigserial PRIMARY KEY,
    resource_type integer NOT NULL,
    operation integer NOT NULL,
    payload bytea NOT NULL,
    create_time timestamptz NOT NULL default now(),
    creator text NOT NULL
);

CREATE INDEX events_paging_idx ON events(create_time DESC, id);