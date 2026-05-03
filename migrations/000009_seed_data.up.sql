-- users
INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES
  ('c44ade7e-de0b-479d-9bc9-32477760a125', 'admin@oficina.com', '$2a$10$P5fXxzPcWqi.XAKQfQL8HOL2fVuTR0J.HXsV29ALMsk7MbYzAdR06', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- customers
INSERT INTO customers (id, name, document, document_type, phone, email, created_at, updated_at) VALUES
  ('16b78307-e693-45c5-bf08-7ee50cf6b922', 'João Carlos Silva',          '12345678909',   'cpf',  '11998765432', 'joao.silva@email.com',                NOW(), NOW()),
  ('d825ab74-8702-4a74-9dc8-3d5548d5611d', 'Transportes Paulistas Ltda', '12345678000195','cnpj', '1132589700',  'contato@transportespaulistas.com.br', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- vehicles
INSERT INTO vehicles (id, customer_id, plate, brand, model, year, created_at, updated_at) VALUES
  ('db2cb4ad-3637-49ea-987e-d1ef1d35c5dc', '16b78307-e693-45c5-bf08-7ee50cf6b922', 'ABC1D23', 'Honda', 'Civic', 2022, NOW(), NOW()),
  ('a3ab0dfb-87a8-4aaf-97df-7c41672d9a6b', '16b78307-e693-45c5-bf08-7ee50cf6b922', 'XYZ9W87', 'Ford',  'Ka',    2019, NOW(), NOW()),
  ('8fe7ce41-3f71-48f9-b288-89aa9270c364', 'd825ab74-8702-4a74-9dc8-3d5548d5611d', 'DEF4E56', 'Fiat',  'Uno',   2020, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- parts
INSERT INTO parts (id, name, description, unit_price, stock_quantity, created_at, updated_at) VALUES
  ('b9268866-0ac0-484c-a024-0e797a0474b3', 'Filtro de óleo',             'Filtro de óleo para motor a gasolina',    25.00, 20, NOW(), NOW()),
  ('f46a1ae9-d02a-4049-946d-4a000deb4f8b', 'Óleo de motor 5W30',         'Óleo sintético 5W30 — unidade 1L',        35.00, 50, NOW(), NOW()),
  ('9e0f3948-8859-435b-8b60-cf0e22d60371', 'Pastilha de freio dianteira','Kit com 4 pastilhas dianteiras',          89.00, 15, NOW(), NOW()),
  ('4e6fa80b-a091-4b27-98a9-42becd44054e', 'Fluido de freio DOT4',       'Fluido de freio DOT4 — 500ml',            28.00, 25, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- services
INSERT INTO services (id, name, description, labor_price, estimated_minutes, created_at, updated_at) VALUES
  ('5c457d6c-4c07-43da-a227-3bf0292b68cd', 'Troca de óleo',               'Troca de óleo e filtro de óleo',                      80.00,  60, NOW(), NOW()),
  ('5f44a580-5af6-4d9d-a215-d1a4cec68ad5', 'Alinhamento e balanceamento', 'Alinhamento das 4 rodas e balanceamento',             120.00,  90, NOW(), NOW()),
  ('0e897709-370a-43ef-9b84-818e28156326', 'Revisão completa de freios',  'Verificação e ajuste completo do sistema de freios',  200.00, 120, NOW(), NOW()),
  ('762c7fd4-0a12-4bd0-9bf3-5a7b80a3569b', 'Troca de pastilhas de freio', 'Substituição das pastilhas de freio dianteiras',      150.00,  90, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- service orders
INSERT INTO service_orders (id, customer_id, vehicle_id, status, notes, total_amount, started_at, finished_at, created_at, updated_at) VALUES
  (
    '92a38a07-605b-4255-99de-11fe9252399f',
    '16b78307-e693-45c5-bf08-7ee50cf6b922',
    'db2cb4ad-3637-49ea-987e-d1ef1d35c5dc',
    'delivered',
    'Cliente relatou barulho nos freios e consumo elevado de óleo.',
    534.00,
    NOW() - INTERVAL '5 days',
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '7 days',
    NOW() - INTERVAL '2 days'
  ),
  (
    '9dd878c0-8118-4c8b-a5d7-9945baed2e65',
    'd825ab74-8702-4a74-9dc8-3d5548d5611d',
    '8fe7ce41-3f71-48f9-b288-89aa9270c364',
    'in_execution',
    'Veículo puxando para o lado direito na frenagem.',
    120.00,
    NOW() - INTERVAL '1 day',
    NULL,
    NOW() - INTERVAL '2 days',
    NOW() - INTERVAL '1 day'
  ),
  (
    'd01a1943-60da-4f8c-9aa1-9657ae6b2de4',
    '16b78307-e693-45c5-bf08-7ee50cf6b922',
    'a3ab0dfb-87a8-4aaf-97df-7c41672d9a6b',
    'waiting_approval',
    'Troca preventiva das pastilhas de freio dianteiras.',
    328.00,
    NULL,
    NULL,
    NOW() - INTERVAL '1 day',
    NOW()
  ),
  (
    '30956ff8-7d7d-41e9-9992-23276e73e8a6',
    '16b78307-e693-45c5-bf08-7ee50cf6b922',
    'db2cb4ad-3637-49ea-987e-d1ef1d35c5dc',
    'received',
    'Revisão geral solicitada pelo cliente.',
    0.00,
    NULL,
    NULL,
    NOW(),
    NOW()
  )
ON CONFLICT DO NOTHING;

-- service order items
INSERT INTO service_order_items (id, service_order_id, service_id, service_name, labor_price) VALUES
  ('f78fc117-fd57-4382-b59d-59bebf1db91c', '92a38a07-605b-4255-99de-11fe9252399f', '5c457d6c-4c07-43da-a227-3bf0292b68cd', 'Troca de óleo',               80.00),
  ('0cd1b385-a37a-478c-b239-27ed8060bca9', '92a38a07-605b-4255-99de-11fe9252399f', '0e897709-370a-43ef-9b84-818e28156326', 'Revisão completa de freios',  200.00),
  ('bced97a4-a31c-4b88-9bbc-00dd580da0db', '9dd878c0-8118-4c8b-a5d7-9945baed2e65', '5f44a580-5af6-4d9d-a215-d1a4cec68ad5', 'Alinhamento e balanceamento', 120.00),
  ('a5df5f86-5c50-4dd5-a20b-516bacb23b8a', 'd01a1943-60da-4f8c-9aa1-9657ae6b2de4', '762c7fd4-0a12-4bd0-9bf3-5a7b80a3569b', 'Troca de pastilhas de freio', 150.00)
ON CONFLICT DO NOTHING;

-- service order parts
INSERT INTO service_order_parts (id, service_order_id, part_id, part_name, quantity, unit_price) VALUES
  ('afd334f3-9f70-4908-9826-7772980b31a4', '92a38a07-605b-4255-99de-11fe9252399f', 'b9268866-0ac0-484c-a024-0e797a0474b3', 'Filtro de óleo',              1, 25.00),
  ('18eb65a2-2703-4486-9333-3b8f31a5d6b5', '92a38a07-605b-4255-99de-11fe9252399f', 'f46a1ae9-d02a-4049-946d-4a000deb4f8b', 'Óleo de motor 5W30',          4, 35.00),
  ('b4080aad-24da-4a42-93da-6c1ed91aa91d', '92a38a07-605b-4255-99de-11fe9252399f', '9e0f3948-8859-435b-8b60-cf0e22d60371', 'Pastilha de freio dianteira', 1, 89.00),
  ('e0e38bfe-0916-412e-9558-4f7fe8dfafca', 'd01a1943-60da-4f8c-9aa1-9657ae6b2de4', '9e0f3948-8859-435b-8b60-cf0e22d60371', 'Pastilha de freio dianteira', 2, 89.00)
ON CONFLICT DO NOTHING;
