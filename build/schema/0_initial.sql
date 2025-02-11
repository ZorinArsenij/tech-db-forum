CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS client, forum, thread, post, vote, forum_client;

-- Client

CREATE UNLOGGED TABLE IF NOT EXISTS client (
  id SERIAL PRIMARY KEY,
  email CITEXT NOT NULL UNIQUE,
  nickname CITEXT NOT NULL UNIQUE,
  fullname TEXT NOT NULL,
  about TEXT NOT NULL DEFAULT ''
) WITH (autovacuum_enabled = FALSE);

CREATE INDEX IF NOT EXISTS client_email_index
  ON client(email);

CREATE INDEX IF NOT EXISTS client_covering_index
  ON client(nickname) INCLUDE (email, fullname, about, id);

-- Forum

CREATE UNLOGGED TABLE IF NOT EXISTS forum (
  id SERIAL PRIMARY KEY,
  slug CITEXT NOT NULL,
  title TEXT NOT NULL,
  threads INTEGER NOT NULL DEFAULT 0,
  posts BIGINT NOT NULL DEFAULT 0,
  user_nickname CITEXT NOT NULL
) WITH (autovacuum_enabled = FALSE);

CREATE INDEX IF NOT EXISTS forum_slug_index
  ON forum(slug) INCLUDE (title, posts, threads, user_nickname, id);

-- Thread

CREATE UNLOGGED TABLE IF NOT EXISTS thread (
  id SERIAL PRIMARY KEY,
  slug CITEXT DEFAULT NULL,
  title TEXT NOT NULL,
  message TEXT NULL,
  forum_id INTEGER NOT NULL,
  forum_slug CITEXT NOT NULL,
  user_nickname CITEXT NOT NULL,
  created TIMESTAMPTZ,
  votes INTEGER NOT NULL DEFAULT 0
) WITH (autovacuum_enabled = FALSE);

CREATE INDEX IF NOT EXISTS thread_slug_index
  ON thread(slug) INCLUDE (id, title, message, forum_slug, user_nickname, created, votes);

CREATE INDEX IF NOT EXISTS thread_func_id_index
  ON thread(text(id));

CREATE INDEX IF NOT EXISTS thread_created_index
  ON thread(forum_slug, created);

-- Post

CREATE UNLOGGED TABLE IF NOT EXISTS post (
  id SERIAL PRIMARY KEY,
  message TEXT NOT NULL,
  created TIMESTAMPTZ,
  is_edited BOOLEAN NOT NULL DEFAULT FALSE,
  user_nickname CITEXT NOT NULL,
  thread_id INTEGER NOT NULL,
  forum_slug CITEXT NOT NULL,
  parent INT DEFAULT 0,
  parents INT [] NOT NULL,
  root INT NOT NULL
) WITH (autovacuum_enabled = FALSE);

CREATE INDEX IF NOT EXISTS post_id_thread_index
  ON post(id, thread_id);

CREATE INDEX IF NOT EXISTS post_tree_index
  ON post(thread_id, array_append(parents, id));

CREATE INDEX IF NOT EXISTS post_parent_tree_index
  ON post(thread_id, id) WHERE parent = 0;

CREATE INDEX IF NOT EXISTS post_thread_id_index
  ON post(thread_id, id);

CREATE INDEX IF NOT EXISTS post_root_parents_func_index
  ON post(root, array_append(parents, id));

-- Vote

CREATE UNLOGGED TABLE IF NOT EXISTS vote (
  id SERIAL PRIMARY KEY,
  voice BOOLEAN,
  user_nickname CITEXT NOT NULL,
  thread_id INTEGER NOT NULL
) WITH (autovacuum_enabled = FALSE);

CREATE INDEX IF NOT EXISTS vote_user_nickname_thread_id_index
  ON vote(user_nickname, thread_id);

-- Forum client

CREATE UNLOGGED TABLE IF NOT EXISTS forum_client (
  forum_slug CITEXT NOT NULL,
  email CITEXT NOT NULL,
  nickname CITEXT NOT NULL,
  fullname TEXT NOT NULL,
  about TEXT NOT NULL DEFAULT ''
) WITH (autovacuum_enabled = FALSE);

CREATE UNIQUE INDEX IF NOT EXISTS forum_client_index
  ON forum_client (forum_slug, nickname);

CREATE INDEX IF NOT EXISTS forum_client_covering_index
  ON forum_client (forum_slug, nickname, email, fullname, about);