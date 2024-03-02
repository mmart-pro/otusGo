-- +goose Up
-- +goose StatementBegin
alter table public.events
add column is_notified bool not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table public.events
drop column if exists is_notified;
-- +goose StatementEnd
