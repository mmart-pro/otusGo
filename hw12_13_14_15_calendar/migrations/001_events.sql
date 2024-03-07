-- +goose Up
-- +goose StatementBegin
create table public.events (
    id bigserial constraint pk_events_id primary key,
    title character varying(255) not null,
    start_date_time timestamp with time zone not null,
	end_date_time timestamp with time zone not null,
    description text not null,
    user_id int not null,
    notify_before_min int not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists public.events;
-- +goose StatementEnd
