CREATE TABLE public.transaction
(
    Id             CHARACTER VARYING(128) PRIMARY KEY,
    BlockHash      CHARACTER VARYING(128) NOT NULL,
    TimeOfCreation int                    NOT NULL,
    Type           CHARACTER VARYING(128) NOT NULL,
    TxOuts         json                   NOT NULL,
    TxIns          json                   NOT NULL
);

ALTER TABLE public.transaction
    ADD CONSTRAINT block_transaction FOREIGN KEY (BlockHash) REFERENCES public.block (Hash) ON DELETE CASCADE;