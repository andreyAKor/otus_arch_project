FROM php:7.4-apache

RUN apt update && \
    apt install -y libfreetype6-dev libjpeg62-turbo-dev libmcrypt-dev libpng-dev libsodium-dev git zlib1g-dev openssl unzip libzip-dev

RUN echo "Europe/Moscow" > /etc/timezone && \
	dpkg-reconfigure --frontend noninteractive tzdata

RUN docker-php-ext-install -j$(nproc) bcmath sodium gd calendar gettext pdo pdo_mysql shmop sockets sysvmsg sysvsem sysvshm pcntl zip

RUN a2enmod rewrite

COPY ./configs/apache/conf/otus_arch_project.conf /etc/apache2/sites-enabled
COPY ./website /var/www/html

WORKDIR /var/www/html

RUN php ./composer.phar install && \
    php ./init --env=Docker --overwrite=All --delete=All

WORKDIR /

EXPOSE 8080

COPY ./build/package/apache/start.sh /
RUN chmod +x /start.sh
CMD ["/start.sh"]
