-- +goose Up
create table refresh_tokens(
    token text primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id UUID not null,
    expires_at timestamp not null,
    revoked_at timestamp,
    constraint fk_uid_users
        foreign key(user_id)
        references users(id) on delete cascade   
);

-- +goose Down
drop table refresh_tokens;
