-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wallet_address TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE,
    username TEXT UNIQUE,
    company TEXT,
    business_size TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Businesses table
CREATE TABLE IF NOT EXISTS businesses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    industry TEXT,
    location TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Market insights table
CREATE TABLE IF NOT EXISTS market_insights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    business_id INTEGER,
    industry TEXT NOT NULL,
    insight_type TEXT NOT NULL,
    data TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE
);

-- Transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    tx_hash TEXT UNIQUE NOT NULL,
    amount REAL NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'confirmed', 'failed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- API access logs table
CREATE TABLE IF NOT EXISTS api_access_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    endpoint TEXT NOT NULL,
    request_data TEXT,
    response_data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Feedback table
CREATE TABLE IF NOT EXISTS feedback (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    insight_id INTEGER NOT NULL,
    rating INTEGER CHECK(rating BETWEEN 1 AND 5),
    comments TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (insight_id) REFERENCES market_insights(id) ON DELETE CASCADE
);

-- Smart contract metadata table
CREATE TABLE IF NOT EXISTS smart_contracts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    contract_address TEXT UNIQUE NOT NULL,
    contract_name TEXT NOT NULL,
    contract_abi TEXT NOT NULL,
    contract_bytecode TEXT NOT NULL,
    deployed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    owner_user_id INTEGER,
    network TEXT NOT NULL,
    FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Business subscription payments
CREATE TABLE IF NOT EXISTS subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    business_id INTEGER NOT NULL,
    plan_type TEXT NOT NULL CHECK(plan_type IN ('basic', 'premium', 'enterprise')),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    payment_tx_hash TEXT UNIQUE,
    amount REAL NOT NULL,
    status TEXT NOT NULL CHECK(status IN ('active', 'expired', 'canceled')),
    FOREIGN KEY (business_id) REFERENCES businesses(id) ON DELETE CASCADE,
    FOREIGN KEY (payment_tx_hash) REFERENCES transactions(tx_hash) ON DELETE SET NULL
);

-- Contract events table
CREATE TABLE IF NOT EXISTS contract_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    contract_id INTEGER NOT NULL,
    event_name TEXT NOT NULL,
    event_data TEXT NOT NULL,
    block_number INTEGER NOT NULL,
    transaction_hash TEXT NOT NULL,
    log_index INTEGER NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (contract_id) REFERENCES smart_contracts(id) ON DELETE CASCADE
);

-- Create indices for performance
CREATE INDEX IF NOT EXISTS idx_users_wallet ON users(wallet_address);
CREATE INDEX IF NOT EXISTS idx_businesses_user ON businesses(user_id);
CREATE INDEX IF NOT EXISTS idx_market_insights_business ON market_insights(business_id);
CREATE INDEX IF NOT EXISTS idx_market_insights_industry ON market_insights(industry);
CREATE INDEX IF NOT EXISTS idx_transactions_user ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_hash ON transactions(tx_hash);
CREATE INDEX IF NOT EXISTS idx_feedback_insight ON feedback(insight_id);
CREATE INDEX IF NOT EXISTS idx_contract_events_contract ON contract_events(contract_id);
