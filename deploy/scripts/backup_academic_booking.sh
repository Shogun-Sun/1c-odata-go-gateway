#!/bin/bash

### PARAMETERS ###
NON_ROOT_USER=                 # Имя пользователя в Системе
STORAGE="/home/${NON_ROOT_USER}/backups"      # Каталог для сохранения резервных копий
IB_BASE_PATH= # Полный путь к папке с базой 1С
SRV="localhost"                       
REF="base1"                                   # Имя базы для имени файла бэкапа
LOGIN=                         # Логин администратора в 1С
KEEP_BACKUPS_FOR_LAST_DAYS="2"            
BACKUP_TIMEOUT_SEC="300"                  

# ------------------------------------------------------------------------------------------------

### BOT PARAMETERS
TELEGRAM_BOT_TOKEN=""             
TELEGRAM_DESTINATION_GROUP_ID=""  

function send_message_telegram(){
exit 0  # отсылка в телеграм отключена
message=${1}
curl -s -X POST \
     -H 'Content-Type: application/json' \
     -d "{\"chat_id\": \"${TELEGRAM_DESTINATION_GROUP_ID}\", \"text\": \"${message}\", \"disable_notification\": true}" \
     https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage > /dev/null
}
# ------------------------------------------------------------------------------------------------

### START ###
BASE=`dirname ${0}`

# Создаем папку для бэкапов, если её нет, и отдаем права текущему юзеру
[ ! -d "${STORAGE}" ] && mkdir -p "${STORAGE}" && chown ${NON_ROOT_USER}: "${STORAGE}"

# Передаем параметры в базовый скрипт (включая путь к базе и имя юзера)
${BASE}/backup_infobase_1c_generic.sh "${STORAGE}" "${SRV}" "${REF}" "${LOGIN}" "${IB_BASE_PATH}" "${KEEP_BACKUPS_FOR_LAST_DAYS}" "${BACKUP_TIMEOUT_SEC}" "${NON_ROOT_USER}"
res=$?

if [ ${res} -eq 0 ]; then
send_message_telegram "[Успешно]\nРезервное копирование информационной базы '${REF}' выполнено успешно"
else
send_message_telegram "[Ошибка]\nРезервное копирование информационной базы '${REF}' не выполнено"
fi
