-- Схема таблиц
CREATE SCHEMA public;

-- Таблица Anime
CREATE TABLE public."Anime" (
    id serial PRIMARY KEY,
    name text NOT NULL UNIQUE,
    img text NOT NULL,
    "imgHeader" text NOT NULL,
    describe text NOT NULL,
    genres text[],
    author text,
    country text NOT NULL,
    published integer NOT NULL,
    "averageRating" double precision DEFAULT 0,
    "ratingCount" integer DEFAULT 0,
    status text NOT NULL,
    popularity integer DEFAULT 0
);

-- Таблица Chapter
CREATE TABLE public."Chapter" (
    chapter integer NOT NULL,
    img text[],
    name text NOT NULL UNIQUE,
    "animeName" text NOT NULL REFERENCES public."Anime"(name) ON UPDATE CASCADE ON DELETE RESTRICT,
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Таблица User
CREATE TABLE public."User" (
    id text PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    email text NOT NULL UNIQUE,
    image text NOT NULL,
    favorite text[],
    "createdAt" timestamp(3) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Таблица _prisma_migrations
CREATE TABLE public._prisma_migrations (
    id character varying(36) PRIMARY KEY,
    checksum character varying(64) NOT NULL,
    finished_at timestamp with time zone,
    migration_name character varying(255) NOT NULL,
    logs text,
    rolled_back_at timestamp with time zone,
    started_at timestamp with time zone DEFAULT now() NOT NULL,
    applied_steps_count integer DEFAULT 0 NOT NULL
);