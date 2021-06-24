CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) NOT NULL,
    created_at BIGINT NOT NULL CHECK (created_at > 0)
);

CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,
    created_at BIGINT NOT NULL CHECK (created_at > 0)
);

CREATE TABLE If NOT EXISTS chat_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID REFERENCES chats (id) ON DELETE CASCADE NOT NULL,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    chat_id UUID REFERENCES chats (id) ON DELETE CASCADE NOT NULL,
    author_id UUID REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    text TEXT NOT NULL,
    created_at BIGINT NOT NULL CHECK (created_at > 0)
);