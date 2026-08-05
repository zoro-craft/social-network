CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,

    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,

    date_of_birth DATE NOT NULL,

    nickname TEXT UNIQUE,
    about_me TEXT,
    avatar TEXT,

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);