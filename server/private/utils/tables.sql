create table if not exists users (
    id serial primary key,
    username varchar(64) not null,
    hashed_password varchar(64) not null,
    salt varchar(64) not null
);

create table if not exists targets (
    id serial primary key,
    name varchar(128) not null,
    owner_id int references users(id) not null
);

create table if not exists questions (
    id serial primary key not null,
    text varchar(128) not null,
    target_id int references targets(id)
);

create table if not exists shots (
    id serial primary key,
    x int not null,
    y int not null,
    target_id int references targets(id) not null,
    shooter_id int references users(id) not null
);
