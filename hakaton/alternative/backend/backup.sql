--
-- PostgreSQL database dump
--

\restrict YtV8EnHnkLGHrw0dXV1gwGLmTBYkFrebFWJEbdDXlUFoim2bAUNkEA8wkm9TJpR

-- Dumped from database version 15.14
-- Dumped by pg_dump version 15.14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: calculation_observations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.calculation_observations (
    orbital_calculation_id bigint NOT NULL,
    observation_id bigint NOT NULL
);


ALTER TABLE public.calculation_observations OWNER TO root;

--
-- Name: comets; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.comets (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    user_id bigint,
    name character varying(255),
    description text
);


ALTER TABLE public.comets OWNER TO root;

--
-- Name: comets_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.comets_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.comets_id_seq OWNER TO root;

--
-- Name: comets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.comets_id_seq OWNED BY public.comets.id;


--
-- Name: observations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.observations (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    comet_id bigint,
    observed_at timestamp with time zone,
    ra numeric,
    "dec" numeric,
    image_url text,
    notes text
);


ALTER TABLE public.observations OWNER TO root;

--
-- Name: observations_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.observations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.observations_id_seq OWNER TO root;

--
-- Name: observations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.observations_id_seq OWNED BY public.observations.id;


--
-- Name: orbital_calculations; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.orbital_calculations (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    comet_id bigint,
    calculated_at timestamp with time zone,
    semi_major_axis numeric,
    eccentricity numeric,
    inclination numeric,
    lon_ascending_node numeric,
    arg_periapsis numeric,
    time_perihelion timestamp with time zone,
    is_latest boolean,
    approach_date timestamp with time zone,
    distance_au numeric,
    distance_km numeric
);


ALTER TABLE public.orbital_calculations OWNER TO root;

--
-- Name: orbital_calculations_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.orbital_calculations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.orbital_calculations_id_seq OWNER TO root;

--
-- Name: orbital_calculations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.orbital_calculations_id_seq OWNED BY public.orbital_calculations.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.users (
    id bigint NOT NULL,
    email text NOT NULL,
    username text NOT NULL,
    avatar_url text,
    password text NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);


ALTER TABLE public.users OWNER TO root;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO root;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: comets id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.comets ALTER COLUMN id SET DEFAULT nextval('public.comets_id_seq'::regclass);


--
-- Name: observations id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.observations ALTER COLUMN id SET DEFAULT nextval('public.observations_id_seq'::regclass);


--
-- Name: orbital_calculations id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.orbital_calculations ALTER COLUMN id SET DEFAULT nextval('public.orbital_calculations_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: calculation_observations; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.calculation_observations (orbital_calculation_id, observation_id) FROM stdin;
\.


--
-- Data for Name: comets; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.comets (id, created_at, updated_at, user_id, name, description) FROM stdin;
1	2025-10-25 13:25:40.662152+00	2025-10-25 13:25:40.662152+00	1	aboba	very aboba
\.


--
-- Data for Name: observations; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.observations (id, created_at, updated_at, comet_id, observed_at, ra, "dec", image_url, notes) FROM stdin;
\.


--
-- Data for Name: orbital_calculations; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.orbital_calculations (id, created_at, updated_at, comet_id, calculated_at, semi_major_axis, eccentricity, inclination, lon_ascending_node, arg_periapsis, time_perihelion, is_latest, approach_date, distance_au, distance_km) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: root
--

COPY public.users (id, email, username, avatar_url, password, created_at, updated_at) FROM stdin;
1	ab@ba	aboba	\N	$2a$10$lYe/WmRZRdjM3P4SqC5oSejvo4MvgE7cLgq/rpyxUN/de3S0U69iq	2025-10-25 13:25:25.203004+00	2025-10-25 13:25:25.203004+00
2	aboba@1	aboba1	\N	$2a$10$2pXk13dqT.g5cZ.vif.q2efRp7X1j4q.1A9LxnLmODqQwSj7AB/Ou	2025-10-25 15:11:35.475699+00	2025-10-25 15:11:35.475699+00
3	user@example.com	SuperUser	\N	$2a$10$GnhBOEWliVD8aRRjDE9W6uYzrhucgMJ/CVIpA9KzsgGuKhwCSvpOm	2025-10-25 15:28:31.635686+00	2025-10-25 15:28:31.635686+00
\.


--
-- Name: comets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.comets_id_seq', 1, true);


--
-- Name: observations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.observations_id_seq', 1, false);


--
-- Name: orbital_calculations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.orbital_calculations_id_seq', 1, false);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.users_id_seq', 3, true);


--
-- Name: calculation_observations calculation_observations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.calculation_observations
    ADD CONSTRAINT calculation_observations_pkey PRIMARY KEY (orbital_calculation_id, observation_id);


--
-- Name: comets comets_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.comets
    ADD CONSTRAINT comets_pkey PRIMARY KEY (id);


--
-- Name: observations observations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.observations
    ADD CONSTRAINT observations_pkey PRIMARY KEY (id);


--
-- Name: orbital_calculations orbital_calculations_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.orbital_calculations
    ADD CONSTRAINT orbital_calculations_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_comets_user_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_comets_user_id ON public.comets USING btree (user_id);


--
-- Name: idx_observations_comet_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_observations_comet_id ON public.observations USING btree (comet_id);


--
-- Name: idx_observations_observed_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_observations_observed_at ON public.observations USING btree (observed_at);


--
-- Name: idx_orbital_calculations_comet_id; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_orbital_calculations_comet_id ON public.orbital_calculations USING btree (comet_id);


--
-- Name: idx_orbital_calculations_is_latest; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_orbital_calculations_is_latest ON public.orbital_calculations USING btree (is_latest);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: calculation_observations fk_calculation_observations_observation; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.calculation_observations
    ADD CONSTRAINT fk_calculation_observations_observation FOREIGN KEY (observation_id) REFERENCES public.observations(id);


--
-- Name: calculation_observations fk_calculation_observations_orbital_calculation; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.calculation_observations
    ADD CONSTRAINT fk_calculation_observations_orbital_calculation FOREIGN KEY (orbital_calculation_id) REFERENCES public.orbital_calculations(id);


--
-- Name: observations fk_comets_observations; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.observations
    ADD CONSTRAINT fk_comets_observations FOREIGN KEY (comet_id) REFERENCES public.comets(id);


--
-- Name: orbital_calculations fk_comets_orbital_calculations; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.orbital_calculations
    ADD CONSTRAINT fk_comets_orbital_calculations FOREIGN KEY (comet_id) REFERENCES public.comets(id);


--
-- PostgreSQL database dump complete
--

\unrestrict YtV8EnHnkLGHrw0dXV1gwGLmTBYkFrebFWJEbdDXlUFoim2bAUNkEA8wkm9TJpR

