FROM nginx:stable-alpine
RUN mkdir -p /www/data/devnet_genesis
COPY Docker/proxy/nginx.conf /etc/nginx/conf.d
COPY Docker/assets/node1/config/genesis.json /www/data/devnet_genesis/genesis.json
EXPOSE 1318/tcp
CMD ["/bin/sh", "-c", "exec nginx -g 'daemon off;';"]
WORKDIR /usr/share/nginx/html
