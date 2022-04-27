CREATE TABLE posts (
    post_id serial PRIMARY KEY,
    title VARCHAR (255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- tsvector = values that stored in ordered list of distinct words
-- setweight = to weight the value of the data
-- coalesce = function that will return the first value that's not null 
-- https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-coalesce/
-- full text search postgres = https://blog.crunchydata.com/blog/postgres-full-text-search-a-search-engine-in-a-database
ALTER TABLE posts ADD COLUMN ts_title_content tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(content, '')), 'B') STORED;
)

CREATE INDEX idx_post_title ON posts(title);

CREATE TABLE users (
    user_id serial PRIMARY KEY,
    name VARCHAR (255) NOT NULL,
    password VARCHAR (255) NOT NULL
);

CREATE INDEX idx_user ON users(user_id, password);

CREATE TABLE favourites (
    FOREIGN KEY (user_id) REFERENCES users (user_id) ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts (post_id) ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT favourites_pkey PRIMARY KEY (user_id, post_id) -- explicit pk
)