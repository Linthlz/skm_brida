ALTER TABLE users
ALTER COLUMN verification_token TYPE UUID USING verification_token::uuid;