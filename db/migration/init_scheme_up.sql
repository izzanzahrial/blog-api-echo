CREATE TABLE posts (
    post_id serial PRIMARY KEY,
    title VARCHAR (255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_post_title ON posts(title);

CREATE TABLE users (
    user_id serial PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    password VARCHAR (255) NOT NULL
);

CREATE INDEX idx_user ON users(user_id, password);

CREATE TABLE favourites (
    fav_id serial PRIMARY KEY,
    FOREIGN KEY (post_id) REFERENCES posts (post_id),
    FOREIGN KEY (user_id) REFERENCES users (user_id)
)