#!/bin/sh
rm -rf /etc/apache2/sites-enabled/*
rm -rf /etc/apache2/sites-available/*
make init
apachectl configtest
/etc/init.d/apache2 start
make run
/etc/init.d/apache2 stop