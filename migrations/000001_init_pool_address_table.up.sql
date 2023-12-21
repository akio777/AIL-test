CREATE TABLE pool_address (
    id SERIAL NOT NULL,
    address TEXT NOT NULL PRIMARY KEY,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Optional: Create a trigger to update the 'updated_at' field whenever a row is updated
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trg_pool_address_updated_at
BEFORE UPDATE ON pool_address
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();