create table user_accounts (
  -- Primary key is UUID here
  id uuid not null primary key,

  -- username needs to be unique
  username text not null unique,

  -- User email needs to be unique as well
  email text not null unique,

  -- Hashed password
  password text, -- leaving this as nullable because I'm not sure what happens for external auth.

  -- Denotes if user is active
  is_active boolean not null default false,

  -- Timestamps
  created_at timestamp with time zone not null,
  updated_at timestamp with time zone -- is null if row was never updated
);
