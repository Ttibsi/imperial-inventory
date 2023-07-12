FROM mysql

# This is fine for a simple thing like this, but should probably be passed from
# a secrets engine like Doppler in production
ENV MYSQL_DATABASE="testdata" \
    MYSQL_ROOT_PASSWORD="rootpass" \
    MYSQL_USER="user" \
    MYSQL_PASSWORD="password"

# On initialisation, mysql can run any sql files alphabetically in this folder
COPY schema.sql /docker-entrypoint-initdb.d/
EXPOSE 3306
