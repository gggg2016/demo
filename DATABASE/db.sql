-- Table: public.users

-- DROP TABLE public.users;

CREATE TABLE users IF NOT EXISTS
(
    id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    created_time timestamp(6) with time zone NOT NULL DEFAULT now(),
    updated_time timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT users_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to postgres;


-- Table: public.relationships

-- DROP TABLE public.relationships;

CREATE TABLE relationships IF NOT EXISTS
(
    id bigint NOT NULL DEFAULT nextval('relationships_id_seq'::regclass),
    user_id bigint NOT NULL,
    other_user_id bigint NOT NULL,
    state "RelationState" NOT NULL,
    create_time timestamp with time zone NOT NULL DEFAULT now(),
    updated_time timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT relationships_pkey PRIMARY KEY (id),
    CONSTRAINT uniq_relationship UNIQUE (user_id, other_user_id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.relationships
    OWNER to postgres;