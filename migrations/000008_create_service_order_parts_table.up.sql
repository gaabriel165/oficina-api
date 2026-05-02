CREATE TABLE service_order_parts (
    id               UUID PRIMARY KEY,
    service_order_id UUID          NOT NULL REFERENCES service_orders(id) ON DELETE CASCADE,
    part_id          UUID          NOT NULL REFERENCES parts(id),
    part_name        VARCHAR(255)  NOT NULL,
    quantity         INT           NOT NULL,
    unit_price       DECIMAL(10,2) NOT NULL
);
