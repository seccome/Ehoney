use mysql;

update user set authentication_string = password('Ehoney2021') where user = 'root';
 
GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY 'Ehoney2021' WITH GRANT OPTION;;
 
flush privileges;