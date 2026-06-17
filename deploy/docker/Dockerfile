# 1. Этап сборки
FROM alt:p11 AS builder

# Установка зависимостей, необходимых для работы инсталлятора 1С
RUN apt-get update && apt-get install -y \
    wget \
    glibc-core \
    glibc-pthread \
    libstdc++6 \
    && apt-get clean

# Копирование и запуск установщика в unattended-режиме
COPY setup-training-8.5.1.1150-x86_64.run /tmp/install.run
RUN chmod +x /tmp/install.run && \
    /tmp/install.run --mode unattended

# 2. Финальный образ
FROM alt:p11

# Установка системных зависимостей для 1С и веб-сервера
RUN apt-get update && apt-get install -y \
    httpd2 \
    libxml2 \
    fontconfig \
    libgsf \
    libglib2 \
    libX11 \
    libXext \
    libXrender \
    libstdc++6 \
    glibc-core \
    glibc-pthread \
    libpthread-stubs \
    libnss-myhostname \
    sudo \
    # Создаем символическую ссылку, чтобы 1С находила httpd по стандартному пути
    && ln -s /usr/sbin/httpd2 /usr/sbin/httpd \
    && apt-get clean

# Копируем установленную платформу из образа-сборщика
COPY --from=builder /opt/1cv8t /opt/1cv8t

# Настройка веб-сервера и рабочих директорий
RUN echo "ServerName localhost" >> /etc/httpd2/conf/httpd2.conf \
    && mkdir -p /var/www/booking /var/run/httpd2 /var/log/httpd2 \
    # Устанавливаем владельцем apache2 (системный пользователь httpd2 в Alt Linux)
    && chown -R apache2:apache2 /var/www/booking /var/run/httpd2 /var/log/httpd2 /etc/httpd2 /opt/1cv8t

# Подготовка скрипта запуска
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]