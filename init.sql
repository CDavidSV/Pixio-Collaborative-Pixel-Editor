create table users (
    user_id char(26) primary key,
    username varchar(32) not null,
    avatar_url text,
    email varchar(255) unique not null,
    hashed_password char(60) not null,
    created_at timestamptz default now()
);

create table canvases (
    canvas_id char(26) primary key,
    owner_id char(26) not null,
    title varchar(32) not null,
    description varchar(512) default "",
    width int not null,
    height int not null,
    data bytea not null,
    last_edited_at timestamptz default now(),
    access_type int not null,
    created_at timestamptz default now(),
    star_count int not null,

    foreign key (owner_id) references users(user_id)
);

create table versions (
    version_id char(26) primary key,
    canvas_id char(26) not null,
    data bytea not null,
    timestamp timestamptz not null default now(),
    edit_count int not null,

    foreign key (canvas_id) references canvases(canvas_id)
);

create table edits (
    version_id char(26) not null,
    user_id char(26) not null,
    data bytea not null,
    timestamp timestamptz not null default now(),

    primary key (version_id, user_id),
    foreign key (user_id) references users(user_id)
);

create type object_type as enum ('canvas', 'collection');

create table access_rules (
    object_id char(26) not null,
    object_type object_type,
    user_id char(26) not null,
    permissions int not null,
    last_modified_at timestamptz default now(),
    last_modified_by char(26),

    primary key (object_id, user_id),
    foreign key (user_id) references users(user_id),
    foreign key (last_modified_by) references users(user_id)
);

create table stars (
    canvas_id char(26) not null,
    user_id char(26) not null,
    added_at timestamptz default now(),

    primary key (canvas_id, user_id),
    foreign key (user_id) references users(user_id),
    foreign key (canvas_id) references canvases(canvas_id)
);

create table collections (
    collection_id char(26) primary key,
    owner_id char(26) not null,
    title varchar(32) not null,
    description varchar(512),
    access_type int not null,
    saves_count int not null,

    foreign key (owner_id) references users(user_id)
);

create table collection_canvas (
    collection_id char(26) not null,
    canvas_id char(26) not null,
    added_by char(26) not null,
    added_at timestamptz not null default now(),

    primary key (collection_id, canvas_id),
    foreign key (canvas_id) references canvases(canvas_id),
    foreign key (collection_id) references collections(collection_id)
);

create table saved_collections (
    user_id char(26) not null,
    collection_id char(26) not null,
    added_at timestamptz default now(),

    primary key (user_id, collection_id),
    foreign key (collection_id) references collections(collection_id),
    foreign key (user_id) references users(user_id)
);

create table user_sessions (
    session_id char(26) primary key,
    user_id char(26) not null,
    expires_at timestamptz not null,
    refresh_token text not null,
    created_at timestamptz default now(),
    last_accessed timestamptz default now(),

    foreign key(user_id) references users(user_id)
);
