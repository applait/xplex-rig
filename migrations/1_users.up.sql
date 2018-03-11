create table user_accounts (
  id uuid not null primary key,
  username text not null,
  email text not null,
  password text,
  is_active boolean not null default false
);
