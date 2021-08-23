```bash
# start a mysql
$ docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -d -p 3306:3306 mysql:5.7
# create database
$ docker exec -it mysql mysql -uroot -p
$ create database orm_test;
$ exit

# run demo
$ go run .
```