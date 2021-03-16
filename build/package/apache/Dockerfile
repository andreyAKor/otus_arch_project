# Базовый образ
FROM php:7.4-apache

# Автор
MAINTAINER Andrey A. Kor. <andrey_kor@inbox.ru>


##########################################################################################################################################


# Подготовка системы
RUN \
	apt update && \
	apt upgrade -y && \
	ulimit -s 524288

# Ставим утилиты
RUN \
	apt install curl -y && \
	apt install wget -y && \
	apt install cmake -y && \
	apt install g++ -y && \
	apt install make -y && \
	apt install openssl -y && \
	apt install libssl-dev -y

# Указываем таймзону в ОС
RUN \
	echo "Europe/Moscow" > /etc/timezone && \
	dpkg-reconfigure --frontend noninteractive tzdata

# Готовим директории
RUN \
	rm -rf /var/www/html && \
	mkdir -p /var/www/html && \
	chown -R www-data:www-data /var/www/html && \
	chmod 0777 /var/www/html && \
	\
	rm -rf /usr/local/etc/php && \
	mkdir -p /usr/local/etc/php && \
	chown -R root:staff /usr/local/etc/php && \
	chmod 0777 /usr/local/etc/php && \
	\
	rm -rf /etc/apache2 && \
	mkdir -p /etc/apache2 && \
	chown -R root:root /etc/apache2 && \
	chmod 0777 /etc/apache2

# Ставим SSH утилиты
RUN \
	apt install ssh -y


##########################################################################################################################################


# Увеличиваем размер стека для процессов:
# Источник:
# 	- https://habrahabr.ru/company/bitrix/blog/153001/
RUN \
	echo "* hard stack 524288" >> /etc/security/limits.conf && \
	echo "* soft stack 524288" >> /etc/security/limits.conf && \
	echo "root hard stack 524288" >> /etc/security/limits.conf && \
	echo "root soft stack 524288" >> /etc/security/limits.conf


##########################################################################################################################################


# Создание SSH-ключей и привязка их к ББ
# Стоит отметить, что SSH-ключи берутся у root'а из хост системы, поэтому контейнеры должны быть запущены от имени root'а
# Источники:
#	- https://toster.ru/q/250345
### {{{

# Создаём ssh директорию
RUN mkdir /root/.ssh/

# Копируем приватный ключ и устанавливаем права доступа к нему
# Приватный ключ взят с площадки dev5.poezd.ru
ADD ./.ssh/id_rsa /root/.ssh/id_rsa

RUN \
	chmod 0700 /root/.ssh/id_rsa && \
	echo "Host youkey"$'\n'"HostName bitbucket.org"$'\n'"IdentityFile /root/.ssh/id_rsa" >> /root/.ssh/config

# Копируем публичный ключ
# Публичный ключ взят с площадки dev5.poezd.ru
ADD ./.ssh/id_rsa.pub /root/.ssh/id_rsa.pub

# Действия:
# - Создаём known_hosts
# - Добавляем ключ bitbuckets в known_hosts
# - Простановка прав доступа к папкам и файлам
RUN \
	touch /root/.ssh/known_hosts && \
	ssh-keyscan bitbucket.org >> /root/.ssh/known_hosts && \
	chmod 0700 /root/.ssh/ && \
	chmod 0600 /root/.ssh/id_rsa && \
	chmod 0644 /root/.ssh/id_rsa.pub && \
	chmod 0644 /root/.ssh/known_hosts
### }}}


##########################################################################################################################################


# Ставим hg
# Источники:
# 	- http://poleshuk.ru/nastrojka-i-rabota-s-mercurial-v-ubuntu/
# 	- http://askubuntu.com/questions/482928/ignore-apt-get-postinstall-scripts-automatically
RUN \
	apt-get download javascript-common && \
	dpkg --unpack javascript-common*.deb && \
	rm /var/lib/dpkg/info/javascript-common.postinst -f && \
	dpkg --configure javascript-common && \
	apt install -yf && \
	apt install mercurial -y

# Ставим git
# Источники:
# 	- https://toster.ru/q/250345
RUN \
	apt install -y git


##########################################################################################################################################


### composer {{{
# Ставим composer
# Источники:
# 	- http://ermin.ru/post/11
RUN \
	cd /usr/local/bin && \
	curl -sS https://getcomposer.org/installer | php && \
	chmod 0755 composer.phar && \
	ln -s composer.phar composer

# Либы для composer
RUN \
	composer global require "fxp/composer-asset-plugin:v1.1.4" && \
	composer global require "hirak/prestissimo:0.3.4"
### }}}


##########################################################################################################################################


# Установка wkhtmltopdf + xvfb-run для его запуска на Debian
# Источники:
# 	- https://www.howtoinstall.co/en/debian/jessie/wkhtmltopdf
#	- http://stackoverflow.com/questions/37100442/wkhtmltopdf-updated-my-ubuntu-server-and-got-qxcbconnection-could-not-connect
### {{{
RUN \
	apt install xvfb -y && \
	apt install wkhtmltopdf -y && \
	chmod 0777 -R /tmp
### }}}


##########################################################################################################################################


### php-расширения {{{

# Подготовка
RUN \
	mkdir /usr/local/etc/php/conf.d/

# Ставим:
# 	- либу rabbitmq
# 	- php-расширение AMQP
RUN \
	wget https://github.com/alanxz/rabbitmq-c/archive/v0.9.0.tar.gz && \
	tar -zxvf v0.9.0.tar.gz && \
	cd rabbitmq-c-0.9.0/ && \
	\
	mkdir build && cd build && \
	cmake -DOPENSSL_ROOT_DIR=/usr/local/ssl -DOPENSSL_LIBRARIES=/usr/local/ssl/lib -DCMAKE_INSTALL_PREFIX=/usr/local .. && \
	\
	cmake --build . --target install && \
	\
	cd ../.. && \
	rm -r rabbitmq-c-0.9.0/ && \
	rm v0.9.0.tar.gz && \
	\
	pecl install amqp && \
	docker-php-ext-enable amqp

# Ставим:
# 	- php-расширения: bcmath, calendar, gettext, mysql, pdo_mysql, shmop, sockets, sysvmsg, sysvsem, sysvshm, pcntl
RUN \
	docker-php-ext-install -j$(nproc) bcmath calendar gettext mysql pdo_mysql shmop sockets sysvmsg sysvsem sysvshm pcntl

# Ставим:
# 	- либы libfreetype6 libjpeg62-turbo libmcrypt libpng
# 	- php-расширения: gd
RUN \
	apt install -y libfreetype6-dev libjpeg62-turbo-dev libmcrypt-dev libpng-dev && \
	docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ && \
	docker-php-ext-install -j$(nproc) gd

# Ставим:
# 	- либу libbz2
# 	- php-расширения: bz2
RUN \
	apt install -y libbz2-dev && \
	docker-php-ext-install -j$(nproc) bz2

# Ставим:
# 	- либу libevent
# 	- php-расширения: event
# Источники:
# 	- https://bitbucket.org/osmanov/pecl-event
RUN \
	wget https://github.com/libevent/libevent/releases/download/release-2.1.8-stable/libevent-2.1.8-stable.tar.gz && \
	tar -zxvf libevent-2.1.8-stable.tar.gz && \
	cd libevent-2.1.8-stable/ && \
	./configure && \
	make && make install && \
	cd .. && \
	rm -r libevent-2.1.8-stable/ && \
	rm libevent-2.1.8-stable.tar.gz && \
	\
	git clone git@bitbucket.org:osmanov/pecl-event.git && \
	cd ./pecl-event && \
	phpize && \
	./configure --with-php-config=/usr/local/bin/php-config --enable-event-debug=no --enable-event-sockets=yes --with-event-libevent-dir=/usr --with-event-pthreads=no --with-event-extra --with-event-openssl=no --with-openssl-dir=no && \
	make && make install && \
	cd .. && \
	rm -r ./pecl-event && \
	docker-php-ext-enable event

# Ставим:
# 	- либу libmcrypt
# 	- php-расширения: mcrypt
RUN \
	apt install -y libmcrypt-dev && \
	docker-php-ext-install -j$(nproc) mcrypt

# Ставим:
# 	- либу libmhash
# 	- php-расширения: mhash
# Млин, не смог поставить эту хрень :)
#RUN \
#	curl -fsSL 'http://netix.dl.sourceforge.net/project/mhash/mhash/0.9.9.9/mhash-0.9.9.9.tar.gz' -o mhash.tgz && \
#	mkdir -p mhash && \
#	tar -xf mhash.tgz -C mhash --strip-components=1 && \
#	rm mhash.tgz && \
#	( \
#		cd mhash && \
#		phpize && \
#		./configure --enable-mhash && \
#		make -j$(nproc) && \
#		make install \
#	) && \
#	rm -r mhash && \
#	docker-php-ext-enable mhash

# Ставим:
# 	- либу recode
# 	- php-расширения: librecode
RUN \
	apt install -y librecode-dev && \
	docker-php-ext-install -j$(nproc) recode

# Ставим:
# 	- либу soap
# 	- php-расширения: libxml2
RUN \
	apt install -y libxml2-dev && \
	docker-php-ext-install -j$(nproc) soap

# Ставим:
# 	- php-расширения: redis
RUN \
	pecl install redis-2.2.8 && \
	docker-php-ext-enable redis

### }}}


##########################################################################################################################################


# Копируем конфигурации Hg
COPY ./hg/hgrc /root/hgrc

# Для hgrc проставляем права:
RUN \
	chown root:root /root/hgrc


##########################################################################################################################################


# Копируем дополнительные файл composer
COPY ./.composer/auth.json /root/.composer/auth.json
COPY ./.composer/config.json /root/.composer/config.json

# Монтируем директорию кеша для composer
VOLUME /root/.composer/cache

# Копируем php.ini + доп конфигурации php
COPY /conf/php/global/php.ini /usr/local/etc/php/php.ini

# Копируем конфигурации Apache
COPY ./conf/apache /etc/apache2

# Монтируем локальные конфигурации Apache, тем самым мы избавляем себя от необходимости, при изменении конфигов, заново пересобирать контейнеры
VOLUME /etc/apache2/local

# Монтируем рабочую директорию с php-проектом
VOLUME /var/www/html

# Монтируем директорию логов
VOLUME /var/log

# Выводим наружу порт
EXPOSE 8080

# Выводим наружу порты для МС-ов
EXPOSE 5000-5500

##########################################################################################################################################


# Docker: запуск нескольких приложений в одном контейнере
# Источники:
#	- https://blog.amartynov.ru/docker-%D0%B7%D0%B0%D0%BF%D1%83%D1%81%D0%BA-%D0%BD%D0%B5%D1%81%D0%BA%D0%BE%D0%BB%D1%8C%D0%BA%D0%B8%D1%85-%D0%BF%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B9/
COPY start.sh /
RUN chmod +x /start.sh
CMD ["/start.sh"]
