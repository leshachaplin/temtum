CREATE TABLE public.block
(
    Index        CHARACTER VARYING(32)  NOT NULL,
    BeaconIndex  int                    NOT NULL,
    BeaconValue  CHARACTER VARYING(128) NOT NULL,
    Hash         CHARACTER VARYING(128) PRIMARY KEY,
    PreviousHash CHARACTER VARYING(128) NOT NULL,
    Time         int                    NOT NULL
);

