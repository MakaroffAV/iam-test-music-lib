-- +goose Up
-- +goose StatementBegin
create table if not exists music_lib.song (
    id uuid primary key,
    youtube text not null,
    released date not null,
    album varchar(100) not null,
    title varchar(100) not null,
    artist varchar(100) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists music_lib.song;
-- +goose StatementEnd
