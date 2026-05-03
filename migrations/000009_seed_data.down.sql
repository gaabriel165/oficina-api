DELETE FROM service_order_parts WHERE id IN (
  'afd334f3-9f70-4908-9826-7772980b31a4',
  '18eb65a2-2703-4486-9333-3b8f31a5d6b5',
  'b4080aad-24da-4a42-93da-6c1ed91aa91d',
  'e0e38bfe-0916-412e-9558-4f7fe8dfafca'
);
DELETE FROM service_order_items WHERE id IN (
  'f78fc117-fd57-4382-b59d-59bebf1db91c',
  '0cd1b385-a37a-478c-b239-27ed8060bca9',
  'bced97a4-a31c-4b88-9bbc-00dd580da0db',
  'a5df5f86-5c50-4dd5-a20b-516bacb23b8a'
);
DELETE FROM service_orders WHERE id IN (
  '92a38a07-605b-4255-99de-11fe9252399f',
  '9dd878c0-8118-4c8b-a5d7-9945baed2e65',
  'd01a1943-60da-4f8c-9aa1-9657ae6b2de4',
  '30956ff8-7d7d-41e9-9992-23276e73e8a6'
);
DELETE FROM services WHERE id IN (
  '5c457d6c-4c07-43da-a227-3bf0292b68cd',
  '5f44a580-5af6-4d9d-a215-d1a4cec68ad5',
  '0e897709-370a-43ef-9b84-818e28156326',
  '762c7fd4-0a12-4bd0-9bf3-5a7b80a3569b'
);
DELETE FROM parts WHERE id IN (
  'b9268866-0ac0-484c-a024-0e797a0474b3',
  'f46a1ae9-d02a-4049-946d-4a000deb4f8b',
  '9e0f3948-8859-435b-8b60-cf0e22d60371',
  '4e6fa80b-a091-4b27-98a9-42becd44054e'
);
DELETE FROM vehicles WHERE id IN (
  'db2cb4ad-3637-49ea-987e-d1ef1d35c5dc',
  'a3ab0dfb-87a8-4aaf-97df-7c41672d9a6b',
  '8fe7ce41-3f71-48f9-b288-89aa9270c364'
);
DELETE FROM customers WHERE id IN (
  '16b78307-e693-45c5-bf08-7ee50cf6b922',
  'd825ab74-8702-4a74-9dc8-3d5548d5611d'
);
DELETE FROM users WHERE id = 'c44ade7e-de0b-479d-9bc9-32477760a125';
