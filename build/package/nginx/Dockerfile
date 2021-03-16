# Версия контейнера nginx на linux Alpine
# Источники:
# 	- https://github.com/nginxinc/docker-nginx/blob/25a3fc7343c6916fce1fba32caa1e8de8409d79f/stable/alpine/Dockerfile
# 	- http://firstwiki.ru/index.php/%D0%94%D0%BE%D0%B1%D0%B0%D0%B2%D0%BB%D0%B5%D0%BD%D0%B8%D0%B5_%D0%BC%D0%BE%D0%B4%D1%83%D0%BB%D0%B5%D0%B9_nginx_%D0%B2_Linux_(Debian/Ubuntu/CentOS)
# 	- https://github.com/nginxinc/docker-nginx/blob/de8822d8d91ff0802989bc0a12f8ab55596a513c/mainline/jessie/Dockerfile
# 	- https://github.com/yzprofile/ngx_http_dyups_module


# Базовый образ
FROM alpine:3.4

# Автор
MAINTAINER Andrey A. Kor. <andrey_kor@inbox.ru>

# Устанавливаемая версия nginx
ENV NGINX_VERSION 1.15.5

##########################################################################################################################################

RUN \
#	GPG_KEYS=B0F4253373F8F6F510D42178520A9993A1C052F8 \
# Конфигурация сборки nginx
	CONFIG="\
		--prefix=/etc/nginx \
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
		--with-http_ssl_module \
		--with-http_realip_module \
#		--with-http_addition_module \
#		--with-http_sub_module \
#		--with-http_dav_module \
#		--with-http_flv_module \
#		--with-http_mp4_module \
#		--with-http_gunzip_module \
		--with-http_gzip_static_module \
#		--with-http_random_index_module \
#		--with-http_secure_link_module \
		--with-http_stub_status_module \
#		--with-http_auth_request_module \
#		--with-http_xslt_module=dynamic \
#		--with-http_image_filter_module=dynamic \
#		--with-http_geoip_module=dynamic \
#		--with-http_perl_module=dynamic \
		--with-threads \
#		--with-stream \
#		--with-stream_ssl_module \
#		--with-http_slice_module \
#		--with-mail \
#		--with-mail_ssl_module \
		--with-file-aio \
#		--with-http_v2_module \
#		--with-ipv6 \
###		--with-cc-opt='-g -Wformat -Werror=format-security -Wp,-D_FORTIFY_SOURCE=2' \
###		--with-ld-opt='-Wl,-z,relro -Wl,-z,now -Wl,--as-needed' \
		--with-pcre \
####		--add-module=./ngx_http_dyups_module \
	" \
	&& addgroup -S nginx \
	&& adduser -D -S -h /var/cache/nginx -s /sbin/nologin -G nginx nginx \
	\
	\
# Ставим пакеты в ОС
	\
	&& apk add --no-cache --virtual .build-deps \
		gcc \
		libc-dev \
		make \
		openssl-dev \
		pcre-dev \
		zlib-dev \
		linux-headers \
		curl \
		gnupg \
#		libxslt-dev \
#		gd-dev \
#		geoip-dev \
#		perl-dev \
		git \
#		mc \
		tzdata \
	\
	\
# Указываем таймзону в ОС
	\
	&& cp /usr/share/zoneinfo/Europe/Moscow /etc/localtime \
	&& echo "Europe/Moscow" > /etc/timezone \
	\
	\
# Выкачиваем исходники nginx
	\
	&& curl -fSL http://nginx.org/download/nginx-$NGINX_VERSION.tar.gz -o nginx.tar.gz \
#	&& curl -fSL http://nginx.org/download/nginx-$NGINX_VERSION.tar.gz.asc -o nginx.tar.gz.asc \
	\
	&& export GNUPGHOME="$(mktemp -d)" \
	\
#	&& gpg --keyserver ha.pool.sks-keyservers.net --recv-keys "$GPG_KEYS" \
#	&& gpg --batch --verify nginx.tar.gz.asc nginx.tar.gz \
	\
#	&& rm -r "$GNUPGHOME" nginx.tar.gz.asc \
	&& mkdir -p /usr/src \
	&& tar -zxC /usr/src -f nginx.tar.gz \
	&& rm nginx.tar.gz \
	\
	&& cd /usr/src/nginx-$NGINX_VERSION \
	\
	\
# Выкачиваем модулья nginx'а для работы с upstream через RESTA API
	\
	&& git clone git://github.com/yzprofile/ngx_http_dyups_module.git \
	\
	&& ./configure $CONFIG --with-debug \
	&& make -j$(getconf _NPROCESSORS_ONLN) \
	\
	&& mv objs/nginx objs/nginx-debug \
#	&& mv objs/ngx_http_xslt_filter_module.so objs/ngx_http_xslt_filter_module-debug.so \
#	&& mv objs/ngx_http_image_filter_module.so objs/ngx_http_image_filter_module-debug.so \
#	&& mv objs/ngx_http_geoip_module.so objs/ngx_http_geoip_module-debug.so \
#	&& mv objs/ngx_http_perl_module.so objs/ngx_http_perl_module-debug.so \
	&& ./configure $CONFIG \
	\
	&& make -j$(getconf _NPROCESSORS_ONLN) \
	&& make install \
	\
	&& rm -rf /etc/nginx/html/ \
	\
	&& mkdir /etc/nginx/conf.d/ \
	&& mkdir -p /usr/share/nginx/html/ \
	\
	&& install -m644 html/index.html /usr/share/nginx/html/ \
	&& install -m644 html/50x.html /usr/share/nginx/html/ \
	&& install -m755 objs/nginx-debug /usr/sbin/nginx-debug \
#	&& install -m755 objs/ngx_http_xslt_filter_module-debug.so /usr/lib/nginx/modules/ngx_http_xslt_filter_module-debug.so \
#	&& install -m755 objs/ngx_http_image_filter_module-debug.so /usr/lib/nginx/modules/ngx_http_image_filter_module-debug.so \
#	&& install -m755 objs/ngx_http_geoip_module-debug.so /usr/lib/nginx/modules/ngx_http_geoip_module-debug.so \
#	&& install -m755 objs/ngx_http_perl_module-debug.so /usr/lib/nginx/modules/ngx_http_perl_module-debug.so \
	\
	&& ln -s ../../usr/lib/nginx/modules /etc/nginx/modules \
	\
	&& strip /usr/sbin/nginx* \
#	&& strip /usr/lib/nginx/modules/*.so \
	\
	&& rm -rf /usr/src/nginx-$NGINX_VERSION \
	\
	\
# Bring in gettext so we can get `envsubst`, then throw
# the rest away. To do this, we need to install `gettext`
# then move `envsubst` out of the way so `gettext` can
# be deleted completely, then move `envsubst` back.
	\
	&& apk add --no-cache --virtual .gettext gettext \
	&& mv /usr/bin/envsubst /tmp/ \
	\
	&& runDeps="$( \
		scanelf --needed --nobanner /usr/sbin/nginx /usr/lib/nginx/modules/*.so /tmp/envsubst \
			| awk '{ gsub(/,/, "\nso:", $2); print "so:" $2 }' \
			| sort -u \
			| xargs -r apk info --installed \
			| sort -u \
	)" \
	&& apk add --no-cache --virtual .nginx-rundeps $runDeps \
	\
	\
# Прибираем постустановочный мусор
	\
	&& apk del .build-deps \
	&& apk del .gettext \
	\
	&& mv /tmp/envsubst /usr/local/bin/ \
	\
	\
# forward request and error logs to docker log collector
	\
	&& ln -sf /dev/stdout /var/log/nginx/access.log \
	&& ln -sf /dev/stderr /var/log/nginx/error.log


##########################################################################################################################################


# Готовим директории
RUN \
	rm -rf /etc/nginx/conf.d && \
	mkdir -p /etc/nginx/conf.d && \
	chown -R root:root /etc/nginx/conf.d && \
	chmod 0777 /etc/nginx/conf.d && \
	\
	chown root:root /etc/nginx/nginx.conf && \
	chmod 0777 /etc/nginx/nginx.conf



##########################################################################################################################################


# Путь к конфиг-файлу
COPY ./configs/nginx/nginx.conf /etc/nginx/nginx.conf

# Путь к дополнительным файлам конфигурации
COPY ./configs/nginx/global /etc/nginx/global
COPY ./configs/nginx/local /etc/nginx/local

# Монтируем локальные конфигурации nginx, тем самым мы избавляем себя от необходимости, при изменении конфигов, заново пересобирать контейнеры
VOLUME /etc/nginx/local

# Путь к ssl-сертификатам
#COPY ./ssl /etc/ssl
#VOLUME /etc/ssl

# Монтируем директорию логов
#VOLUME /var/log

# Выводим наружу порт
EXPOSE 80 443


##########################################################################################################################################


# Docker: запуск нескольких приложений в одном контейнере
# Источники:
#	- https://blog.amartynov.ru/docker-%D0%B7%D0%B0%D0%BF%D1%83%D1%81%D0%BA-%D0%BD%D0%B5%D1%81%D0%BA%D0%BE%D0%BB%D1%8C%D0%BA%D0%B8%D1%85-%D0%BF%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B9/
COPY ./build/package/nginx/start.sh /
RUN chmod +x /start.sh
CMD ["/start.sh"]

#CMD ["nginx", "-g", "daemon off;"]
