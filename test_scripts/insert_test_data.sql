TRUNCATE messages;
TRUNCATE chat_users;
TRUNCATE chats CASCADE;
TRUNCATE users CASCADE;

--users
INSERT INTO users (username, created_at)
VALUES ('User 1', 100);

INSERT INTO users (username, created_at)
VALUES ('User 2', 200);

INSERT INTO users (username, created_at)
VALUES ('User 3', 300);

-- chats
INSERT INTO chats (name, created_at)
VALUES ('Chat 1', 400);

INSERT INTO chats (name, created_at)
VALUES ('Chat 2', 500);

--chat_users
INSERT INTO chat_users (chat_id, user_id)
VALUES (1, 1);

INSERT INTO chat_users (chat_id, user_id)
VALUES (1, 2);

INSERT INTO chat_users (chat_id, user_id)
VALUES (2, 1);

INSERT INTO chat_users (chat_id, user_id)
VALUES (2, 2);

INSERT INTO chat_users (chat_id, user_id)
VALUES (2, 3);

-- messages
INSERT INTO messages (chat_id, author_id, text, created_at)
VALUES (1, 2, 'Message 1 in chat 1', 600);

INSERT INTO messages (chat_id, author_id, text, created_at)
VALUES (1, 1, 'Message 2 in chat 1', 650);

INSERT INTO messages (chat_id, author_id, text, created_at)
VALUES (1, 2, 'Message 3 in chat 1', 700);

INSERT INTO messages (chat_id, author_id, text, created_at)
VALUES (2, 2, 'Message 1 in chat 2', 550);

INSERT INTO messages (chat_id, author_id, text, created_at)
VALUES (2, 1, 'Message 2 in chat 2', 600);

INSERT INTO messages (chat_id, author_id, text, created_at)
VALUES (2, 3, 'Message 3 in chat 2', 650);
