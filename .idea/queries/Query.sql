-- departments
create table departments (
                             id serial primary key,
                             name text not null,
                             created_at timestamp default now(),
                             updated_at timestamp default now(),
                             deleted_at timestamp
);

-- users
create table users (
                       id serial primary key,
                       full_name text not null,
                       email text not null unique,
                       password text not null,
                       department_id int references departments(id) on delete set null,
                       created_at timestamp default now(),
                       updated_at timestamp default now(),
                       deleted_at timestamp
);

-- roles
create table roles (
                       id serial primary key,
                       name text not null,
                       department_id int references departments(id) on delete set null
);

-- users_roles
create table users_roles (
                             id serial primary key,
                             user_id int references users(id) on delete cascade,
                             role_id int references roles(id) on delete cascade
);

-- task status enum
create type task_status as enum ('todo', 'in_progress', 'done');

-- tasks
create table tasks (
                       id serial primary key,
                       title text not null,
                       description text not null,
                       deadline timestamp,
                       department_id int references departments(id) on delete set null,
                       creator_id int references users(id) on delete set null,
                       assignee_id int references users(id) on delete set null,
                       status task_status not null default 'todo',
                       created_at timestamp default now(),
                       updated_at timestamp default now(),
                       deleted_at timestamp
);

-- task comments
create table tasks_comments (
                                id serial primary key,
                                comment text not null,
                                created_at timestamp default now(),
                                updated_at timestamp default now(),
                                deleted_at timestamp,
                                task_id int references tasks(id) on delete cascade,
                                user_id int references users(id) on delete set null
);

-- attendance status enum
create type attendance_status as enum ('present', 'absent', 'excused');

-- attendance
create table attendance (
                            id serial primary key,
                            user_id int references users(id) on delete cascade,
                            department_id int references departments(id) on delete set null,
                            status attendance_status,
                            comment text,
                            marked_by int references users(id) on delete set null,
                            created_at timestamp default now(),
                            updated_at timestamp default now(),
                            deleted_at timestamp
);
