#!/bin/bash
# Скрипт управления контейнером 1С-публикации

CONTAINER_NAME="academic-test"

case "$1" in
    republish)
        echo "Сброс публикации..."
        # Удаляем файл публикации, чтобы entrypoint.sh пересоздал его при перезапуске
        docker exec -it $CONTAINER_NAME rm -f /var/www/booking/default.vrd
        docker compose restart app
        echo "Публикация успешно сброшена."
        ;;
        
    edit-vrd)
        echo "Открываю конфигурацию..."
        # Извлекаем текущий файл конфигурации из контейнера на хост для редактирования
        docker cp $CONTAINER_NAME:/var/www/booking/default.vrd ./deploy/docker/default.vrd.tmp
        
        # Редактируем файл локально
        nano ./deploy/docker/default.vrd.tmp
        
        # Копируем отредактированный файл обратно в контейнер
        docker cp ./deploy/docker/default.vrd.tmp $CONTAINER_NAME:/var/www/booking/default.vrd
        rm ./deploy/docker/default.vrd.tmp
        
        echo "Обновление конфигурации Apache (httpd2)..."
        # Перезагружаем конфигурацию веб-сервера без остановки контейнера
        docker exec -it $CONTAINER_NAME httpd2 -k restart
        ;;
        
    logs)
        # Вывод последних 50 строк логов с отслеживанием изменений
        docker compose logs -f --tail=50
        ;;
        
    *)
        echo "Использование: $0 {republish|edit-vrd|logs}"
        exit 1
        ;;
esac