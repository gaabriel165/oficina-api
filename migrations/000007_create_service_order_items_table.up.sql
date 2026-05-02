CREATE TABLE service_order_items (
    id               UUID PRIMARY KEY,
    service_order_id UUID          NOT NULL REFERENCES service_orders(id) ON DELETE CASCADE,
    service_id       UUID          NOT NULL REFERENCES services(id),
    service_name     VARCHAR(255)  NOT NULL,
    labor_price      DECIMAL(10,2) NOT NULL
);
