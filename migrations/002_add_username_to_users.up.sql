-- Add username column to users table
ALTER TABLE users ADD COLUMN username VARCHAR(255) UNIQUE;

-- Create index for faster lookups
CREATE INDEX idx_users_username ON users(username);

-- Update existing users: set username = email (without @domain) as a fallback
UPDATE users SET username = split_part(email, '@', 1) WHERE username IS NULL;