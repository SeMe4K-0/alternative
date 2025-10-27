# # api_app.py
# from fastapi import FastAPI, HTTPException
# from pydantic import BaseModel, Field, validator
# from typing import Optional
# import uvicorn
# import math
#
# # Импортируем функции из вашего модуля (положите main_gpt_with_mars_jpl.py рядом)
# # В нём должны быть: mechanism2_closest_approach_to_earth и validate_elements_input (или аналог)
# # from main_gpt import mechanism2_closest_approach_to_earth, validate_elements_input
# from main_gpt import mechanism2_closest_approach_to_earth
#
# # astropy check (для валидации времени)
# try:
#     from astropy.time import Time
#
#     ASTROPY_AVAILABLE = True
# except Exception:
#     ASTROPY_AVAILABLE = False
#
# app = FastAPI(
#     title="Closest Approach API",
#     description="Endpoint computes closest approach (t_closest_rfc3339, distance_AU) "
#                 "for orbital elements a,e,i,Omega,omega,Tp.",
#     version="0.1"
# )
#
#
# # Pydantic model for input payload (no 'notes')
# class ElementsPayload(BaseModel):
#     a: float = Field(..., gt=0.0, description="Semi-major axis in AU")
#     e: float = Field(..., ge=0.0, description="Eccentricity (>=0)")
#     i: float = Field(..., ge=0.0, lt=360.0, description="Inclination in degrees")
#     Omega: float = Field(..., ge=0.0, lt=360.0, description="Longitude of ascending node in degrees")
#     omega: float = Field(..., ge=0.0, lt=360.0, description="Argument of perihelion in degrees")
#     Tp: str = Field(..., description="Time of perihelion (RFC3339/ISO e.g. 2025-05-28T13:42:39Z)")
#     # optional control parameters
#     coarse_N: Optional[int] = Field(1000, gt=50, lt=20000, description="Coarse grid size for search (default 1000)")
#     search_years: Optional[float] = Field(None, gt=0.0, description="Optional search span in years")
#
#     @validator('Tp')
#     def tp_must_be_rfc3339(cls, v):
#         # minimal syntactic check — astropy will do full parse if available
#         if not isinstance(v, str) or 'T' not in v or 'Z' not in v:
#             raise ValueError("Tp must be RFC3339 string like 'YYYY-MM-DDTHH:MM:SSZ'")
#         return v
#
#
# @app.on_event("startup")
# def startup_event():
#     # Preconfigure astropy ephemeris if available to reduce first-request delay.
#     if ASTROPY_AVAILABLE:
#         try:
#             from astropy.coordinates import solar_system_ephemeris
#             # Prefer local DE file if you already have one; set('de440') will download if needed.
#             solar_system_ephemeris.set('de440')
#         except Exception:
#             # ignore — endpoint will fallback as in core module
#             pass
#
#
# @app.post("/closest-approach", tags=["compute"])
# def closest_approach(payload: ElementsPayload):
#     # Build elements dict for core function
#     elements = {
#         'semi_major_axis': payload.a,
#         'eccentricity': payload.e,
#         'inclination': payload.i,
#         'lon_ascending_node': payload.Omega,
#         'arg_periapsis': payload.omega,
#         'time_perihelion': payload.Tp
#     }
#
#     # Validate (use your validate function if present)
#     # try:
#     #     validate_elements_input(elements)
#     # except Exception as ex:
#     #     raise HTTPException(status_code=400, detail=f"Invalid input elements: {ex}")
#
#     # Call core computation
#     try:
#         # Pass coarse_N and search_years if present
#         res = mechanism2_closest_approach_to_earth(elements,
#                                                    search_years=payload.search_years,
#                                                    coarse_N=payload.coarse_N)
#     except Exception as ex:
#         # Unexpected internal error
#         raise HTTPException(status_code=500, detail=f"Computation failed: {ex}")
#
#     # Return only required fields as requested by you
#     return {
#         "approach_date": res.get('t_closest_rfc3339'),
#         "distance_au": res.get('distance_AU'),
#         "distance_km": res.get('distance_AU') * 149597870.7
#     }
#
#
# # Optional health endpoint
# @app.get("/health")
# def health():
#     return {"status": "ok", "astropy": ASTROPY_AVAILABLE}
#
#
# # If you want to run locally via `python api_app.py`
# if __name__ == "__main__":
#     # For debugging only. In production use uvicorn/gunicorn with workers.
#     uvicorn.run("api_app:app", host="0.0.0.0", port=5001, reload=False)
#
# # api_app.py
# from fastapi import FastAPI, HTTPException
# from pydantic import BaseModel, Field
# from typing import List
# import uvicorn
# import datetime
#
# # Импортируем только новую функцию-обертку из вашего файла с логикой
#
# app = FastAPI(
#     title="Orbit Calculation Service",
#     description="API, совместимое с Go-клиентом для расчета орбит и сближений.",
#     version="1.0.0"
# )
#
#
# # --- Модели, точно соответствующие структурам в Go ---
#
# # Ручка 1: /calculate-orbit (заглушка)
# class Observation(BaseModel):
#     # Это предположение о структуре, т.к. в Go она импортируется из models
#     observed_at: str
#     right_ascension: float
#     declination: float
#
#
# class OrbitCalculationRequest(BaseModel):
#     observations: List[Observation]
#
#
# class OrbitCalculationResponse(BaseModel):
#     semi_major_axis: float
#     eccentricity: float
#     inclination: float
#     lon_ascending_node: float
#     arg_periapsis: float
#     time_perihelion: str  # В Go это time.Time, здесь будет строка RFC3339
#
#
# # Ручка 2: /calculate-approach
# class ApproachCalculationRequest(BaseModel):
#     semi_major_axis: float
#     eccentricity: float
#     inclination: float
#     lon_ascending_node: float
#     arg_periapsis: float
#     time_perihelion: str
#     # Добавляем опциональные поля для сохранения логики
#     coarse_N: Optional[int] = Field(1000, gt=50, lt=20000)
#     search_years: Optional[float] = Field(None, gt=0.0)
#
#
# class ApproachCalculationResponse(BaseModel):
#     approach_date: str
#     distance_au: float
#     distance_km: float
#
#
# # --- Эндпоинты API ---
#
# @app.post("/calculate-orbit", response_model=OrbitCalculationResponse, tags=["Orbit Calculation"])
# def calculate_orbit_stub(request: OrbitCalculationRequest):
#     """
#     ЗАГЛУШКА для первой ручки. Возвращает статичные данные,
#     аналогичные mock-функции в Go.
#     """
#     print(f"Received {len(request.observations)} observations for '/calculate-orbit'")
#     # Возвращаем примерные данные
#     future_time = datetime.datetime.utcnow() + datetime.timedelta(days=240)
#     return OrbitCalculationResponse(
#         semi_major_axis=3.1,
#         eccentricity=0.75,
#         inclination=22.0,
#         lon_ascending_node=150.0,
#         arg_periapsis=110.0,
#         time_perihelion=future_time.strftime("%Y-%m-%dT%H:%M:%SZ")
#     )
#
# from main_gpt import mechanism2_closest_approach_to_earth
#
# # --- Запуск сервера ---
# if __name__ == "__main__":
#     uvicorn.run("api_app:app", host="0.0.0.0", port=5001, reload=True)


# api_app.py
import datetime
import hashlib
import random
from typing import List, Optional

import uvicorn
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

# Импортируем вашу основную функцию вычислений из main_gpt.py
from main_gpt import mechanism2_closest_approach_to_earth

# astropy check (для валидации времени)
try:
    from astropy.time import Time
    from astropy.coordinates import solar_system_ephemeris

    ASTROPY_AVAILABLE = True
except Exception:
    ASTROPY_AVAILABLE = False

app = FastAPI(
    title="Orbit and Approach API",
    description="API для расчета орбит (заглушка) и сближений.",
    version="1.3.0"
)


# --- МОДЕЛИ ДАННЫХ (ЕДИНЫЕ ДЛЯ ВСЕГО ПРИЛОЖЕНИЯ) ---

# Модель для ручки /calculate-orbit
class Observation(BaseModel):
    observed_at: str
    right_ascension: float
    declination: float


class OrbitCalculationRequest(BaseModel):
    observations: List[Observation]


# Модель для ручки /calculate-approach (ТЕПЕРЬ ИСПОЛЬЗУЕТСЯ И В /closest-approach)
class ApproachCalculationRequest(BaseModel):
    # Эти имена точно соответствуют тому, что отправляет Go
    semi_major_axis: float
    eccentricity: float
    inclination: float
    lon_ascending_node: float
    arg_periapsis: float
    time_perihelion: str
    # Опциональные параметры
    coarse_N: Optional[int] = Field(1000, gt=50, lt=20000)
    search_years: Optional[float] = Field(None, gt=0.0)


class ApproachCalculationResponse(BaseModel):
    approach_date: str
    distance_au: float
    distance_km: float


# --- ЭНДПОИНТЫ API ---

@app.on_event("startup")
def startup_event():
    if ASTROPY_AVAILABLE:
        try:
            solar_system_ephemeris.set('de440')
        except Exception:
            pass


@app.post("/calculate-orbit", tags=["Orbit Calculation"])
def calculate_orbit_stub(payload: OrbitCalculationRequest):
    """
    РУЧКА 1: Заглушка, которая возвращает детерминированные, но
    псевдослучайные параметры орбиты. Для одного и того же набора
    наблюдений результат всегда будет одинаковым.
    """
    print(f"Получено {len(payload.observations)} наблюдений для '/calculate-orbit'")

    if not payload.observations:
        return {"error": "Не получено ни одного наблюдения."}

    # 1. Создаем стабильную строку из входных данных.
    # Чтобы порядок наблюдений в JSON не влиял на результат,
    # мы сначала сортируем их.
    # Преобразуем каждое наблюдение в кортеж (tuple), чтобы их можно было отсортировать.
    sorted_observations_tuples = sorted(
        (
            obs.observed_at,
            obs.right_ascension,
            obs.declination
        )
        for obs in payload.observations
    )

    # Преобразуем отсортированный список кортежей в строку.
    # Это будет нашей уникальной "подписью" для данного набора наблюдений.
    data_string = str(sorted_observations_tuples)

    # 2. Создаем "сид" (seed) путем хеширования этой строки.
    # Алгоритм SHA-256 гарантирует, что даже малейшее изменение во входных данных
    # приведет к совершенно другому хешу (и, соответственно, другому сиду).
    seed = hashlib.sha256(data_string.encode('utf-8')).hexdigest()

    # 3. Устанавливаем этот сид для генератора случайных чисел.
    # Теперь все последующие вызовы random.uniform() и random.randint()
    # будут генерировать одну и ту же последовательность чисел для этого сида.
    random.seed(seed)

    # 4. Генерируем "случайные" параметры, которые теперь полностью предсказуемы.
    semi_major_axis = random.uniform(0.5, 40.0)
    eccentricity = random.uniform(0.0, 0.99)
    inclination = random.uniform(0.0, 180.0)
    lon_ascending_node = random.uniform(0.0, 360.0)
    arg_periapsis = random.uniform(0.0, 360.0)

    random_days_in_future = random.randint(30, 730)
    # Используем фиксированную дату как точку отсчета, чтобы время тоже было детерминированным
    base_time = datetime.datetime(2025, 1, 1, tzinfo=datetime.timezone.utc)
    future_time = base_time + datetime.timedelta(days=random_days_in_future)
    time_perihelion_str = future_time.strftime("%Y-%m-%dT%H:%M:%SZ")

    # --- Возвращаем результат ---
    return {
        "semi_major_axis": semi_major_axis,
        "eccentricity": eccentricity,
        "inclination": inclination,
        "lon_ascending_node": lon_ascending_node,
        "arg_periapsis": arg_periapsis,
        "time_perihelion": time_perihelion_str
    }


@app.post("/closest-approach", response_model=ApproachCalculationResponse, tags=["Close Approach"])
def closest_approach(payload: ApproachCalculationRequest):  # <-- ИСПРАВЛЕНО ЗДЕСЬ
    """
    РУЧКА 2: Теперь использует ту же модель, что и Go-клиент.
    """
    # Преобразуем поля из payload в словарь `elements` с короткими именами,
    # который ожидает ваша функция `mechanism2_closest_approach_to_earth`.
    elements = {
        'a': payload.semi_major_axis,
        'e': payload.eccentricity,
        'i': payload.inclination,
        'Omega': payload.lon_ascending_node,
        'omega': payload.arg_periapsis,
        'Tp': payload.time_perihelion
    }

    try:
        res = mechanism2_closest_approach_to_earth(
            elements,
            search_years=payload.search_years,
            coarse_N=payload.coarse_N
        )
    except Exception as ex:
        raise HTTPException(status_code=500, detail=f"Computation failed: {str(ex)}")

    distance_au_val = res.get('distance_AU')
    if distance_au_val is None:
        raise HTTPException(status_code=500, detail="Computation did not return a valid distance.")

    return {
        "approach_date": res.get('t_closest_rfc3339'),
        "distance_au": distance_au_val,
        "distance_km": distance_au_val * 149597870.7
    }


@app.get("/health", tags=["Health"])
def health():
    return {"status": "ok", "astropy": ASTROPY_AVAILABLE}


# --- ЗАПУСК СЕРВЕРА ---
if __name__ == "__main__":
    uvicorn.run("api_app:app", host="0.0.0.0", port=5001, reload=True)