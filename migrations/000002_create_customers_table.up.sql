CREATE TABLE customers (
    id            UUID PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    document      VARCHAR(14)  NOT NULL UNIQUE,
    document_type VARCHAR(4)   NOT NULL,
    phone         VARCHAR(20)  NOT NULL,
    email         VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
);
