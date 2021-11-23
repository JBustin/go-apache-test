#!/bin/sh
rm -rf /etc/apache2/sites-enabled/*
rm -rf /etc/apache2/sites-available/*
go build
./go-apache-test -x init
apachectl configtest
/etc/init.d/apache2 start
./go-apache-test -x run
/etc/init.d/apache2 stop