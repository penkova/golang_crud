# golang-restful, Simple CRUD With Golang Mysql
Simple Golang REST API with MySQL


# Install
```
go get -u github.com/gorilla/mux
```
```
go get -u github.com/go-sql-driver/mysql
```
### Creating database:
```
mysql> create database cars;
mysql> use cars;
mysql> create table cars(id int(11) not null primary key auto_increment, model varchar(20),age int(11), price int(20));

mysql> insert into cars(model,age,price) values('BMW', 3, 200000);
mysql> insert into cars(model,age,price) values('Audi', 1, 99999);
mysql> insert into cars(model,age,price) values('KIA', 5, 40000);
```
