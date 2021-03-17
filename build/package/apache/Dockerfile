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


# Ставим git
# Источники:
# 	- https://toster.ru/q/250345
RUN \
	apt install -y git


##########################################################################################################################################



# Подготовка
RUN \
	mkdir /usr/local/etc/php/conf.d/

# Ставим:
# 	- php-расширения: bcmath, calendar, gettext, mysql, pdo_mysql, shmop, sockets, sysvmsg, sysvsem, sysvshm, pcntl
#RUN docker-php-ext-install -j$(nproc) bcmath calendar gettext mysql pdo_mysql shmop sockets sysvmsg sysvsem sysvshm pcntl



##########################################################################################################################################




# Копируем php.ini + доп конфигурации php
#COPY ./configs/apache/conf/php/global/php.ini /usr/local/etc/php/php.ini

# Копируем конфигурации Apache
#COPY ./configs/apache/conf/apache /etc/apache2

# Монтируем локальные конфигурации Apache, тем самым мы избавляем себя от необходимости, при изменении конфигов, заново пересобирать контейнеры
VOLUME /etc/apache2/local

# Выводим наружу порт
EXPOSE 8080

##########################################################################################################################################


# Docker: запуск нескольких приложений в одном контейнере
# Источники:
#	- https://blog.amartynov.ru/docker-%D0%B7%D0%B0%D0%BF%D1%83%D1%81%D0%BA-%D0%BD%D0%B5%D1%81%D0%BA%D0%BE%D0%BB%D1%8C%D0%BA%D0%B8%D1%85-%D0%BF%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B9/
COPY ./build/package/apache/start.sh /
RUN chmod +x /start.sh
CMD ["/start.sh"]
