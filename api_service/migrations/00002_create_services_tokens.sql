-- +goose Up
-- +goose StatementBegin
CREATE TABLE services_tokens (
    id bigserial primary key,
    service_id int references services(id) not null,
    expires_at timestamp with time zone,
    value varchar unique not null CHECK ( length(value) > 20 )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE services_tokens;
-- +goose StatementEnd
