DROP DATABASE IF EXIST customer;
CREATE DATABASE customer;
USE customer;

CREATE TABLE customers(id int NOT NULL AUTO_INCREMENT, name varchar(20),PRIMARY KEY(id));

INSERT INTO customers VALUES(1,'abc');
INSERT INTO customers VALUES(2,'pqr');