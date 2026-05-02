CREATE TABLE services (
    id                UUID PRIMARY KEY,
    name              VARCHAR(255)  NOT NULL,
    description       TEXT,
    labor_price       DECIMAL(10,2) NOT NULL,
    estimated_minutes INT           NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP NOT NULL DEFAULT NOW()
);
