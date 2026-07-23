-- Password for both seeded users: password123  (bcrypt, cost 10)

INSERT INTO users (email, password_hash) VALUES
    ('jeremy@example.com', '$2b$10$IfrZadPMDw2YB3WPQW1ZVeSrpgmu9SXjJKYk.2lsfxNZh3U1ACD6u'),
    ('vijay@example.com',  '$2b$10$a28RE9dTn.fd4EDOj.QoouGtxDHcm3SOY8bXZCjMlMJrJ3s2IIdUu');

INSERT INTO orders (user_id, description, total_amount, status, transaction_id, created_at)
SELECT u.id, v.description, v.total_amount, v.status, v.transaction_id, NOW() - v.age
FROM users u
JOIN (VALUES
    ('jeremy@example.com', 'Mechanical keyboard',   1250000.00, 'completed',  'txn_mock_0001', INTERVAL '3 days'),
    ('jeremy@example.com', 'Standing desk',         4800000.00, 'processing', NULL,            INTERVAL '2 days'),
    ('jeremy@example.com', 'USB-C hub',              350000.00, 'failed',     'txn_mock_0002', INTERVAL '1 day'),
    ('jeremy@example.com', 'Monitor arm',            920000.00, 'pending',    NULL,            INTERVAL '2 hours'),
    ('vijay@example.com',  'Noise cancelling buds', 2100000.00, 'completed',  'txn_mock_0003', INTERVAL '5 days'),
    ('vijay@example.com',  'Laptop sleeve',          275000.00, 'pending',    NULL,            INTERVAL '30 minutes')
) AS v(email, description, total_amount, status, transaction_id, age)
  ON u.email = v.email;

INSERT INTO order_logs (order_id, event_id, event_type, payload)
SELECT o.id, 'evt_seed_0001', 'payment.completed',
       jsonb_build_object('transaction_id', o.transaction_id, 'status', 'completed')
FROM orders o WHERE o.transaction_id = 'txn_mock_0001';