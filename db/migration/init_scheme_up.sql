CREATE TABLE posts (
    post_id serial PRIMARY KEY,
    title VARCHAR (200) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_post_title ON posts(title);