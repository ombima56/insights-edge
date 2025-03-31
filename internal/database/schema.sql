CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	password_hash TEXT NOT NULL,
    account_type TEXT NOT NULL,
	company_name TEXT,
	industry TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	token TEXT UNIQUE NOT NULL,
	expires_at DATETIME NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE market_insights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    industry TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    trend_value REAL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ratings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    insight_id INT NOT NULL REFERENCES market_insights(id) ON DELETE CASCADE,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE rewards (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    insight_id INT NOT NULL REFERENCES market_insights(id) ON DELETE CASCADE,
    tokens_awarded DECIMAL(10,2) NOT NULL,
    transaction_hash VARCHAR(66) UNIQUE NOT NULL, -- TX hash from blockchain
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);