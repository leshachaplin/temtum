CREATE TABLE public.status
(
    Address CHARACTER VARYING(128),
    Hash    CHARACTER VARYING(128) PRIMARY KEY ,
    Type    CHARACTER VARYING(32),
    Ban     boolean,
    Status  int
);