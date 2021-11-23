FROM golang:1.17.3
RUN apt-get update && apt-get install -y apache2 vim curl && a2enmod rewrite lbmethod_byrequests proxy proxy_http proxy_balancer
RUN mkdir -p /usr/go-apache-test/app /opt
COPY .docker/dev/entrypoint.sh /opt
RUN chmod +x /opt/entrypoint.sh
WORKDIR /usr/go-apache-test/app
CMD ["sh", "/opt/entrypoint.sh"]