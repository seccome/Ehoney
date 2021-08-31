use mysql;

update user set authentication_string = password('123456') where user = 'root';
 
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '123456' WITH GRANT OPTION;;
 
flush privileges;