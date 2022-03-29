

# mysql

docker pull mysql

## run docker:
### with specific network:
docker run -d --name mysql --network urlshortner-network -p 3306:3306 -v ~/workspace/data/mysql-docker-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=admin123 mysql:latest

### without specific network
docker run -d --name mysql -p 3306:3306 -v ~/workspace/data/mysql-docker-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=admin123 mysql:latest

## view mysql logs:
docker logs mysql

## connect to mysql docker and start client console:

docker exec -it mysql bash

then in container prompt:
mysql -u root -p

then create a new databse :
CREATE DATABASE MYSQLTEST;


## connect to mysql as client through mysql container (require run mysql docker with speceific network)
following command starts another mysql container instance and runs the mysql command line client against your original mysql container

docker run -it --network urlshortner-network --rm mysql mysql -hmysql -uroot -p

# database/sql package

## driver 
golang database drivers list: https://github.com/golang/go/wiki/SQLDrivers

we use following driver:
github.com/go-sql-driver/mysql

# get driver:



# mysql cli 

## login to mysql with user root. (password will be asked in prompt):
mysql -u root -p

### create database
create database testDB;

### drop database
drop database testDB;

### databse list:
show databases;

### set current database (test is a database). by running this command next commands will apply to test database
use testDB;

### list tables (in currebt data base)
show tables;

### create table
create table psGroup ( 
    id INT NOT NULL, 
    name varchar(100) NOT NULL, 
    PRIMARY KEY (`id`) 
);

### show table schema
describe psgroup;

### insert into table
insert into psgroup  
    (id, name) 
values 
    (1, 'brand_1'), 
    (2, 'brand_2');

### select table rows
select * from psgruop;

## run file script

### create a file (ps_create_script.sql) and copy following script to it:

DROP DATABASE IF EXISTS powerstation;
create database powerstation;

use powerstation;

DROP TABLE IF EXISTS psGroup;
DROP TABLE IF EXISTS psRegistery;
create table psGroup ( 
    id INT NOT NULL, 
    name varchar(100) NOT NULL, 
    PRIMARY KEY (`id`) 
);

create table psRegistery ( 
    id INT AUTO_INCREMENT NOT NULL, 
    stGroup INT NOT NULL,
    stId varchar(256) NOT NULL,
    slotCount INT NOT NULL,
    PRIMARY KEY (`id`), 
    FOREIGN KEY (stGroup) REFERENCES psGroup(id)
);

insert into psgroup  
    (id, name) 
values 
    (1, 'brand_1'), 
    (2, 'brand_2');

insert into psRegistery  
    (stGroup, stId, slotCount) 
values 
    (1, 'a1a1a1a1a1a1a1a1', 8), 
    (1, 'b1b1b1b1b1b1b1b1', 8),
    (2, 'a2a2a2a2a2a2a2a2', 10),
    (2, 'b2b2b2b2b2b2b2b2', 10);



## then use foloowing commqnd to run it (path is example)
source \root\scripts\ps_create_script.sql;

## copy from local machine to docker container
copy file ps_create_script.sql from ~/myfiles folder to container with id: 135950565ad8 folder /root/scripts
sudo docker cp ~/myfiles/ps_create_script.sql 135950565ad8:/root/scripts