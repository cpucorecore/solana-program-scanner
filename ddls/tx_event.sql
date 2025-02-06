-- auto-generated definition
create type tx_event as enum ('add', 'remove', 'buy', 'sell', 'burn');

alter type tx_event owner to postgres;

