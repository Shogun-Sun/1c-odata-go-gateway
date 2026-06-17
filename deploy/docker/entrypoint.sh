#!/bin/bash
set -e

# Конфигурация для Alt Linux p11
WEB_USER="apache2"
BASE_DIR="/base"
WWW_DIR="/var/www/booking"
HTTPD_CONF="/etc/httpd2/conf/httpd2.conf"
HTTPD_PID="/var/run/httpd2/httpd2.pid"

# Функция логирования
log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] [entrypoint] $1"
}

# Динамическая настройка UID для синхронизации с хостом
if [ -n "$UID" ] && [ "$UID" != "0" ]; then
    log "Adjusting $WEB_USER UID to $UID..."
    usermod -u "$UID" "$WEB_USER" || log "Warning: UID update skipped (already set or conflict)."
fi

# Очистка базы данных от старых блокировок
if [ -d "$BASE_DIR" ]; then
    log "Cleaning up 1C lock files in $BASE_DIR..."
    find "$BASE_DIR" -type f \( -name "*.lck" -o -name "*.cfl" \) -delete
fi

# Применение прав доступа для Alt Linux
log "Fixing permissions..."
mkdir -p "$BASE_DIR" "$WWW_DIR" /var/run/httpd2 /var/log/httpd2
chown -R "$WEB_USER":"$WEB_USER" "$BASE_DIR" "$WWW_DIR" /var/run/httpd2 /var/log/httpd2
chmod -R 755 "$BASE_DIR" "$WWW_DIR"

# Поиск утилиты публикации
WEBINST_PATH=$(find /opt/1cv8t -name "webinstt" | head -n 1)
if [ -z "$WEBINST_PATH" ]; then
    log "ERROR: webinstt not found in /opt/1cv8t."
    exit 1
fi

# Автоматическая публикация базы (путь к конфигу Alt Linux)
if [ ! -f "$WWW_DIR/default.vrd" ]; then
    log "Publishing database..."
    "$WEBINST_PATH" -apache24 -wsdir booking -dir "$WWW_DIR" -connstr "File=$BASE_DIR;" -confpath "$HTTPD_CONF"
    
    cat > "$WWW_DIR/default.vrd" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<point xmlns="http://v8.1c.ru/8.2/virtual-resource-system" 
       xmlns:xs="http://www.w3.org/2001/XMLSchema" 
       xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
       base="/booking" 
       ib="File=$BASE_DIR;">
    <ws pointEnableCommon="true"/>
    <standardOdata enable="true" reuseSessions="dontuse" sessionMaxAge="0" poolSize="0" poolTimeout="0"/>
    <analytics enable="true" sessionMaxAge="1200" poolSize="500" poolTimeout="5"/>
</point>
EOF
    chown "$WEB_USER":"$WEB_USER" "$WWW_DIR/default.vrd"
    chmod 644 "$WWW_DIR/default.vrd"
    
    log "Database published successfully."
fi

# Подготовка и запуск Apache (httpd2)
log "Starting httpd2..."
rm -f "$HTTPD_PID"

# Запуск httpd2
exec /usr/sbin/httpd2 -D FOREGROUND