FROM nginx
RUN echo "application/wasm wasm" >> /etc/mime.types

COPY site /usr/share/nginx/html
