FROM nginx:stable-alpine
COPY Docker/proxy/nginx.conf /etc/nginx/conf.d
EXPOSE 1318/tcp
CMD ["/bin/sh", "-c", "exec nginx -g 'daemon off;';"]
WORKDIR /usr/share/nginx/html