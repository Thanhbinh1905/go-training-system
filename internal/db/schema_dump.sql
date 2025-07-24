--
-- PostgreSQL database dump
--

-- Dumped from database version 13.21
-- Dumped by pg_dump version 13.21

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
-- Name: folder_shares; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.folder_shares (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    folder_id uuid NOT NULL,
    user_id uuid NOT NULL,
    access character varying(10) NOT NULL,
    shared_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    shared_by_id uuid NOT NULL,
    CONSTRAINT chk_folder_shares_access CHECK (((access)::text = ANY ((ARRAY['read'::character varying, 'write'::character varying])::text[])))
);


ALTER TABLE public.folder_shares OWNER TO root;

--
-- Name: folders; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.folders (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    description text,
    owner_id uuid NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.folders OWNER TO root;

--
-- Name: note_shares; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.note_shares (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    note_id uuid NOT NULL,
    user_id uuid NOT NULL,
    access character varying(10) NOT NULL,
    shared_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    shared_by_id uuid NOT NULL,
    CONSTRAINT chk_note_shares_access CHECK (((access)::text = ANY ((ARRAY['read'::character varying, 'write'::character varying])::text[])))
);


ALTER TABLE public.note_shares OWNER TO root;

--
-- Name: notes; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.notes (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title text NOT NULL,
    body text,
    folder_id uuid NOT NULL,
    owner_id uuid NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.notes OWNER TO root;

--
-- Name: team_managers; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.team_managers (
    team_id uuid NOT NULL,
    user_id uuid NOT NULL,
    added_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    added_by_id uuid
);


ALTER TABLE public.team_managers OWNER TO root;

--
-- Name: team_members; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.team_members (
    team_id uuid NOT NULL,
    user_id uuid NOT NULL,
    added_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    added_by_id uuid
);


ALTER TABLE public.team_members OWNER TO root;

--
-- Name: teams; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.teams (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    team_name text NOT NULL,
    created_by_id uuid NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.teams OWNER TO root;

--
-- Name: users; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username text NOT NULL,
    email text NOT NULL,
    password_hash text NOT NULL,
    role character varying(20) NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT chk_users_role CHECK (((role)::text = ANY ((ARRAY['manager'::character varying, 'member'::character varying])::text[])))
);


ALTER TABLE public.users OWNER TO root;

--
-- Name: folder_shares folder_shares_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.folder_shares
    ADD CONSTRAINT folder_shares_pkey PRIMARY KEY (id);


--
-- Name: folders folders_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.folders
    ADD CONSTRAINT folders_pkey PRIMARY KEY (id);


--
-- Name: note_shares note_shares_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.note_shares
    ADD CONSTRAINT note_shares_pkey PRIMARY KEY (id);


--
-- Name: notes notes_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.notes
    ADD CONSTRAINT notes_pkey PRIMARY KEY (id);


--
-- Name: team_managers team_managers_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.team_managers
    ADD CONSTRAINT team_managers_pkey PRIMARY KEY (team_id, user_id);


--
-- Name: team_members team_members_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.team_members
    ADD CONSTRAINT team_members_pkey PRIMARY KEY (team_id, user_id);


--
-- Name: teams teams_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.teams
    ADD CONSTRAINT teams_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_folders_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_folders_deleted_at ON public.folders USING btree (deleted_at);


--
-- Name: idx_notes_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_notes_deleted_at ON public.notes USING btree (deleted_at);


--
-- Name: idx_teams_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_teams_deleted_at ON public.teams USING btree (deleted_at);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: root
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: idx_users_email; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_users_email ON public.users USING btree (email);


--
-- Name: idx_users_username; Type: INDEX; Schema: public; Owner: root
--

CREATE UNIQUE INDEX idx_users_username ON public.users USING btree (username);


--
-- Name: folder_shares fk_folder_shares_shared_by; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.folder_shares
    ADD CONSTRAINT fk_folder_shares_shared_by FOREIGN KEY (shared_by_id) REFERENCES public.users(id);


--
-- Name: notes fk_folders_notes; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.notes
    ADD CONSTRAINT fk_folders_notes FOREIGN KEY (folder_id) REFERENCES public.folders(id);


--
-- Name: folder_shares fk_folders_shares; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.folder_shares
    ADD CONSTRAINT fk_folders_shares FOREIGN KEY (folder_id) REFERENCES public.folders(id);


--
-- Name: note_shares fk_note_shares_shared_by; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.note_shares
    ADD CONSTRAINT fk_note_shares_shared_by FOREIGN KEY (shared_by_id) REFERENCES public.users(id);


--
-- Name: note_shares fk_notes_shares; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.note_shares
    ADD CONSTRAINT fk_notes_shares FOREIGN KEY (note_id) REFERENCES public.notes(id);


--
-- Name: team_managers fk_team_managers_added_by; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.team_managers
    ADD CONSTRAINT fk_team_managers_added_by FOREIGN KEY (added_by_id) REFERENCES public.users(id);


--
-- Name: team_managers fk_team_managers_team; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.team_managers
    ADD CONSTRAINT fk_team_managers_team FOREIGN KEY (team_id) REFERENCES public.teams(id);


--
-- Name: team_managers fk_team_managers_user; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.team_managers
    ADD CONSTRAINT fk_team_managers_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: team_members fk_team_members_added_by; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.team_members
    ADD CONSTRAINT fk_team_members_added_by FOREIGN KEY (added_by_id) REFERENCES public.users(id);


--
-- Name: team_members fk_team_members_team; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.team_members
    ADD CONSTRAINT fk_team_members_team FOREIGN KEY (team_id) REFERENCES public.teams(id);


--
-- Name: team_members fk_team_members_user; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.team_members
    ADD CONSTRAINT fk_team_members_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: folder_shares fk_users_folder_shares; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.folder_shares
    ADD CONSTRAINT fk_users_folder_shares FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: note_shares fk_users_note_shares; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.note_shares
    ADD CONSTRAINT fk_users_note_shares FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: folders fk_users_owned_folders; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.folders
    ADD CONSTRAINT fk_users_owned_folders FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: notes fk_users_owned_notes; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.notes
    ADD CONSTRAINT fk_users_owned_notes FOREIGN KEY (owner_id) REFERENCES public.users(id);


--
-- Name: teams fk_users_owned_teams; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.teams
    ADD CONSTRAINT fk_users_owned_teams FOREIGN KEY (created_by_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

