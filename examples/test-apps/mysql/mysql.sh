#!/bin/bash
while true
do
    mysql -u $MYSQL_DBUSER_ENV -h $MYSQL_SERVER_ENV --password=$MYSQL_DBPASS_ENV
    if [ "$?" -eq 0 ]; then
        echo "Connection to remote mysql server was successful"
    else 
        echo "Connection to remote mysql server failed"
    fi
    sleep 5
done
