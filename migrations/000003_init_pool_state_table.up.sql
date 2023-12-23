CREATE TABLE
    pool_state (
        id SERIAL PRIMARY KEY,
        pool_address TEXT NOT NULL,
        date TIMESTAMP UNIQUE NOT NULL,
        tvl_usd TEXT,
        fees_usd TEXT,
        created_at TIMESTAMP
        WITH
            TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
            CONSTRAINT fk_pool_address FOREIGN KEY (pool_address) REFERENCES pool_address (address) ON DELETE CASCADE
    );