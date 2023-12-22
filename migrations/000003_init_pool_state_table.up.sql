CREATE TABLE pool_state (
    id SERIAL PRIMARY KEY,
    pool_address TEXT NOT NULL,
    start_block BIGINT NOT NULL,
    stop_block BIGINT NOT NULL,
    block_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    total_fees BIGINT DEFAULT 0,
    total_liquidity BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_pool_address
        FOREIGN KEY (pool_address)
        REFERENCES pool_address (address)
        ON DELETE CASCADE
);

