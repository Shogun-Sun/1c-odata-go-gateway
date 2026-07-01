#!/bin/bash
export LANG=ru_RU.UTF-8

STORAGE="${1}"
SRV="${2}"
REF="${3}"
LOGIN="${4}"
IB_BASE_PATH="${5}"
KEEP_BACKUPS_FOR_LAST_DAYS="${6}"
BACKUP_TIMEOUT_SEC="${7}"
NON_ROOT_USER="${8}"

DISPLAY=":0"
TIMEOUT_COUNTER=0

#-----------------------------
export DISPLAY="${DISPLAY}"
export XAUTHORITY="/home/${NON_ROOT_USER}/.Xauthority"

bin1c="/opt/1cv8t/x86_64/8.5.1.1150/1cv8t"

infobase_suffix="$(date +%Y-%m-%d-%H-%M-%S).dt"
TMP_STORAGE="/opt/1c_backup_scripts/tmp_storage"
ib_tmp_filepath="${TMP_STORAGE}/${REF}__${infobase_suffix}"
ib_tmp_log="${ib_tmp_filepath}.log"
ib_filepath="${STORAGE}/${REF}__${infobase_suffix}"
ib_log="${ib_filepath}.log"

function rotate_backups(){
filenames_cmd="${STORAGE}/${REF}*"
current_unixtime=`date +%s --date="$(date +%Y-%m-%d) 00:00:00"`
test_unixtime=$(($current_unixtime-$KEEP_BACKUPS_FOR_LAST_DAYS*60*60*24))

for f in $(/bin/ls $filenames_cmd); do
datetime_from_file=`echo ${f} | awk -F '__' {'print $2'} | awk -F '.' {'print $1'}`

Y=`echo ${datetime_from_file} | awk -F '-' {'print $1'}`
M=`echo ${datetime_from_file} | awk -F '-' {'print $2'}`
D=`echo ${datetime_from_file} | awk -F '-' {'print $3'}`
h=`echo ${datetime_from_file} | awk -F '-' {'print $4'}`
m=`echo ${datetime_from_file} | awk -F '-' {'print $5'}`
s=`echo ${datetime_from_file} | awk -F '-' {'print $6'}`

unixtime_from_file=`date +%s --date="${Y}-${M}-${D} ${h}:${m}:${s}"`

if [[ $unixtime_from_file -lt $test_unixtime ]]; then
  echo "delete file: ${f}"
  /bin/rm ${f}
else
  echo "keep file: ${f}"
fi
done
}

function wait(){

  if [[ $TIMEOUT_COUNTER -ge $BACKUP_TIMEOUT_SEC ]]; then
    kill -9 ${1}
  fi

  if [ ! -d /proc/${1}/ ]; then
    nice_text="Выгрузка информационной базы успешно завершена"
    res_text=$(grep -o -i "${nice_text}" ${ib_tmp_log})

    # start apache2
    sudo systemctl start httpd2

    if [ "${res_text}" = "${nice_text}" ]; then
      /bin/cp ${ib_tmp_filepath} ${ib_filepath} && /bin/rm ${ib_tmp_filepath}
      /bin/cp ${ib_tmp_log} ${ib_log} && /bin/rm ${ib_tmp_log}
      rotate_backups
      exit 0
    fi

    exit 1
  fi
  sleep 1
  TIMEOUT_COUNTER=$(($TIMEOUT_COUNTER+1))
  wait ${1}
}

# prepare tmp storage
[ ! -d ${TMP_STORAGE} ] && sudo mkdir "${TMP_STORAGE}" && sudo chown ${NON_ROOT_USER}: "${TMP_STORAGE}" -R

# stop apache2
sudo systemctl stop httpd2
sleep 2

# Запуск файловой базы. Путь берется из параметров
${bin1c} CONFIG /N "${LOGIN}" /F "${IB_BASE_PATH}" /DumpIB "${ib_tmp_filepath}" /OUT "${ib_tmp_log}" &
main_process_pid=$!

wait ${main_process_pid}