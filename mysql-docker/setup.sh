#!/bin/bash

echo 'checking mysql status.'
service mysql status

echo '1.start mysql....'
service mysql start
sleep 3
service mysql status

db_path=/var/lib/mysql/sec_ehoneypot
if [ ! -d "${db_path}" ]; then
    echo '2.start importing data....'
	mysql < /mysql/schema.sql
	echo '3.end importing data....'
	sleep 3
	service mysql status
	echo '4.start changing password....'
	mysql < /mysql/privileges.sql
	echo '5.end changing password....'
fi

sleep 3
service mysql status
echo 'mysql is ready'

tail -f /dev/null