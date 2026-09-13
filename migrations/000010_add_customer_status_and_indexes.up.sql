ALTER TABLE customers ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'active';

UPDATE customers SET document_type = UPPER(document_type);

CREATE INDEX idx_customers_status ON customers (status);
CREATE INDEX idx_vehicles_customer_id ON vehicles (customer_id);
CREATE INDEX idx_service_orders_status ON service_orders (status);
CREATE INDEX idx_service_orders_customer_id ON service_orders (customer_id);
CREATE INDEX idx_service_orders_vehicle_id ON service_orders (vehicle_id);
CREATE INDEX idx_service_order_items_service_order_id ON service_order_items (service_order_id);
CREATE INDEX idx_service_order_parts_service_order_id ON service_order_parts (service_order_id);
