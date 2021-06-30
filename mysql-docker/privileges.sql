use mysql;

update user set authentication_string = password('Security#123456#') where user = 'root';
 
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'Security#123456#' WITH GRANT OPTION;;
 
flush privileges;