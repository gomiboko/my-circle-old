FROM mysql:8.0.25
COPY init/* /docker-entrypoint-initdb.d/
COPY ./my.cnf /etc/mysql/conf.d/my.cnf
CMD [ "mysqld" ]
