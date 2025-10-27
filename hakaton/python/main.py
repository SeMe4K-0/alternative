# main_gpt_with_mars_jpl.py
# Механизм №1: заглушка (Mars elements)
# Механизм №2: ищет ближайшее приближение к Земле — теперь позиция Земли и Марса берутся из astropy/JPL ephemeris (если доступно)

import math
import datetime
import numpy as np

# --- astropy imports (для точной позиции Земли и Марса) ---
try:
    from astropy.time import Time
    from astropy.coordinates import get_body_barycentric, solar_system_ephemeris
    import astropy.units as u
    ASTROPY_AVAILABLE = True
except Exception:
    ASTROPY_AVAILABLE = False

# --- Вспомогательные функции ---
# def rfc3339_to_jd(rfc3339_str):
#     dt = datetime.datetime.strptime(rfc3339_str, "%Y-%m-%dT%H:%M:%SZ")
#     unix = (dt - datetime.datetime(1970,1,1)).total_seconds()
#     jd = unix / 86400.0 + 2440587.5
#     return jd

def rfc3339_to_jd(rfc3339_str: str) -> float:
    """
    Конвертирует строку времени RFC3339/ISO8601 в юлианскую дату.
    Корректно обрабатывает часовые пояса (и 'Z', и '+03:00').
    """
    # Python's fromisoformat идеально подходит для этого, но не любит 'Z'.
    # Заменим 'Z' на стандартный эквивалент, чтобы парсер был един.
    if rfc3339_str.endswith('Z'):
        rfc3339_str = rfc3339_str[:-1] + '+00:00'

    # Эта функция автоматически парсит строки с часовыми поясами, как "+03:00"
    dt_aware = datetime.datetime.fromisoformat(rfc3339_str)

    # .timestamp() корректно преобразует время в Unix timestamp с учетом часового пояса,
    # приводя его к UTC.
    unix_timestamp = dt_aware.timestamp()

    jd = unix_timestamp / 86400.0 + 2440587.5
    return jd

def jd_to_datetime_utc(jd):
    unix = (jd - 2440587.5) * 86400.0
    return datetime.datetime.utcfromtimestamp(unix)

def jd_to_rfc3339(jd):
    dt = jd_to_datetime_utc(jd)
    return dt.strftime("%Y-%m-%dT%H:%M:%SZ")

# Gaussian gravitational constant (AU^(3/2)/day)
k_gauss = 0.01720209895
mu = k_gauss**2  # AU^3 / day^2

def kepler_E_from_M(M, e, tol=1e-12):
    M = (M + math.pi) % (2*math.pi) - math.pi
    if abs(e) < 1e-8:
        return M
    if abs(M) < math.pi/6:
        E = M + e*math.sin(M)
    else:
        E = M
    for _ in range(200):
        f = E - e*math.sin(E) - M
        fp = 1 - e*math.cos(E)
        dE = -f / fp
        E = E + dE
        if abs(dE) < tol:
            break
    return E

def perifocal_to_ecliptic(r_pf, i_deg, Omega_deg, omega_deg):
    i = math.radians(i_deg)
    Omega = math.radians(Omega_deg)
    omega = math.radians(omega_deg)
    cosO = math.cos(Omega); sinO = math.sin(Omega)
    cosi = math.cos(i); sini = math.sin(i)
    cosw = math.cos(omega); sinw = math.sin(omega)
    R = np.array([
        [ cosO*cosw - sinO*sinw*cosi, -cosO*sinw - sinO*cosw*cosi, sinO*sini],
        [ sinO*cosw + cosO*sinw*cosi, -sinO*sinw + cosO*cosw*cosi, -cosO*sini],
        [ sinw*sini,                    cosw*sini,                    cosi     ]
    ])
    return R.dot(r_pf)

def position_from_elements(a, e, i_deg, Omega_deg, omega_deg, Tp_jd, t_jd):
    """Analytical Kepler position (heliocentric, ecliptic) in AU"""
    n = math.sqrt(mu / (a**3))
    dt = t_jd - Tp_jd
    M = n * dt
    M = M % (2*math.pi)
    E = kepler_E_from_M(M, e)
    cosE = math.cos(E); sinE = math.sin(E)
    sqrt_1_e2 = math.sqrt(1 - e*e)
    r_scalar = a * (1 - e*cosE)
    nu = 2*math.atan2(sqrt_1_e2 * sinE, cosE - e)
    x_pf = r_scalar * math.cos(nu)
    y_pf = r_scalar * math.sin(nu)
    r_pf = np.array([x_pf, y_pf, 0.0])
    r_ecl = perifocal_to_ecliptic(r_pf, i_deg, Omega_deg, omega_deg)
    return r_ecl  # AU, heliocentric, ecliptic frame

# --- Механизм №1 (заглушка) ---
def mechanism1_mars_stub():

    res = {
            "semi_major_axis": 1.523666696900072,
            "eccentricity": 0.09349252902810183,
            "inclination": 24.67729278671966,
            "lon_ascending_node": 3.365413082947441,
            "arg_periapsis": 333.0466125077485,
            "time_perihelion": "2025-05-28T13:42:39Z",
            "epoch_jd": 2460974.947532121
            # "notes": "Mars osculating elements derived from JPL/astropy at epoch (TDB) = 2460974.947532 JD"
    } # Это реальные данные взятые из get_mars_orb.py, там они брались из astroproxy


    print(res)

    return res


# --- Conversion: ecliptic heliocentric -> ICRS heliocentric (approx using mean obliquity) ---
EPSILON_DEG = 23.439291111
def ecliptic_helio_to_icrs_helio(r_ecl):
    eps = math.radians(EPSILON_DEG)
    cose = math.cos(eps); sine = math.sin(eps)
    x, y, z = r_ecl
    x_eq = x
    y_eq = cose * y - sine * z
    z_eq = sine * y + cose * z
    return np.array([x_eq, y_eq, z_eq])

# --- Astropy-based heliocentric vectors (Earth and Mars) ---
def earth_helio_icrs_from_jd(jd):
    if not ASTROPY_AVAILABLE:
        raise RuntimeError("Astropy is required for precise Earth positions.")
    try:
        solar_system_ephemeris.set('de440')
    except Exception:
        try:
            solar_system_ephemeris.set('de432s')
        except Exception:
            solar_system_ephemeris.set('builtin')
    t = Time(jd, format='jd', scale='tdb')
    earth_bary = get_body_barycentric('earth', t)
    sun_bary = get_body_barycentric('sun', t)
    earth_bary_xyz = earth_bary.xyz.to(u.AU).value.reshape(3,)
    sun_bary_xyz = sun_bary.xyz.to(u.AU).value.reshape(3,)
    return earth_bary_xyz - sun_bary_xyz  # AU, heliocentric ICRS

def mars_helio_icrs_from_jd(jd):
    """
    Uses astropy/JPL to return Mars heliocentric vector (ICRS) in AU.
    Fallback: if astropy unavailable, return analytical Kepler result converted to ICRS.
    """
    if ASTROPY_AVAILABLE:
        try:
            t = Time(jd, format='jd', scale='tdb')
            mars_bary = get_body_barycentric('mars', t)
            sun_bary = get_body_barycentric('sun', t)
            mars_bary_xyz = mars_bary.xyz.to(u.AU).value.reshape(3,)
            sun_bary_xyz = sun_bary.xyz.to(u.AU).value.reshape(3,)
            return mars_bary_xyz - sun_bary_xyz  # AU
        except Exception:
            # fallback to analytical model below
            pass

    # Fallback analytic: requires elements; we will use stub elements to produce a heliocentric ICRS vector
    # (this branch used only if astropy unavailable)
    stub = mechanism1_mars_stub()
    a = stub['a']; e = stub['e']; i_deg = stub['i']; Omega = stub['Omega']; omega = stub['omega']
    Tp_jd = rfc3339_to_jd(stub['Tp'])
    r_ecl = position_from_elements(a, e, i_deg, Omega, omega, Tp_jd, jd)  # AU, ecliptic
    r_icrs = ecliptic_helio_to_icrs_helio(r_ecl)
    return r_icrs

# --- Механизм №2 (основная функция) ---
def mechanism2_closest_approach_to_earth(elements, search_years=None, coarse_N=4000):
    a = float(elements['a']); e = float(elements['e'])
    i_deg = float(elements['i']); Omega = float(elements['Omega']); omega = float(elements['omega'])
    Tp_rfc = elements['Tp']
    Tp_jd = rfc3339_to_jd(Tp_rfc)
    P_days = 2*math.pi / math.sqrt(mu / (a**3))

    if search_years is None:
        # span_days = max(2.0 * P_days, 365.25*5.0)
        span_days = P_days + 365.25
    else:
        span_days = abs(search_years) * 365.25

    now_dt = datetime.datetime.utcnow()
    now_jd = (now_dt - datetime.datetime(1970,1,1)).total_seconds() / 86400.0 + 2440587.5
    epsilon_days = 1e-6
    t_start = now_jd + epsilon_days
    t_end = now_jd + span_days

    AU_km = 149597870.7

    ts = np.linspace(t_start, t_end, coarse_N)
    dists = np.empty_like(ts)
    for idx, t in enumerate(ts):
        # Mars position: use JPL (astropy) if available, fallback to analytic
        r_obj_icrs = mars_helio_icrs_from_jd(t)  # AU, heliocentric ICRS
        # Earth position: astropy (expected available); if exception, fallback to circular approx
        try:
            r_earth_icrs = earth_helio_icrs_from_jd(t)
        except Exception:
            theta = 2*math.pi * ((t - 2451545.0) / 365.25)
            r_earth_icrs = np.array([math.cos(theta), math.sin(theta), 0.0])
        dists[idx] = np.linalg.norm(r_obj_icrs - r_earth_icrs)

    min_idx = int(np.argmin(dists))
    left_idx = max(0, min_idx - 5)
    right_idx = min(len(ts)-1, min_idx + 5)
    t_left = ts[left_idx]; t_right = ts[right_idx]

    from scipy.optimize import minimize_scalar
    def dist_at_t(t):
        r_obj_icrs = mars_helio_icrs_from_jd(t)
        try:
            r_earth_icrs = earth_helio_icrs_from_jd(t)
        except Exception:
            theta = 2*math.pi * ((t - 2451545.0) / 365.25)
            r_earth_icrs = np.array([math.cos(theta), math.sin(theta), 0.0])
        return np.linalg.norm(r_obj_icrs - r_earth_icrs)

    res = minimize_scalar(dist_at_t, method='brent', bracket=(t_left, ts[min_idx], t_right))
    t_min = res.x if res.success else ts[min_idx]

    # ensure future
    if t_min <= now_jd:
        if min_idx+1 < len(ts):
            try:
                next_idx = min_idx+1 + int(np.argmin(dists[min_idx+1:]))
                left_idx = max(0, next_idx-5); right_idx = min(len(ts)-1, next_idx+5)
                t_left = ts[left_idx]; t_right = ts[right_idx]
                res2 = minimize_scalar(dist_at_t, method='brent', bracket=(t_left, ts[next_idx], t_right))
                if res2.success and res2.x > now_jd:
                    t_min = res2.x
            except Exception:
                pass

    d_min_AU = dist_at_t(t_min)
    d_min_km = d_min_AU * AU_km
    result = {
        't_closest_rfc3339': jd_to_rfc3339(t_min),
        't_jd': float(t_min),
        'distance_AU': float(d_min_AU),
        'distance_thousand_km': float(d_min_km / 1000.0),
        'search_span_days': float(span_days),
        'period_days': float(P_days),
        'search_now_rfc3339': now_dt.strftime("%Y-%m-%dT%H:%M:%SZ")
    }
    return result




# --- Демонстрация ---
if __name__ == "__main__":
    mars = mechanism1_mars_stub()
    print("Mars stub elements:")
    for k,v in mars.items():
        print(f"  {k}: {v}")
    print()

    if ASTROPY_AVAILABLE:
        print("Astropy available: using JPL ephemerides (DE).")
    else:
        print("Astropy NOT available: falling back to analytic approximations for Mars/Earth.")

    res = mechanism2_closest_approach_to_earth(mars, coarse_N=3000)
    print("Closest approach (JPL/astropy if available):")
    for k,v in res.items():
        if isinstance(v, float):
            print(f"  {k}: {v:.6f}")
        else:
            print(f"  {k}: {v}")
