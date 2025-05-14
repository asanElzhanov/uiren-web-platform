--
-- PostgreSQL database dump
--

-- Dumped from database version 16.8 (Ubuntu 16.8-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.8 (Ubuntu 16.8-0ubuntu0.24.04.1)

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
-- Name: achievements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.achievements (
    id integer NOT NULL,
    name character varying(60) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone
);


ALTER TABLE public.achievements OWNER TO postgres;

--
-- Name: achievements_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.achievements_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.achievements_id_seq OWNER TO postgres;

--
-- Name: achievements_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.achievements_id_seq OWNED BY public.achievements.id;


--
-- Name: achievements_levels; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.achievements_levels (
    achievement_id integer NOT NULL,
    level integer NOT NULL,
    description character varying(100) NOT NULL,
    threshold integer NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone,
    CONSTRAINT achievements_levels_level_check CHECK ((level > 0)),
    CONSTRAINT achievements_levels_threshold_check CHECK ((threshold >= 0))
);


ALTER TABLE public.achievements_levels OWNER TO postgres;

--
-- Name: badges; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.badges (
    badge text NOT NULL,
    description text NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.badges OWNER TO postgres;

--
-- Name: friendships; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.friendships (
    id integer NOT NULL,
    user1_username character varying(50) NOT NULL,
    user2_username character varying(50) NOT NULL,
    status character varying(20) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    recipient character varying(50) DEFAULT ''::character varying NOT NULL,
    CONSTRAINT check_username_order CHECK (((user1_username)::text < (user2_username)::text)),
    CONSTRAINT friendships_status_check CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'accepted'::character varying, 'declined'::character varying])::text[])))
);


ALTER TABLE public.friendships OWNER TO postgres;

--
-- Name: friendships_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.friendships_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.friendships_id_seq OWNER TO postgres;

--
-- Name: friendships_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.friendships_id_seq OWNED BY public.friendships.id;


--
-- Name: users_achievements; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_achievements (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    achievement_id integer NOT NULL,
    achievement_level integer DEFAULT 1 NOT NULL,
    progress integer DEFAULT 0 NOT NULL,
    last_updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT user_achievements_achievement_level_check CHECK ((achievement_level >= 1)),
    CONSTRAINT user_achievements_progress_check CHECK ((progress >= 0))
);


ALTER TABLE public.users_achievements OWNER TO postgres;

--
-- Name: user_achievements_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.user_achievements_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.user_achievements_id_seq OWNER TO postgres;

--
-- Name: user_achievements_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.user_achievements_id_seq OWNED BY public.users_achievements.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(255) NOT NULL,
    password character varying(60) NOT NULL,
    first_name character varying(50),
    last_name character varying(50),
    phone character varying(20),
    is_active boolean DEFAULT false,
    is_admin boolean DEFAULT false,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_badges; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_badges (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    badge character varying(50) NOT NULL,
    awarded_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users_badges OWNER TO postgres;

--
-- Name: users_badges_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_badges_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_badges_id_seq OWNER TO postgres;

--
-- Name: users_badges_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_badges_id_seq OWNED BY public.users_badges.id;


--
-- Name: users_progress; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_progress (
    id integer NOT NULL,
    user_id uuid NOT NULL,
    xp integer DEFAULT 0 NOT NULL,
    last_updated timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_progress_xp_check CHECK ((xp >= 0))
);


ALTER TABLE public.users_progress OWNER TO postgres;

--
-- Name: users_progress_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_progress_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_progress_id_seq OWNER TO postgres;

--
-- Name: users_progress_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_progress_id_seq OWNED BY public.users_progress.id;


--
-- Name: users_verification_codes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users_verification_codes (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    email character varying(255) NOT NULL,
    verification_code character varying(50) NOT NULL,
    expires_at timestamp without time zone NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.users_verification_codes OWNER TO postgres;

--
-- Name: users_verification_codes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_verification_codes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_verification_codes_id_seq OWNER TO postgres;

--
-- Name: users_verification_codes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_verification_codes_id_seq OWNED BY public.users_verification_codes.id;


--
-- Name: achievements id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.achievements ALTER COLUMN id SET DEFAULT nextval('public.achievements_id_seq'::regclass);


--
-- Name: friendships id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.friendships ALTER COLUMN id SET DEFAULT nextval('public.friendships_id_seq'::regclass);


--
-- Name: users_achievements id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_achievements ALTER COLUMN id SET DEFAULT nextval('public.user_achievements_id_seq'::regclass);


--
-- Name: users_badges id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_badges ALTER COLUMN id SET DEFAULT nextval('public.users_badges_id_seq'::regclass);


--
-- Name: users_progress id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_progress ALTER COLUMN id SET DEFAULT nextval('public.users_progress_id_seq'::regclass);


--
-- Name: users_verification_codes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_verification_codes ALTER COLUMN id SET DEFAULT nextval('public.users_verification_codes_id_seq'::regclass);


--
-- Data for Name: achievements; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.achievements (id, name, created_at, updated_at, deleted_at) FROM stdin;
3	asn	2025-03-05 01:37:26.392046	2025-03-05 11:41:14.346905	\N
4	testAch	2025-03-05 10:40:13.806747	2025-03-05 10:40:13.806747	2025-03-05 22:54:53.930766
5	test_asan	2025-03-05 22:55:31.189625	2025-03-05 22:55:31.189625	2025-03-05 22:57:45.289659
1	asan	2025-03-05 01:37:17.432851	2025-03-05 11:23:16.258867	\N
\.


--
-- Data for Name: achievements_levels; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.achievements_levels (achievement_id, level, description, threshold, created_at, updated_at, deleted_at) FROM stdin;
3	1	earn 62	62	2025-03-05 22:21:34.891941	2025-03-05 22:21:34.891941	\N
3	2	en 62	63	2025-05-10 15:55:11.733415	2025-05-10 15:55:11.733415	\N
3	3	en 62	64	2025-05-10 15:55:14.971756	2025-05-10 15:55:14.971756	\N
1	1	earn	100	2025-05-13 15:43:29.894828	2025-05-13 15:43:29.894828	\N
1	2	earn	200	2025-05-13 15:43:37.190391	2025-05-13 15:43:37.190391	\N
1	3	earn	300	2025-05-13 15:43:43.000507	2025-05-13 15:43:43.000507	\N
1	4	earn	400	2025-05-13 15:43:48.533489	2025-05-13 15:43:48.533489	\N
1	5	earn	599	2025-05-13 15:43:53.675232	2025-05-13 15:43:53.675232	\N
\.


--
-- Data for Name: badges; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.badges (badge, description, created_at) FROM stdin;
new_badges	dec	2025-05-13 15:02:09.791508
badge1	badge1	2025-05-13 15:45:42.232317
badge2	badge2	2025-05-13 15:45:42.232317
badge3	badge3	2025-05-13 15:45:42.232317
badge4	badge4	2025-05-13 15:45:42.232317
badge5	badge5	2025-05-13 15:45:42.232317
badge6	badge6	2025-05-13 15:45:42.232317
badge7	badge7	2025-05-13 15:45:42.232317
badge8	badge8	2025-05-13 15:45:42.232317
badge9	badge9	2025-05-13 15:45:42.232317
badge10	badge10	2025-05-13 15:45:42.232317
badge11	badge11	2025-05-13 15:45:42.232317
badge12	badge12	2025-05-13 15:45:42.232317
badge13	badge13	2025-05-13 15:45:42.232317
badge14	badge14	2025-05-13 15:45:42.232317
badge15	badge15	2025-05-13 15:45:42.232317
badge16	badge16	2025-05-13 15:45:42.232317
badge17	badge17	2025-05-13 15:45:42.232317
badge18	badge18	2025-05-13 15:45:42.232317
badge19	badge19	2025-05-13 15:45:42.232317
badge20	badge20	2025-05-13 15:45:42.232317
badge21	badge21	2025-05-13 15:45:42.232317
badge22	badge22	2025-05-13 15:45:42.232317
badge23	badge23	2025-05-13 15:45:42.232317
badge24	badge24	2025-05-13 15:45:42.232317
badge25	badge25	2025-05-13 15:45:42.232317
badge26	badge26	2025-05-13 15:45:42.232317
badge27	badge27	2025-05-13 15:45:42.232317
badge28	badge28	2025-05-13 15:45:42.232317
badge29	badge29	2025-05-13 15:45:42.232317
badge30	badge30	2025-05-13 15:45:42.232317
badge31	badge31	2025-05-13 15:45:42.232317
badge32	badge32	2025-05-13 15:45:42.232317
badge33	badge33	2025-05-13 15:45:42.232317
badge34	badge34	2025-05-13 15:45:42.232317
badge35	badge35	2025-05-13 15:45:42.232317
badge36	badge36	2025-05-13 15:45:42.232317
badge37	badge37	2025-05-13 15:45:42.232317
badge38	badge38	2025-05-13 15:45:42.232317
badge39	badge39	2025-05-13 15:45:42.232317
badge40	badge40	2025-05-13 15:45:42.232317
badge41	badge41	2025-05-13 15:45:42.232317
badge42	badge42	2025-05-13 15:45:42.232317
badge43	badge43	2025-05-13 15:45:42.232317
badge44	badge44	2025-05-13 15:45:42.232317
badge45	badge45	2025-05-13 15:45:42.232317
badge46	badge46	2025-05-13 15:45:42.232317
badge47	badge47	2025-05-13 15:45:42.232317
badge48	badge48	2025-05-13 15:45:42.232317
badge49	badge49	2025-05-13 15:45:42.232317
badge50	badge50	2025-05-13 15:45:42.232317
\.


--
-- Data for Name: friendships; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.friendships (id, user1_username, user2_username, status, created_at, recipient) FROM stdin;
66	asan1	s4ab	accepted	2025-05-14 00:50:25.151358	asan1
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, username, email, password, first_name, last_name, phone, is_active, is_admin, created_at, updated_at, deleted_at) FROM stdin;
9b5b565e-7ee5-4f5c-abe7-cc047fe60f0b	asan3	asn@gmail.com	$2a$10$QQl6Hp0STOm2aPCMBsR32.t.nMRDKFKrSMbPE9sG9La8ZdChEP2DC	\N	\N	\N	f	f	2025-02-16 01:36:00.471728	2025-02-16 01:36:00.471728	\N
ec55e15c-b7f1-4ac1-8b46-fca647c2e2ca	niggaaф	xlm9877@gmail.com	$2a$10$SizhcHUeDAW6dDw3k.lX/.r2jPtD6xLhtCX.WmlnBvnK0BZg48s2G	\N	\N	\N	f	f	2025-02-23 15:51:26.761191	2025-02-23 15:51:26.761191	\N
52d4abf7-a846-48d1-93f3-d651e41beda7	s4ab	some_test@gmail.com	$2a$10$9y1Js.sdzxRgeKW3lyOI..z5tjc8t3RM8KQxNV.o1EH9TPPWuYI0y			87474434210	t	t	2025-02-16 15:06:55.433497	2025-02-16 16:35:14.410738	\N
90ba1dc7-e68f-4427-9298-cd57f9b0bd0d	niggaa	asanelzhanov2@gmail.com	$2a$10$5tIIwikxyCljy1RCQb/09.8DP9GADcZig7zaK2EQWge8oTTeDJ6Ui	\N	\N	\N	t	t	2025-02-22 14:49:16.575542	2025-02-22 14:49:16.575542	\N
41c3beaa-a9c8-4db9-a0bb-9b447ecda079	seab_test1	seab@seb.kz	$2a$10$qN/Fv/wGeTJm.0Tlz/X6RuAekhyWvo4U0R17Wv.hmCeECX6t6uZ8m	\N	\N	\N	f	f	2025-05-10 20:26:34.93027	2025-05-10 20:26:34.93027	\N
c1eeb883-0726-45ac-84a1-040525662e3a	asan1	fdasfadsf	$2a$10$9y1Js.sdzxRgeKW3lyOI..z5tjc8t3RM8KQxNV.o1EH9TPPWuYI0y	fdasfas	fdasfadsfff	fdas	t	f	2025-02-15 21:55:29.213846	2025-02-16 10:53:00.40047	\N
cceb89db-81ce-44a3-82d7-40aa68e02a6b	asan2	asaan@gmail.com	$2a$10$9y1Js.sdzxRgeKW3lyOI..z5tjc8t3RM8KQxNV.o1EH9TPPWuYI0y	\N	\N	\N	t	f	2025-02-15 22:31:34.459659	2025-02-15 22:31:34.459659	\N
ca820365-be55-4a39-b23b-098ca0e6dab5	seab_test	seab@seab.kz	$2a$10$9y1Js.sdzxRgeKW3lyOI..z5tjc8t3RM8KQxNV.o1EH9TPPWuYI0y	\N	\N	\N	t	f	2025-05-10 20:25:54.902171	2025-05-10 20:25:54.902171	\N
\.


--
-- Data for Name: users_achievements; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users_achievements (id, user_id, achievement_id, achievement_level, progress, last_updated) FROM stdin;
46	c1eeb883-0726-45ac-84a1-040525662e3a	1	5	500	2025-05-14 00:08:12.668301
51	cceb89db-81ce-44a3-82d7-40aa68e02a6b	1	2	100	2025-05-14 15:12:39.21107
52	ca820365-be55-4a39-b23b-098ca0e6dab5	1	3	200	2025-05-14 17:10:46.97513
\.


--
-- Data for Name: users_badges; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users_badges (id, user_id, badge, awarded_at) FROM stdin;
67	c1eeb883-0726-45ac-84a1-040525662e3a	badge14	2025-05-14 00:06:11.440867
68	c1eeb883-0726-45ac-84a1-040525662e3a	badge1	2025-05-14 00:06:47.415158
70	c1eeb883-0726-45ac-84a1-040525662e3a	badge21	2025-05-14 00:08:03.683645
71	c1eeb883-0726-45ac-84a1-040525662e3a	badge22	2025-05-14 00:08:12.668301
73	c1eeb883-0726-45ac-84a1-040525662e3a	badge3	2025-05-14 15:09:18.395771
75	c1eeb883-0726-45ac-84a1-040525662e3a	badge4	2025-05-14 15:10:25.968613
76	c1eeb883-0726-45ac-84a1-040525662e3a	badge7	2025-05-14 15:11:31.478016
78	c1eeb883-0726-45ac-84a1-040525662e3a	badge8	2025-05-14 15:12:16.158123
79	cceb89db-81ce-44a3-82d7-40aa68e02a6b	badge8	2025-05-14 15:12:39.21107
80	ca820365-be55-4a39-b23b-098ca0e6dab5	badge8	2025-05-14 17:10:46.97513
81	ca820365-be55-4a39-b23b-098ca0e6dab5	badge1	2025-05-14 17:10:57.403587
43	52d4abf7-a846-48d1-93f3-d651e41beda7	badge1	2025-05-13 15:45:53.857406
45	52d4abf7-a846-48d1-93f3-d651e41beda7	badge2	2025-05-13 15:46:23.839568
47	52d4abf7-a846-48d1-93f3-d651e41beda7	badge3	2025-05-13 15:46:34.726191
50	52d4abf7-a846-48d1-93f3-d651e41beda7	badge4	2025-05-13 15:49:34.134384
52	52d4abf7-a846-48d1-93f3-d651e41beda7	badge5	2025-05-13 15:49:53.818934
53	52d4abf7-a846-48d1-93f3-d651e41beda7	badge6	2025-05-13 15:50:06.070432
54	52d4abf7-a846-48d1-93f3-d651e41beda7	badge7	2025-05-13 15:50:17.280572
55	52d4abf7-a846-48d1-93f3-d651e41beda7	badge8	2025-05-13 15:50:23.592481
56	52d4abf7-a846-48d1-93f3-d651e41beda7	badge9	2025-05-13 15:50:29.277147
61	52d4abf7-a846-48d1-93f3-d651e41beda7	badge11	2025-05-13 22:59:24.189663
63	52d4abf7-a846-48d1-93f3-d651e41beda7	badge12	2025-05-13 23:02:54.315621
65	52d4abf7-a846-48d1-93f3-d651e41beda7	badge14	2025-05-13 23:03:22.393012
\.


--
-- Data for Name: users_progress; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users_progress (id, user_id, xp, last_updated) FROM stdin;
1	90ba1dc7-e68f-4427-9298-cd57f9b0bd0d	48	2025-05-10 15:59:18.741826
8	52d4abf7-a846-48d1-93f3-d651e41beda7	3280	2025-05-12 19:52:40.247303
40	c1eeb883-0726-45ac-84a1-040525662e3a	13480	2025-05-14 00:06:11.440867
48	cceb89db-81ce-44a3-82d7-40aa68e02a6b	3280	2025-05-14 15:12:39.21107
49	ca820365-be55-4a39-b23b-098ca0e6dab5	321434521	2025-05-14 17:10:46.97513
\.


--
-- Data for Name: users_verification_codes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users_verification_codes (id, username, email, verification_code, expires_at, created_at) FROM stdin;
10	niggaa	asanelzhanov2@gmail.com	BdNybUgVfo	2025-02-22 14:59:16.577209	2025-02-22 14:49:16.577671
11	niggaaф	xlm9877@gmail.com	EdgZsP6L88	2025-02-23 15:51:26.762562	2025-02-23 15:51:26.762905
12	seab_test	seab@seab.kz	ieo05v2aN4	2025-05-10 20:25:54.904226	2025-05-10 20:25:54.904656
13	seab_test1	seab@seb.kz	N0906FhGCS	2025-05-10 20:26:34.93134	2025-05-10 20:26:34.931683
\.


--
-- Name: achievements_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.achievements_id_seq', 5, true);


--
-- Name: friendships_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.friendships_id_seq', 67, true);


--
-- Name: user_achievements_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.user_achievements_id_seq', 53, true);


--
-- Name: users_badges_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_badges_id_seq', 81, true);


--
-- Name: users_progress_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_progress_id_seq', 50, true);


--
-- Name: users_verification_codes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_verification_codes_id_seq', 13, true);


--
-- Name: achievements_levels achievements_levels_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.achievements_levels
    ADD CONSTRAINT achievements_levels_pkey PRIMARY KEY (achievement_id, level);


--
-- Name: achievements achievements_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.achievements
    ADD CONSTRAINT achievements_name_key UNIQUE (name);


--
-- Name: achievements achievements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.achievements
    ADD CONSTRAINT achievements_pkey PRIMARY KEY (id);


--
-- Name: badges badges_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.badges
    ADD CONSTRAINT badges_pkey PRIMARY KEY (badge);


--
-- Name: friendships friendships_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.friendships
    ADD CONSTRAINT friendships_pkey PRIMARY KEY (id);


--
-- Name: friendships friendships_user1_username_user2_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.friendships
    ADD CONSTRAINT friendships_user1_username_user2_username_key UNIQUE (user1_username, user2_username);


--
-- Name: users_achievements unique_ach; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_achievements
    ADD CONSTRAINT unique_ach UNIQUE (user_id, achievement_id);


--
-- Name: users_progress unique_user_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_progress
    ADD CONSTRAINT unique_user_id UNIQUE (user_id);


--
-- Name: users_achievements user_achievements_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_achievements
    ADD CONSTRAINT user_achievements_pkey PRIMARY KEY (id);


--
-- Name: users_achievements user_achievements_user_id_achievement_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_achievements
    ADD CONSTRAINT user_achievements_user_id_achievement_id_key UNIQUE (user_id, achievement_id);


--
-- Name: users_badges users_badges_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_badges
    ADD CONSTRAINT users_badges_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_phone_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_key UNIQUE (phone);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users_progress users_progress_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_progress
    ADD CONSTRAINT users_progress_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: users_verification_codes users_verification_codes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_verification_codes
    ADD CONSTRAINT users_verification_codes_pkey PRIMARY KEY (id);


--
-- Name: idx_email; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_email ON public.users_verification_codes USING btree (email);


--
-- Name: idx_expires_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_expires_at ON public.users_verification_codes USING btree (expires_at);


--
-- Name: unique_user_badge; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX unique_user_badge ON public.users_badges USING btree (user_id, badge);


--
-- Name: achievements_levels achievements_levels_achievement_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.achievements_levels
    ADD CONSTRAINT achievements_levels_achievement_id_fkey FOREIGN KEY (achievement_id) REFERENCES public.achievements(id) ON DELETE CASCADE;


--
-- Name: users_badges fk_badge; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_badges
    ADD CONSTRAINT fk_badge FOREIGN KEY (badge) REFERENCES public.badges(badge) ON DELETE CASCADE;


--
-- Name: friendships friendships_user1_username_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.friendships
    ADD CONSTRAINT friendships_user1_username_fkey FOREIGN KEY (user1_username) REFERENCES public.users(username) ON DELETE CASCADE;


--
-- Name: friendships friendships_user2_username_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.friendships
    ADD CONSTRAINT friendships_user2_username_fkey FOREIGN KEY (user2_username) REFERENCES public.users(username) ON DELETE CASCADE;


--
-- Name: users_achievements user_achievements_achievement_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_achievements
    ADD CONSTRAINT user_achievements_achievement_id_fkey FOREIGN KEY (achievement_id) REFERENCES public.achievements(id) ON DELETE CASCADE;


--
-- Name: users_achievements user_achievements_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_achievements
    ADD CONSTRAINT user_achievements_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users_badges users_badges_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_badges
    ADD CONSTRAINT users_badges_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: users_progress users_progress_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users_progress
    ADD CONSTRAINT users_progress_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

