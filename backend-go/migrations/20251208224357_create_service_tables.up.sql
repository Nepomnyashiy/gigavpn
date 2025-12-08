-- Subscriptions Table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan_id INT NOT NULL REFERENCES plans(id),
    start_date TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    auto_renew BOOLEAN NOT NULL DEFAULT FALSE
);

-- Servers Table
CREATE TABLE IF NOT EXISTS servers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ip_address INET UNIQUE NOT NULL,
    country_code VARCHAR(10) NOT NULL,
    api_url VARCHAR(255),
    xray_config JSONB,
    capacity_limit INT NOT NULL DEFAULT 100,
    current_load INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Access Keys Table (VLESS UUIDs)
CREATE TABLE IF NOT EXISTS access_keys (
    id UUID PRIMARY KEY, -- Это VLESS UUID, он не генерируется в БД
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    server_id INT NOT NULL REFERENCES servers(id),
    key_private TEXT NOT NULL,
    key_public TEXT NOT NULL,
    access_link TEXT NOT NULL,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE
);

-- Transactions Table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    amount BIGINT NOT NULL,
    source VARCHAR(50),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Traffic Logs Table
CREATE TABLE IF NOT EXISTS traffic_logs (
    id BIGSERIAL PRIMARY KEY,
    access_key_id UUID NOT NULL REFERENCES access_keys(id) ON DELETE CASCADE,
    server_id INT NOT NULL REFERENCES servers(id),
    upload_bytes BIGINT NOT NULL,
    download_bytes BIGINT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL
);

-- Foreign Key Indexes
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_access_keys_subscription_id ON access_keys(subscription_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_traffic_logs_access_key_id ON traffic_logs(access_key_id);
