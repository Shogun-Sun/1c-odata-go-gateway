#!/bin/bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

BACKUP_SCRIPT="${SCRIPT_DIR}/backup_academic_booking.sh"

CRON_JOB="1 23 * * * $BACKUP_SCRIPT"

case "${1}" in
    on)
        if crontab -l 2>/dev/null | grep -Fq "$BACKUP_SCRIPT"; then
            echo "Помидорка уже на месте: задача для cron уже активирована."
        else
            (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
            echo "Задача успешно добавлена в crontab! Текущее расписание:"
            crontab -l | grep -F "$BACKUP_SCRIPT"
        fi
        ;;
        
    off)
        if ! crontab -l 2>/dev/null | grep -Fq "$BACKUP_SCRIPT"; then
            echo "Удалять нечего: задача в crontab отсутствовала."
        else
            crontab -l 2>/dev/null | grep -v -F "$BACKUP_SCRIPT" | crontab -
            echo "Задача бэкапа успешно удалена из crontab."
        fi
        ;;
        
    *)
        echo "Ошибка: неверный параметр."
        echo "Использование: $0 {on|off}"
        exit 1
        ;;
esac