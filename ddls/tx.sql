-- auto-generated definition
create table tx
(
    tx_hash        varchar(128)                        not null,
    event          tx_event default 'buy'::tx_event    not null,
    token0_amount  numeric(70, 18)                     not null,
    token1_amount  numeric(70, 18)                     not null,
    maker          varchar(64)                         not null,
    token0_address varchar(64)                         not null,
    token1_address varchar(64)                         not null,
    amount_usd     numeric(70, 18)                     not null,
    price_usd      numeric(70, 18)                     not null,
    block          bigint                              not null,
    block_at       timestamp(6) with time zone         not null,
    created_at     timestamp(6) with time zone         not null,
    block_index    integer                             not null,
    tx_index       integer,
    id             uuid     default uuid_generate_v4() not null
        primary key,
    pair_id        uuid
);

alter table tx
    owner to postgres;

create unique index tx_tx_hash_tx_index_block_index_uindex
    on tx (tx_hash, tx_index, block_index);

