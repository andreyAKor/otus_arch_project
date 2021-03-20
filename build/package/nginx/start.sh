#!/bin/sh

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
chown -R root:root /var/log/
chown -R nginx:nginx /var/log/nginx/
chmod -R 0777 /var/log

# Всё запускаем
nginx -g 'daemon off;'

# Ждём SIGTERM или SIGINT
wait_signal

# Запрашиваем остановку
pkill -x nginx

# Ждём завершения процессов по их названию
wait_exit "nginx"
