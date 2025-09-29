CREATE TABLE public.author(
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              name VARCHAR(100) NOT NULL
);

CREATE TABLE public.book(
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            name VARCHAR(100) NOT NULL,
                            author_id UUID NOT NULL,
                            CONSTRAINT author_fk FOREIGN KEY (author_id) REFERENCES public.author(id)
);

INSERT INTO author(name) VALUES ('Народ');  --efebc31d-990c-4b15-8c85-189e063a5a1f
INSERT INTO author(name) VALUES ('Джон Роулинг');  --046b5565-9904-4f0c-9163-2b051947e52d
INSERT INTO author(name) VALUES ('Джек Лондон');  --e09948e6-57c8-4d44-a84c-d20c51ebec85

INSERT INTO book(name, author_id) VALUES ('колобок', 'efebc31d-990c-4b15-8c85-189e063a5a1f');
INSERT INTO book(name, author_id) VALUES ('гарри поттер', '046b5565-9904-4f0c-9163-2b051947e52d');
INSERT INTO book(name, author_id) VALUES ('бриллианты', 'e09948e6-57c8-4d44-a84c-d20c51ebec85');