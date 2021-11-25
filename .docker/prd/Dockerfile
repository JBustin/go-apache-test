FROM golang:1.17.3
RUN apt-get update && apt-get install -y apache2 vim curl && a2enmod rewrite lbmethod_byrequests proxy proxy_http proxy_balancer
RUN mkdir -p /usr/app /opt
COPY .docker/prd/entrypoint.sh /opt
RUN chmod +x /opt/entrypoint.sh
WORKDIR /usr/go-apache-test/app
COPY . .
RUN go build
CMD ["sh", "/opt/entrypoint.sh"]