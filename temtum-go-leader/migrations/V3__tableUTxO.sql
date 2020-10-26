CREATE TABLE public.txOuts
(
    OutputIndex int PRIMARY KEY,
    OutputId    CHARACTER VARYING(128),
    Address     CHARACTER VARYING(128),
    Amount      integer
);