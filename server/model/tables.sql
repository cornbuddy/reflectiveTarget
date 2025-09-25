create table if not exists users (
    id serial primary key,
    name varchar(100)
);

create table if not exists targets (
    id serial primary key,
    name varchar(100),
    owner_id int references users(id)
);

create table if not exists questions (
    id serial primary key,
    text varchar(200),
    target_id int references targets(id)
);

create table if not exists shots (
    id serial primary key,
    x int,
    y int,
    target_id int references targets(id),
    shooter_id int references users(id)
);
