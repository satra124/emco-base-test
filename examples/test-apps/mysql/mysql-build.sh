#!/bin/bash
docker build --no-cache -t mysql-client --build-arg MYSQL_DBUSER=$1 --build-arg MYSQL_SERVER=$2 --build-arg MYSQL_DBPASS=$3 -f Dockerfile_mysql .
docker tag mysql-client:latest $DOCKER_REG/mysql-client:1.0
docker push $DOCKER_REG/mysql-client:1.0
