#!/bin/bash
#
# Starts the mysqld. Runs inside the Docker container.
#
echo "Creating user ${DB_USER} for databases loaded from ${SQL_URL}"

# Import database provided via 'docker run --env
# url="http://[host]/db.sql"' or via ENV variable in Dockerfile.
echo "Starting MySQL.."
/usr/sbin/mysqld &
# TODO: Do something nicer to wait for DB to come up than just sleeping.
sleep 5
echo "Creating DB from ${SQL_URL}.."
curl ${SQL_URL} --silent --max-time 1 | mysql --default-character-set=utf8
mysqladmin shutdown
echo "done."

/usr/sbin/mysqld &
sleep 5
echo "Creating user '${DB_USER}'.."
echo "CREATE USER '${DB_USER}' IDENTIFIED BY '${DB_PASSWORD}'" | mysql --default-character-set=utf8
echo "REVOKE ALL PRIVILEGES ON *.* FROM '${DB_USER}'@'%'; FLUSH PRIVILEGES" | mysql --default-character-set=utf8
echo "GRANT SELECT ON *.* TO '${DB_USER}'@'%'; FLUSH PRIVILEGES" | mysql --default-character-set=utf8
echo "done."

if [ "${DB_ACCESS}" = "WRITE" ]; then
		echo "Adding write access for '${DB_USER}'"
		echo "GRANT ALL PRIVILEGES ON *.* TO '${DB_USER}'@'%' WITH GRANT OPTION; FLUSH PRIVILEGES" | mysql --default-character-set=utf8
fi

echo "Restarting MySQL.."
mysqladmin shutdown
/usr/sbin/mysqld
