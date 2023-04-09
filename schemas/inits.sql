CREATE TABLE users
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    phone TEXT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);

CREATE TABLE tokens
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token TEXT NOT NULL
);

CREATE TABLE accounts
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number TEXT NOT NULL,
    user_id UUID NOT NULL
        REFERENCES users ON DELETE CASCADE,
    balance DECIMAL NOT NULL DEFAULT 0.0,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);

CREATE TABLE transactions
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL
        REFERENCES accounts ON DELETE CASCADE,
    type TEXT NOT NULL,
    amount DECIMAL NOT NULL DEFAULT 0.0,
    created_at timestamptz NOT NULL DEFAULT current_timestamp,
    updated_at timestamptz,
    deleted_at timestamptz
);