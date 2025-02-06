-- auto-generated definition
create table token
(
    address      varchar(64)                 default ''::character varying not null,
    name         varchar(255)                default ''::character varying not null,
    decimal      smallint                    default 0                     not null,
    chain_id     smallint                    default 0                     not null,
    total_supply varchar(255)                default ''::character varying not null,
    symbol       varchar(255)                default ''::character varying not null,
    block        bigint                      default 0                     not null,
    block_at     timestamp(6) with time zone default NULL::timestamp with time zone,
    created_at   timestamp(6) with time zone                               not null,
    logo         varchar(255)                default ''::character varying not null,
    creator      varchar(255)                default ''::character varying not null,
    holder_count bigint                      default 0                     not null,
    id           uuid                        default uuid_generate_v4()    not null
        primary key
);

alter table token owner to postgres;

create index id_token_symbol
    on token (symbol);

create unique index idx_token_chain_addr
    on token (chain_id, address);

create index idx_token_name
    on token (name);

