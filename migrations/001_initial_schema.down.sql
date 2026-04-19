-- Drop initial schema
DROP TRIGGER IF EXISTS update_api_keys_updated_at ON api_keys;
DROP TRIGGER IF EXISTS update_links_updated_at ON links;
DROP TRIGGER IF EXISTS update_domains_updated_at ON domains;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS click_events;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS links;
DROP TABLE IF EXISTS domains;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";