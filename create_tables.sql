DROP DATABASE IF EXISTS
GO;

CREATE DATABASE
GO;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS userdata;

\c GO;

CREATE TABLE users (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    username VARCHAR(100) PRIMARY KEY
);

CREATE TABLE userdata (
    useriD INTEGER NOT NULL,
    name VARCHAR(100),
    surname VARCHAR(100),
    description VARCHAR(200)
);