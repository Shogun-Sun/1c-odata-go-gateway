# 1. Этап сборки
FROM alt:p11 AS builder

# Выносим имя установщика в ARG
ARG INSTALLER_NAME=setup-training-8.5.1.1150-x86_64.run

RUN apt-get update && apt-get install -y wget glibc-core glibc-pthread libstdc++6 && apt-get clean

# Копируем и запускаем установщик
COPY $INSTALLER_NAME /tmp/install.run
RUN chmod +x /tmp/install.run && /tmp/install.run --mode unattended

# 2. Финальный образ
FROM alt:p11

# Переменная окружения для легкого изменения конфигурации на этапе запуска
ENV WS_DIR=default_ws
ENV BASE_DIR=/default_base

# Установка зависимостей (оставляем как есть, это база)
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
    && ln -s /usr/sbin/httpd2 /usr/sbin/httpd \
    && apt-get clean

COPY --from=builder /opt/1cv8t /opt/1cv8t

# Настройка (статические пути лучше оставить здесь, так как они зависят от дистрибутива)
RUN echo "ServerName localhost" >> /etc/httpd2/conf/httpd2.conf \
    && mkdir -p /var/www/$WS_DIR /var/run/httpd2 /var/log/httpd2 \
    && chown -R apache2:apache2 /var/www/$WS_DIR /var/run/httpd2 /var/log/httpd2 /etc/httpd2 /opt/1cv8t

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]