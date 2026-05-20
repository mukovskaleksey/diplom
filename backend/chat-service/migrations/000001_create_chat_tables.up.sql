
create table if not exists chats (
                                     id bigserial primary key,
                                     ticket_id bigint not null unique,
                                     created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
    );

create table if not exists messages (
                                        id bigserial primary key,
                                        chat_id bigint not null references chats(id) on delete cascade,
    sender_type varchar(20) not null,
    sender_id bigint not null,
    body text not null,
    created_at timestamptz not null default now()
    );

create index if not exists idx_messages_chat_id on messages(chat_id);
create index if not exists idx_messages_created_at on messages(created_at);
create index if not exists idx_chats_ticket_id on chats(ticket_id);