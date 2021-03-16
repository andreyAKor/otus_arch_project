#!/bin/sh

# Взято отсюда: https://blog.amartynov.ru/docker-%D0%B7%D0%B0%D0%BF%D1%83%D1%81%D0%BA-%D0%BD%D0%B5%D1%81%D0%BA%D0%BE%D0%BB%D1%8C%D0%BA%D0%B8%D1%85-%D0%BF%D1%80%D0%B8%D0%BB%D0%BE%D0%B6%D0%B5%D0%BD%D0%B8%D0%B9/
# {{{
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
# }}}

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
