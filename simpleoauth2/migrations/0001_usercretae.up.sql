create table users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    oauth_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    phone TEXT NOT NULL
);
