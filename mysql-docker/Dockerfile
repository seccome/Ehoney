
FROM mysql:5.6
 
ENV MYSQL_ALLOW_EMPTY_PASSWORD yes
 
COPY setup.sh /mysql/setup.sh
COPY schema.sql /mysql/schema.sql
COPY privileges.sql /mysql/privileges.sql
 
CMD ["sh", "/mysql/setup.sh"]