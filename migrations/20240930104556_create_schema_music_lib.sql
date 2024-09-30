-- +goose Up
-- +goose StatementBegin
create schema if not exists music_lib;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop schema if exists music_lib;
-- +goose StatementEnd
