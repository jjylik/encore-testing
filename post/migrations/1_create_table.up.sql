
CREATE TABLE public.user (
    id varchar(255) PRIMARY KEY,
    enabled boolean DEFAULT TRUE,
    name text DEFAULT 'unknown'
);


CREATE TABLE public.post (
    id SERIAL PRIMARY KEY,
    user_id varchar(255) REFERENCES public.user (id),
    title text NOT NULL,
    content text,
    created_at timestamp with time zone DEFAULT now()
);
