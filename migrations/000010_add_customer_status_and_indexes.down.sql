DROP INDEX IF EXISTS idx_service_order_parts_service_order_id;
DROP INDEX IF EXISTS idx_service_order_items_service_order_id;
DROP INDEX IF EXISTS idx_service_orders_vehicle_id;
DROP INDEX IF EXISTS idx_service_orders_customer_id;
DROP INDEX IF EXISTS idx_service_orders_status;
DROP INDEX IF EXISTS idx_vehicles_customer_id;
DROP INDEX IF EXISTS idx_customers_status;

ALTER TABLE customers DROP COLUMN status;
