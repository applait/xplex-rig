-- Custom sequence to explicitly handle destinations counter
create sequence destinations_seq
  start with 1
  increment by 1
  no minvalue
  no maxvalue
  cache 1
  cycle;

-- The table for destinations
create table destinations (
  -- ID of the destionation. This is more for internal use, so we use an auto incrementing serial here
  id bigint default nextval('destinations_seq') not null primary key,

  -- Name of the service. This will be something like "YouTube" or "Twitch"
  service text not null,

  -- Stream key of the destination. This will be the streaming key given by the user for the service
  stream_key text not null,

  -- Which specific server to push the stream to
  server text not null default 'default',

  -- Whether this destination is active or not
  is_active boolean not null default false,

  -- The multi stream to which this destination belongs
  multi_stream_id uuid not null references multi_streams (id)
);
