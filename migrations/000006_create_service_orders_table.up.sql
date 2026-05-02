CREATE TABLE service_orders (
    id           UUID PRIMARY KEY,
    customer_id  UUID          NOT NULL REFERENCES customers(id),
    vehicle_id   UUID          NOT NULL REFERENCES vehicles(id),
    status       VARCHAR(50)   NOT NULL DEFAULT 'received',
    notes        TEXT,
    total_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    started_at   TIMESTAMP,
    finished_at  TIMESTAMP,
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
);
