-- Drop the trigger if it exists
DROP TRIGGER IF EXISTS trg_pool_address_updated_at ON pool_address;

-- Drop the function if it exists
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop the table if it exists
DROP TABLE IF EXISTS pool_address;