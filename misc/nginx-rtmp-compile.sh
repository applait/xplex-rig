#!/bin/sh

set -ex

NGINX_VERSION=1.12.0
NGINX_RTMP_MODULE_VERSION=1.1.11
NGINX_RTMP_MODULE_PATH="../nginx-rtmp-module-${NGINX_RTMP_MODULE_VERSION}"
TMP_DIR="data"

GPG_KEYS="B0F4253373F8F6F510D42178520A9993A1C052F8"
CONFIG="\
		--prefix=\"\" \
		--sbin-path=/usr/sbin/nginx \
		--modules-path=/usr/lib/nginx/modules \
		--conf-path=/etc/nginx/nginx.conf \
		--error-log-path=/var/log/nginx/error.log \
		--http-log-path=/var/log/nginx/access.log \
		--pid-path=/var/run/nginx.pid \
		--lock-path=/var/run/nginx.lock \
		--http-client-body-temp-path=/var/cache/nginx/client_temp \
		--http-proxy-temp-path=/var/cache/nginx/proxy_temp \
		--http-fastcgi-temp-path=/var/cache/nginx/fastcgi_temp \
		--http-uwsgi-temp-path=/var/cache/nginx/uwsgi_temp \
		--http-scgi-temp-path=/var/cache/nginx/scgi_temp \
		--user=nginx \
		--group=nginx \
        --add-module=${NGINX_RTMP_MODULE_PATH} \
		--with-http_ssl_module \
		--with-http_stub_status_module \
		--with-http_v2_module \
		--with-ipv6"

mkdir -p $TMP_DIR
mkdir -p /usr/lib/nginx/modules
mkdir -p /var/cache/nginx
mkdir -p /var/log/nginx
rm -rf $TMP_DIR/*
pushd $TMP_DIR

dnf builddep nginx --assumeyes
dnf install gnupg gcc automake

curl -fSL http://nginx.org/download/nginx-$NGINX_VERSION.tar.gz -o nginx.tar.gz
curl -fSL http://nginx.org/download/nginx-$NGINX_VERSION.tar.gz.asc  -o nginx.tar.gz.asc
curl -fSL https://github.com/arut/nginx-rtmp-module/archive/v${NGINX_RTMP_MODULE_VERSION}.tar.gz -o nginx-rtmp-module.tar.gz

gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$GPG_KEYS"
gpg --batch --verify nginx.tar.gz.asc nginx.tar.gz
rm -r nginx.tar.gz.asc

tar -zxf nginx.tar.gz
tar -zxf nginx-rtmp-module.tar.gz

pushd nginx-$NGINX_VERSION
./configure $CONFIG
make -j$(getconf _NPROCESSORS_ONLN)
popd
popd
cp -v $TMP_DIR/nginx-$NGINX_VERSION/objs/nginx $TMP_DIR/

# rm -rf /etc/nginx/html/
# mkdir /etc/nginx/conf.d/
# ln -s ../../usr/lib/nginx/modules /etc/nginx/modules
# strip /usr/sbin/nginx
# rm -rf /usr/src/nginx-$NGINX_VERSION
