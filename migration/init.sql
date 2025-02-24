--
-- PostgreSQL database dump
--

-- Dumped from database version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.6 (Ubuntu 16.6-0ubuntu0.24.04.1)

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
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- PostgreSQL database dump complete
--




CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    name VARCHAR(60) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone
);

CREATE TABLE achievements_levels (
    id SERIAL PRIMARY KEY,
    achievement_id INT NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    level INT NOT NULL,
    description VARCHAR(100) NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone
);
CREATE TABLE users_achievements (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id INT NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    level INT NOT NULL DEFAULT 1,
    progress INT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, achievement_id),
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now(),
    deleted_at timestamp without time zone
);

ALTER TABLE achievements_levels ADD COLUMN to_complete INT NOT NULL;


INSERT INTO users_achievements (user_id, achievement_id, progress)
VALUES ('52d4abf7-a846-48d1-93f3-d651e41beda7', 1, 1)
ON CONFLICT (user_id, achievement_id) DO UPDATE
SET progress = users_achievements.progress + EXCLUDED.progress;


CREATE TABLE users_verification_codes (
    id SERIAL PRIMARY KEY,
	username varchar(50) NOT NULL,
    email VARCHAR(255) NOT NULL,
    verification_code VARCHAR(50) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_email ON users_verification_codes(email);
CREATE INDEX idx_expires_at ON users_verification_codes(expires_at);
