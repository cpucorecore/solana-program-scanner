-- auto-generated definition
create table block_err
(
    slot       bigint default 0 not null,
    err_code   integer          not null,
    created_at timestamp(6)     not null,
    constraint slot_status_pk
        primary key (slot, err_code)
);

alter table block_err
    owner to postgres;