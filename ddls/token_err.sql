-- auto-generated definition
create table token_err
(
    address    varchar(80)                 not null
        constraint token_get_failed_pkey
            primary key,
    reason     integer default 0           not null,
    created_at timestamp(6) with time zone not null
);

alter table token_err
    owner to postgres;

