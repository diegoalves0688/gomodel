CREATE TABLE message (
    id         uuid DEFAULT public.uuid_generate_v4() NOT NULL PRIMARY KEY,
    receiver   VARCHAR(35) NOT NULL,
    sender     VARCHAR(35) NOT NULL,
    content    VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);