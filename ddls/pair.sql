-- auto-generated definition
create table pair
(
    address    varchar(255)                default ''::character varying not null,
    token0     varchar(255)                default ''::character varying not null,
    token1     varchar(255)                default ''::character varying not null,
    chain_id   integer                     default 0                     not null,
    reserve0   numeric(78)                 default 0                     not null,
    reserve1   numeric(78)                 default 0                     not null,
    block      bigint                      default 0                     not null,
    block_at   timestamp(6) with time zone default NULL::timestamp with time zone,
    created_at timestamp(6) with time zone                               not null,
    id         uuid                        default uuid_generate_v4()    not null
        primary key,
    name       varchar(64)                 default ''::character varying not null
);

alter table pair
    owner to postgres;

create unique index idx_pair_chain_addr
    on pair (chain_id, address);

