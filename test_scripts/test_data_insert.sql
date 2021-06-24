TRUNCATE messages;
TRUNCATE chat_users;
TRUNCATE chats CASCADE;
TRUNCATE users CASCADE;

CREATE OR REPLACE FUNCTION data_insert() RETURNS void AS $$
    DECLARE
        user1_id UUID;
        user2_id UUID;
        user3_id UUID;
        chat1_id UUID;
        chat2_id UUID;
    BEGIN
        --users
        INSERT INTO users (username, created_at)
        VALUES ('User 1', 100) RETURNING id INTO user1_id;

        INSERT INTO users (username, created_at)
        VALUES ('User 2', 200) RETURNING id INTO user2_id;

        INSERT INTO users (username, created_at)
        VALUES ('User 3', 300) RETURNING id INTO user3_id;

        -- chats
        INSERT INTO chats (name, created_at)
        VALUES ('Chat 1', 400) RETURNING id INTO chat1_id;

        INSERT INTO chats (name, created_at)
        VALUES ('Chat 2', 500) RETURNING id INTO chat2_id;

        --chat_users
        INSERT INTO chat_users (chat_id, user_id)
        VALUES (chat1_id, user1_id);

        INSERT INTO chat_users (chat_id, user_id)
        VALUES (chat1_id, user2_id);

        INSERT INTO chat_users (chat_id, user_id)
        VALUES (chat2_id, user1_id);

        INSERT INTO chat_users (chat_id, user_id)
        VALUES (chat2_id, user2_id);

        INSERT INTO chat_users (chat_id, user_id)
        VALUES (chat2_id, user3_id);

        -- messages
        INSERT INTO messages (chat_id, author_id, text, created_at)
        VALUES (chat1_id, user2_id, 'Message 1 in chat 1', 600);

        INSERT INTO messages (chat_id, author_id, text, created_at)
        VALUES (chat1_id, user1_id, 'Message 2 in chat 1', 650);

        INSERT INTO messages (chat_id, author_id, text, created_at)
        VALUES (chat1_id, user2_id, 'Message 3 in chat 1', 700);

        INSERT INTO messages (chat_id, author_id, text, created_at)
        VALUES (chat2_id, user2_id, 'Message 1 in chat 2', 550);

        INSERT INTO messages (chat_id, author_id, text, created_at)
        VALUES (chat2_id, user1_id, 'Message 2 in chat 2', 600);

        INSERT INTO messages (chat_id, author_id, text, created_at)
        VALUES (chat2_id, user3_id, 'Message 3 in chat 2', 650);

        RETURN;
    END;
$$ LANGUAGE plpgsql;

SELECT data_insert()




