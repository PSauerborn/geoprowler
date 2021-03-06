--
-- PostgreSQL database dump
--

-- Dumped from database version 12.4 (Debian 12.4-1.pgdg100+1)
-- Dumped by pg_dump version 12.4 (Debian 12.4-1.pgdg100+1)

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
-- Name: entities; Type: TABLE; Schema: public; Owner: geoprowler
--

CREATE TABLE public.entities (
    entity_id uuid NOT NULL,
    meta json DEFAULT '{}'::json
);


ALTER TABLE public.entities OWNER TO geoprowler;

--
-- Name: locations; Type: TABLE; Schema: public; Owner: geoprowler
--

CREATE TABLE public.locations (
    entity_id uuid NOT NULL,
    location json,
    last_updated timestamp without time zone
);


ALTER TABLE public.locations OWNER TO geoprowler;

--
-- Name: entities entities_pkey; Type: CONSTRAINT; Schema: public; Owner: geoprowler
--

ALTER TABLE ONLY public.entities
    ADD CONSTRAINT entities_pkey PRIMARY KEY (entity_id);


--
-- Name: locations locations_pkey; Type: CONSTRAINT; Schema: public; Owner: geoprowler
--

ALTER TABLE ONLY public.locations
    ADD CONSTRAINT locations_pkey PRIMARY KEY (entity_id);


--
-- PostgreSQL database dump complete
--

