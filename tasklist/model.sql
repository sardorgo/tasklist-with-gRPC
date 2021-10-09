drop database tasklist;
create database tasklist;
\c tasklist;

create extension "uuid-ossp";

create table users (
    user_id uuid not null primary key,
    first_name varchar(32)
);

create table tasks (
    task_id uuid not null primary key,
    task_name varchar(32)
);

create table task_controller (
    task_controller_id uuid not null primary key,
    user_id uuid not null references users (user_id),
    task_id uuid not null references tasks (task_id)
);

select
    u.user_id,
    u.first_name,
    t.task_name

from users as u
join task_controller as tc on tc.user_id = u.user_id
join tasks as t on tc.task_id = t.task_id
;

insert into users (user_id, first_name) values (uuid_generate_v4(), 'Sarah');
insert into users (user_id, first_name) values (uuid_generate_v4(), 'John');

insert into tasks (task_id, task_name) values (uuid_generate_v4(), 'Create Database');
insert into tasks (task_id, task_name) values (uuid_generate_v4(), 'Create Server');

insert into task_controller(task_controller_id, user_id, task_id) values (uuid_generate_v4(), 'edf3b3b7-a59d-4664-9ffa-13c240f6b500', '5b9ba479-4ca4-4f11-b980-1b00e5a885da');
insert into task_controller(task_controller_id, user_id, task_id) values (uuid_generate_v4(), 'edf3b3b7-a59d-4664-9ffa-13c240f6b500', '66800637-6580-4544-a2c2-332b1539e3d6');


-------------------------------------------------


select
    u.user_id,
    u.first_name,
    array_agg(t.task_name)

from users as u
join task_controller as tc on tc.user_id = u.user_id
join tasks as t on tc.task_id = t.task_id
group by tc.task_id, u.user_id
;

insert into task_controller(task_controller_id, user_id, task_id) values (uuid_generate_v4(), '48e77112-9649-457a-823f-8c40c01f8bcd', '61bb09a2-d804-48f0-a9ef-4dbf7ba80707');

select task_id from task_controller
where user_id = '48e77112-9649-457a-823f-8c40c01f8bcd';



select 
    tc.user_id
from 
    task_controller as tc
where tc.task_id is not null;


