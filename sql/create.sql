create table profiles
(
	uid text primary key ,
	name text not null,
	coins int not null
);

create table tasks
(
	id serial not null
		constraint tasks_pk
			primary key,
	created_by text not null,
	date_to_complete date,
	task_type text,
	time_to_complete time,
	lat real,
	long real,
	reward int not null,
    description text,
    status int
);

create table tasks_accepted
(
	uid text
		constraint tasks_accepted_profiles_uid_fk
			references profiles,
	task_id serial
		constraint tasks_accepted_tasks_id_fk
			references tasks,
	status int not null,
	constraint tasks_accepted_pk
		primary key (uid, task_id)
);
