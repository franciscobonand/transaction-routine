create schema if not exists pismo;

create table if not exists pismo.account (
    id serial primary key,
    document_number varchar(255) not null
);

create table if not exists pismo.operation_type (
    id serial primary key,
    description varchar(255) not null,
    positive_amount boolean not null
);

insert into pismo.operation_type (description, positive_amount) values ('COMPRA A VISTA', false);
insert into pismo.operation_type (description, positive_amount) values ('COMPRA PARCELADA', false);
insert into pismo.operation_type (description, positive_amount) values ('SAQUE', false);
insert into pismo.operation_type (description, positive_amount) values ('PAGAMENTO', true);

create table if not exists pismo.transaction (
    id serial primary key,
    account_id integer not null,
    operation_type_id integer not null,
    amount numeric not null,
    event_date timestamp not null,
    foreign key (account_id) references pismo.account(id),
    foreign key (operation_type_id) references pismo.operation_type(id)
);