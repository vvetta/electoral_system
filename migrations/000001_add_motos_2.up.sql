CREATE TABLE motos (
  id INT PRIMARY KEY,
  year INT NOT NULL,
  mileage INT NOT NULL,
  engine_size INT NOT NULL,
  moto_type VARCHAR(255) NOT NULL,
  location VARCHAR(255) NOT NULL,
  price INT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_motos_deleted_at ON motos (deleted_at);


