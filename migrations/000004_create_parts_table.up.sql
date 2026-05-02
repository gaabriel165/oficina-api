CREATE TABLE parts (
    id             UUID PRIMARY KEY,
    name           VARCHAR(255)    NOT NULL,
    description    TEXT,
    unit_price     DECIMAL(10,2)   NOT NULL,
    stock_quantity INT             NOT NULL DEFAULT 0,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW()
);
