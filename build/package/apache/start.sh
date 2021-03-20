#!/bin/bash

stop_requested=false
trap "stop_requested=true" TERM INT

wait_signal() {
    while ! $stop_requested; do
        sleep 1
    done
}

wait_exit() {
    while pidof $1; do
        sleep 1
    done
}

# Инициализация
chown -R www-data:www-data /var/www/html
sleep 20
php /var/www/html/yii migrate --interactive=0

# Всё запускаем
apache2ctl start

# Ждём SIGTERM или SIGINT
wait_signal

# Запрашиваем остановку
apache2ctl stop

# Ждём завершения процессов по их названию
wait_exit "apache2"
