CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    username varchar(50) NOT NULL,
    created_at bigint CHECK (created_at > 0) NOT NULL 
);

CREATE TABLE IF NOT EXISTS chats (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL,
    created_at bigint CHECK (created_at > 0) NOT NULL 
);

CREATE TABLE If NOT EXISTS chat_users (
    id serial PRIMARY KEY,
    chat_id INT REFERENCES chats (id) ON DELETE CASCADE NOT NULL,
    user_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
    id serial PRIMARY KEY,
    chat_id INT REFERENCES chats (id) ON DELETE CASCADE NOT NULL,
    author_id INT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    text TEXT NOT NULL,
    created_at bigint CHECK (created_at > 0) NOT NULL 
);