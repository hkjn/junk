FROM ubuntu:latest

MAINTAINER Henrik Jonsson <me@hkjn.me>

# Install latest updates.
RUN apt-get update

# Install mysql client and server.
RUN apt-get -y install mysql-client mysql-server curl

# Enable remote access (default is localhost only, we change this
# otherwise our database would not be reachable from outside the container)
RUN sed -i -e"s/^bind-address\s*=\s*127.0.0.1/bind-address = 0.0.0.0/" /etc/mysql/my.cnf

# Create database, with some dummy data.
ADD ./create_db.sql /var/db/create_db.sql

# Uncomment to set DB credentials. Note: Not secure.
# This can also be specified by passing --env DB_USER=foo --env DB_PASSWORD=bar to "docker run".
# ENV DB_USER dbuser
# ENV DB_PASSWORD dbsecret
ENV SQL_URL file:/var/db/create_db.sql
ENV DB_ACCESS WRITE

# Install starting script.
ADD ./start_db.sh /usr/local/bin/start_db.sh
RUN chmod +x /usr/local/bin/start_db.sh

EXPOSE 3306

CMD ["/usr/local/bin/start_db.sh"]
