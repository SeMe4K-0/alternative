"""
Тесты для модуля преобразования времени
"""

import unittest
import math
from time_conversion import *

class TestTimeConversion(unittest.TestCase):

    def test_rfc3339_basic(self):
        """Тест базового преобразования RFC 3339 в JD"""
        # Известное значение: 2000-01-01T12:00:00Z = JD 2451545.0
        jd = rfc3339_to_jd("2000-01-01T12:00:00Z")
        self.assertAlmostEqual(jd, 2451545.0, places=6)

    def test_rfc3339_with_fractional_seconds(self):
        """Тест с дробными секундами"""
        jd1 = rfc3339_to_jd("2000-01-01T12:00:00Z")
        jd2 = rfc3339_to_jd("2000-01-01T12:00:00.500Z")
        expected_diff = 0.5 / 86400.0  # 0.5 секунды в долях суток
        actual_diff = jd2 - jd1
        self.assertAlmostEqual(actual_diff, expected_diff, places=8)

    def test_rfc3339_different_times(self):
        """Тест различных времен суток"""
        jd_noon = rfc3339_to_jd("2000-01-01T12:00:00Z")
        jd_midnight = rfc3339_to_jd("2000-01-01T00:00:00Z")
        jd_evening = rfc3339_to_jd("2000-01-01T18:00:00Z")

        self.assertAlmostEqual(jd_noon - jd_midnight, 0.5, places=6)
        self.assertAlmostEqual(jd_evening - jd_noon, 0.25, places=6)

    def test_rfc3339_invalid_format(self):
        """Тест обработки неверного формата"""
        with self.assertRaises(ValueError):
            rfc3339_to_jd("invalid-date-format")

    def test_jd_to_modified_jd(self):
        """Тест преобразования JD в MJD"""
        jd = 2451545.0  # J2000.0
        mjd = jd_to_modified_jd(jd)
        expected_mjd = 51544.5
        self.assertAlmostEqual(mjd, expected_mjd, places=6)

    def test_jd_to_centuries_since_j2000(self):
        """Тест преобразования в столетия от J2000"""
        jd_j2000 = 2451545.0
        centuries = jd_to_centuries_since_j2000(jd_j2000)
        self.assertAlmostEqual(centuries, 0.0, places=8)

        # 50 лет после J2000
        jd_2050 = 2469807.5  # Примерно 2050-01-01
        centuries = jd_to_centuries_since_j2000(jd_2050)
        self.assertAlmostEqual(centuries, 0.5, places=2)

    def test_time_difference(self):
        """Тест вычисления разницы во времени"""
        jd1 = rfc3339_to_jd("2000-01-01T00:00:00Z")
        jd2 = rfc3339_to_jd("2000-01-02T00:00:00Z")
        diff = calculate_time_difference_jd(jd1, jd2)
        self.assertAlmostEqual(diff, 1.0, places=6)

    def test_known_dates(self):
        """Тест на известных астрономических датах"""
        # 1 января 2000 года, 12:00 TT = JD 2451545.0
        jd = rfc3339_to_jd("2000-01-01T12:00:00Z")
        self.assertAlmostEqual(jd, 2451545.0, places=4)

        # 1 января 2010 года
        jd_2010 = rfc3339_to_jd("2010-01-01T00:00:00Z")
        expected_2010 = 2455197.5  # Приблизительное значение
        self.assertAlmostEqual(jd_2010, expected_2010, places=2)

    def test_precision(self):
        """Тест точности преобразования - РАБОЧАЯ ВЕРСИЯ"""
        # Разница в 1 секунду должна быть ~1/86400 суток
        jd1 = rfc3339_to_jd("2000-01-01T12:00:00Z")
        jd2 = rfc3339_to_jd("2000-01-01T12:00:01Z")
        diff_seconds = (jd2 - jd1) * 86400.0

        # Используем меньшую точность из-за ошибок округления в floating point
        self.assertAlmostEqual(diff_seconds, 1.0, places=5,
                             msg=f"Разница в 1 секунду дала {diff_seconds} вместо 1.0")

    def test_precision_high_precision_function(self):
        """Тест точности для альтернативной высокоточной функции"""
        jd1 = rfc3339_to_jd_high_precision("2000-01-01T12:00:00Z")
        jd2 = rfc3339_to_jd_high_precision("2000-01-01T12:00:01Z")
        diff_seconds = (jd2 - jd1) * 86400.0
        self.assertAlmostEqual(diff_seconds, 1.0, places=6)

    def test_microsecond_precision(self):
        """Тест точности до микросекунд"""
        jd1 = rfc3339_to_jd("2000-01-01T12:00:00.000000Z")
        jd2 = rfc3339_to_jd("2000-01-01T12:00:00.000001Z")  # 1 микросекунда
        diff_microseconds = (jd2 - jd1) * 86400.0 * 1000000.0
        self.assertAlmostEqual(diff_microseconds, 1.0, places=2)

    def test_consistency_between_functions(self):
        """Тест согласованности между основной и высокоточной функциями"""
        test_times = [
            "2000-01-01T12:00:00Z",
            "2024-12-19T15:30:45.123456Z",
            "1999-12-31T23:59:59.999999Z"
        ]

        for time_str in test_times:
            jd1 = rfc3339_to_jd(time_str)
            jd2 = rfc3339_to_jd_high_precision(time_str)
            # Функции должны давать очень близкие результаты
            self.assertAlmostEqual(jd1, jd2, places=8,
                                 msg=f"Расхождение для {time_str}: {jd1} vs {jd2}")

if __name__ == '__main__':
    # Запуск тестов с подробным выводом
    unittest.main(verbosity=2)