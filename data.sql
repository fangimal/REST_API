DROP TABLE IF EXISTS author CASCADE;
DROP TABLE IF EXISTS book CASCADE;
DROP TABLE IF EXISTS book_authors CASCADE;

CREATE TABLE public.author(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL
);

CREATE TABLE public.book(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    age INT
);

CREATE TABLE public.book_authors(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID NOT NULL,
    author_id UUID NOT NULL,

    CONSTRAINT book_fk FOREIGN KEY (book_id) REFERENCES public.book(id),
    CONSTRAINT author_fk FOREIGN KEY (author_id) REFERENCES public.author(id),
    CONSTRAINT book_author_unique UNIQUE (book_id, author_id)
);

INSERT INTO author(id, name) VALUES ('efebc31d-990c-4b15-8c85-189e063a5a1f','Народ');
INSERT INTO author(id, name) VALUES ('046b5565-9904-4f0c-9163-2b051947e52d','Джон Роулинг');
INSERT INTO author(id, name) VALUES ('e09948e6-57c8-4d44-a84c-d20c51ebec85','Джек Лондон');

INSERT INTO book(id, name, age) VALUES ('ecd3f1f8-6272-4c95-9375-9af0b1bce063','колобок', 1000);
INSERT INTO book(id, name, age) VALUES ('71872f9c-fdb7-46d8-8248-a24532058c0a','гарри поттер', 22);
INSERT INTO book(id, name) VALUES ('25950564-ff01-4188-bcf4-12f3cecc35da','бриллианты');

-- колобок
INSERT INTO book_authors(book_id, author_id) VALUES ('ecd3f1f8-6272-4c95-9375-9af0b1bce063', 'efebc31d-990c-4b15-8c85-189e063a5a1f');
INSERT INTO book_authors(book_id, author_id) VALUES ('ecd3f1f8-6272-4c95-9375-9af0b1bce063', '046b5565-9904-4f0c-9163-2b051947e52d');

-- HP
INSERT INTO book_authors(book_id, author_id) VALUES ('71872f9c-fdb7-46d8-8248-a24532058c0a', 'efebc31d-990c-4b15-8c85-189e063a5a1f');
INSERT INTO book_authors(book_id, author_id) VALUES ('71872f9c-fdb7-46d8-8248-a24532058c0a', 'e09948e6-57c8-4d44-a84c-d20c51ebec85');
INSERT INTO book_authors(book_id, author_id) VALUES ('71872f9c-fdb7-46d8-8248-a24532058c0a', '046b5565-9904-4f0c-9163-2b051947e52d');

SELECT
    b.id, b.name, ARRAY ((
    SELECT ba.author_id
    FROM book_authors ba
    WHERE ba.book_id = b.id)
                  ) AS authors
FROM book b
GROUP BY b.id, b.name;

SELECT
    a.id, a.name
FROM book_authors
JOIN public.author a ON a.id = book_authors.author_id
WHERE book_id = 'ecd3f1f8-6272-4c95-9375-9af0b1bce063';

SELECT
    *, (SELECT count(*) FROM book_authors where book_id = b.id) AS authors_count
FROM book b

