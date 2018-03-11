create table multi_streams (
  -- ID of the stream
  id uuid primary key,

  -- Streaming key used by user
  stream_key text not null unique,

  -- Whether stream is active
  is_active boolean not null default false,

  -- Whether stream is currently live
  is_streaming boolean not null default false,

  -- User Account of this stream
  user_account_id uuid not null references user_accounts (id),

  -- Timestamps
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone
);
