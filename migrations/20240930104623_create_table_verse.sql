-- +goose Up
-- +goose StatementBegin
create table if not exists music_lib.verse (
    id uuid primary key,
    text text not null,
    song_id uuid not null references music_lib.song(id) on delete cascade,
    order_num int not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists music_lib.verse;
-- +goose StatementEnd
