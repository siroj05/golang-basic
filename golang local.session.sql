
-- SELECT * FROM customer;

-- DELETE FROM customer;

-- ALTER TABLE customer
--   ADD COLUMN email VARCHAR(100),
--   ADD COLUMN balance INT DEFAULT 0,
--   ADD COLUMN rating DOUBLE DEFAULT 0.0,
--   ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--   ADD COLUMN birth_date DATE,
--   ADD COLUMN marriage BOOLEAN DEFAULT false;

-- DESC customer;

-- INSERT INTO customer(id, name, email, balance, rating, birth_date, marriage)
-- VALUES ('eko', 'EKO', 'eko@gmail.com', 100000, 5.0, '1999-9-9', true);

-- INSERT INTO customer(id, name, email, balance, rating, birth_date, marriage)
-- VALUES ('budi', 'Budi', 'budi@gmail.com', 100000, 5.0, '1999-9-9', true);

-- INSERT INTO customer(id, name, email, balance, rating, birth_date, marriage)
-- VALUES ('budi', 'Budi', null, 750000, 88.5, '1999-9-9', false);

UPDATE customer
  SET email = NULL,
  birth_date = NULL,
  WHERE id = 'eko'
