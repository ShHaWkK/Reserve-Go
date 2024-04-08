/*
* File : BDD.sql
*/

CREATE TABLE rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,
    available BOOLEAN DEFAULT TRUE
);

CREATE TABLE reservations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    room_id INT,
    date DATE NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULLs
);


INSERT INTO rooms (name, capacity) VALUES ('Salle A', 40);
INSERT INTO rooms (name, capacity) VALUES ('Salle B', 30);
INSERT INTO rooms (name, capacity) VALUES ('Salle C', 50);