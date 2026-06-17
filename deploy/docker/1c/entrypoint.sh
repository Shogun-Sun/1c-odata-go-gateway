#!/bin/bash
set -e

# Используем переменные среды (с дефолтными значениями, если они не заданы в .env)
WEB_USER="${WEB_USER:-apache2}"
BASE_DIR="${BASE_DIR:-/base}"
WS_DIR="${WS_DIR:-booking}"
WWW_DIR="/var/www/$WS_DIR" # Путь к папке публикации теперь зависит от имени WS_DIR

HTTPD_CONF="/etc/httpd2/conf/httpd2.conf"
HTTPD_PID="/var/run/httpd2/httpd2.pid"

# Функция логирования
log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] [entrypoint] $1"
}

# 1. Динамическая настройка UID/GID (используем переданные параметры)
# Добавили проверку на GID, если нужно
if [ -n "$UID" ] && [ "$UID" != "0" ]; then
    log "Adjusting $WEB_USER UID to $UID..."
    usermod -u "$UID" "$WEB_USER" || log "Warning: UID update skipped."
fi
if [ -n "$GID" ]; then
    log "Adjusting $WEB_USER GID to $GID..."
    groupmod -g "$GID" "$WEB_USER" || log "Warning: GID update skipped."
fi

# 2. Очистка базы данных
if [ -d "$BASE_DIR" ]; then
    log "Cleaning up 1C lock files in $BASE_DIR..."
    find "$BASE_DIR" -type f \( -name "*.lck" -o -name "*.cfl" \) -delete
fi

# 3. Применение прав доступа
log "Fixing permissions..."
mkdir -p "$BASE_DIR" "$WWW_DIR" /var/run/httpd2 /var/log/httpd2
chown -R "$WEB_USER":"$WEB_USER" "$BASE_DIR" "$WWW_DIR" /var/run/httpd2 /var/log/httpd2
chmod -R 755 "$BASE_DIR" "$WWW_DIR"

# 4. Поиск утилиты публикации
WEBINST_PATH=$(find /opt/1cv8t -name "webinstt" | head -n 1)
if [ -z "$WEBINST_PATH" ]; then
    log "ERROR: webinstt not found in /opt/1cv8t."
    exit 1
fi

# 5. Автоматическая публикация базы
if [ ! -f "$WWW_DIR/default.vrd" ]; then
    log "Publishing database '$WS_DIR'..."
    "$WEBINST_PATH" -apache24 -wsdir "$WS_DIR" -dir "$WWW_DIR" -connstr "File=$BASE_DIR;" -confpath "$HTTPD_CONF"
    
    cat > "$WWW_DIR/default.vrd" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<point xmlns="http://v8.1c.ru/8.2/virtual-resource-system" 
       base="/$WS_DIR" 
       ib="File=$BASE_DIR;">
    <ws pointEnableCommon="true"/>
    <standardOdata enable="true" reuseSessions="dontuse" sessionMaxAge="0" poolSize="0" poolTimeout="0"/>
    <analytics enable="true" sessionMaxAge="1200" poolSize="500" poolTimeout="5"/>
</point>
EOF
    chown "$WEB_USER":"$WEB_USER" "$WWW_DIR/default.vrd"
    chmod 644 "$WWW_DIR/default.vrd"
    
    log "Database '$WS_DIR' published successfully."
fi

# 6. Запуск Apache
log "Starting httpd2..."
rm -f "$HTTPD_PID"
exec /usr/sbin/httpd2 -D FOREGROUND