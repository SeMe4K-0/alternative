"""
Модуль для преобразования временных форматов в астрономические величины
"""

import math
from datetime import datetime, timezone
from typing import Tuple

def rfc3339_to_jd(time_str: str) -> float:
    """
    Конвертирует время из формата RFC 3339 в Юлианскую дату (JD)

    Исправленная версия с повышенной точностью для малых интервалов времени.

    Args:
        time_str: Строка времени в формате RFC 3339 (например, "2024-12-19T12:00:00Z")

    Returns:
        float: Юлианская дата

    Raises:
        ValueError: Если строка имеет неверный формат
    """
    try:
        # Парсим строку RFC 3339
        if time_str.endswith('Z'):
            dt = datetime.fromisoformat(time_str[:-1]).replace(tzinfo=timezone.utc)
        else:
            dt = datetime.fromisoformat(time_str)

        # Приводим к UTC если нужно
        dt_utc = dt.astimezone(timezone.utc)

        # Извлекаем компоненты даты
        year = dt_utc.year
        month = dt_utc.month
        day = dt_utc.day
        hour = dt_utc.hour
        minute = dt_utc.minute
        second = dt_utc.second
        microsecond = dt_utc.microsecond

        # Вычисляем общее количество секунд в сутках с высокой точностью
        total_seconds = (hour * 3600.0 +
                        minute * 60.0 +
                        second +
                        microsecond / 1000000.0)

        fractional_day = total_seconds / 86400.0

        # Вычисляем Юлианскую дату с использованием алгоритма из Meeus "Astronomical Algorithms"
        if month <= 2:
            year -= 1
            month += 12

        a = year // 100
        b = 2 - a + (a // 4)

        # Используем точные вычисления с float для избежания целочисленных ошибок
        jd_day = (int(365.25 * (year + 4716)) +
                 int(30.6001 * (month + 1)) +
                 day + b - 1524.5)

        jd = jd_day + fractional_day

        return jd

    except (ValueError, AttributeError) as e:
        raise ValueError(f"Неверный формат времени RFC 3339: {time_str}") from e


def rfc3339_to_jd_high_precision(time_str: str) -> float:
    """
    Альтернативная реализация с еще более высокой точностью
    используя алгоритм из NASA/JPL стандартов
    """
    try:
        if time_str.endswith('Z'):
            dt = datetime.fromisoformat(time_str[:-1]).replace(tzinfo=timezone.utc)
        else:
            dt = datetime.fromisoformat(time_str)

        dt_utc = dt.astimezone(timezone.utc)

        year = dt_utc.year
        month = dt_utc.month
        day = dt_utc.day

        # Вычисляем время в секундах с наносекундной точностью
        total_nanoseconds = (dt_utc.hour * 3600000000000 +
                           dt_utc.minute * 60000000000 +
                           dt_utc.second * 1000000000 +
                           dt_utc.microsecond * 1000)

        fractional_day = total_nanoseconds / 86400000000000.0

        # Алгоритм из Meeus с модификациями для высокой точности
        a = (14 - month) // 12
        y = year + 4800 - a
        m = month + 12 * a - 3

        jd_day = (day +
                 (153 * m + 2) // 5 +
                 365 * y +
                 y // 4 -
                 y // 100 +
                 y // 400 -
                 32045.5)

        return jd_day + fractional_day

    except (ValueError, AttributeError) as e:
        raise ValueError(f"Неверный формат времени RFC 3339: {time_str}") from e


def jd_to_modified_jd(jd: float) -> float:
    """
    Конвертирует Юлианскую дату в модифицированную Юлианскую дату (MJD)
    """
    return jd - 2400000.5


def jd_to_centuries_since_j2000(jd: float) -> float:
    """
    Конвертирует Юлианскую дату в юлианские столетия от эпохи J2000.0
    """
    return (jd - 2451545.0) / 36525.0


def calculate_time_difference_jd(jd1: float, jd2: float) -> float:
    """
    Вычисляет разницу во времени между двумя Юлианскими датами в сутках
    """
    return jd2 - jd1