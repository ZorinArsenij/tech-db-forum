CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS forum, thread, client, post, vote;

CREATE TABLE client (
  id SERIAL PRIMARY KEY,
  email CITEXT NOT NULL UNIQUE,
  nickname CITEXT NOT NULL UNIQUE,
  fullname TEXT NOT NULL,
  about TEXT NOT NULL DEFAULT ''
);

CREATE INDEX client_covering_index
  ON client(nickname, email, fullname, about);

CREATE INDEX client_email_index
  ON client(email);

CREATE INDEX client_nickname_index
  ON client(nickname);

CREATE TABLE forum (
  id SERIAL PRIMARY KEY,
  slug CITEXT NOT NULL UNIQUE,
  title TEXT NOT NULL,
  threads INTEGER NOT NULL DEFAULT 0,
  posts BIGINT NOT NULL DEFAULT 0,
  user_id INTEGER NOT NULL REFERENCES client(id),
  user_nickname CITEXT NOT NULL
);

CREATE INDEX forum_slug_index
  ON forum(slug);

CREATE INDEX forum_covering_index
  ON forum(slug, title, posts, threads, user_nickname);

CREATE TABLE thread (
  id SERIAL PRIMARY KEY,
  slug CITEXT DEFAULT NULL,
  title TEXT NOT NULL,
  message TEXT NULL,
  forum_id INTEGER NOT NULL REFERENCES forum(id),
  forum_slug CITEXT NOT NULL,
  user_id INTEGER NOT NULL REFERENCES client(id),
  user_nickname CITEXT NOT NULL,
  created TIMESTAMPTZ,
  votes INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX thread_slug_index
  ON thread(slug);

CREATE INDEX thread_func_id_index
  ON thread(text(id));

CREATE INDEX thread_created_index
  ON thread(forum_slug, created);

CREATE INDEX thread_created_desc_index
  ON thread(forum_slug, created DESC);

-- ???
CREATE INDEX thread_covering_index
  ON thread(id, slug, title, message, forum_slug, user_nickname, created, votes);

CREATE TABLE post (
  id SERIAL PRIMARY KEY,
  message TEXT NOT NULL,
  created TIMESTAMPTZ,
  is_edited BOOLEAN NOT NULL DEFAULT FALSE,
  user_id INTEGER NOT NULL REFERENCES client(id),
  user_nickname CITEXT NOT NULL,
  thread_id INTEGER NOT NULL REFERENCES thread(id),
  forum_slug CITEXT NOT NULL REFERENCES forum(slug),
  parent INT DEFAULT 0,
  parents INT [] NOT NULL,
  root INT NOT NULL
);

CREATE INDEX post_id_thread_index
  ON post(id, thread_id);

CREATE INDEX post_tree_index
  ON post(thread_id, array_append(parents, id));

CREATE INDEX post_tree_desc_index
  ON post(thread_id, array_append(parents, id) DESC);

CREATE INDEX post_parent_tree_index
  ON post(thread_id, id) WHERE parent = 0;

CREATE INDEX post_parent_tree_desc_index
  ON post(thread_id, id DESC) WHERE parent = 0;

CREATE INDEX post_thread_index
  ON post(thread_id, id);

CREATE INDEX post_thread_desc_index
  ON post(thread_id, id DESC);


CREATE TABLE vote (
  id SERIAL PRIMARY KEY,
  voice BOOLEAN,
  user_id INTEGER NOT NULL REFERENCES client(id),
  thread_id INTEGER NOT NULL REFERENCES thread(id)
);