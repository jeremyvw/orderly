DELETE FROM orders
WHERE user_id IN (
    SELECT id FROM users WHERE email IN ('jeremy@example.com', 'vijay@example.com')
);

DELETE FROM users
WHERE email IN ('jeremy@example.com', 'vijay@example.com');