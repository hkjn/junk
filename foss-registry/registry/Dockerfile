#
# Dockerfile for registry-nginx.
#
FROM hkjn/nginx

COPY registry.io.conf sites-available/
# TODO(hkjn): Experimental auth added with:
# docker run --entrypoint htpasswd httpd:2.4 -bn testuser testpassword
RUN echo "testuser:$apr1$vu8kCLUm$IIUVexAr0Ph6BcRv2NWfk." > conf.d/nginx.htpasswd
RUN ln -s sites-available/registry.io.conf /etc/nginx/sites-enabled/
