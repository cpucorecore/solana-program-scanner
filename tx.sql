create table tx
(
    tx_hash        varchar(88)                 not null,
    event          smallint default 0          not null,
    token0_amount  varchar(70)                 not null,
    token1_amount  varchar(70)                 not null,
    maker          varchar(64)                 not null,
    token0_address varchar(64)                 not null,
    token1_address varchar(64)                 not null,
    amount_usd     numeric(70, 18)             not null,
    price_usd      numeric(70, 18)             not null,
    block          bigint                      not null,
    block_at       timestamp(6) with time zone not null,
    created_at     timestamp(6) with time zone not null,
    index          integer                     not null
);

alter table tx
    owner to postgres;
