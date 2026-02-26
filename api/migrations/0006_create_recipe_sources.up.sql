CREATE TABLE recipe_sources (
    recipe_id UUID PRIMARY KEY REFERENCES recipes(id) ON DELETE CASCADE,    
    type INT NOT NULL,
    url TEXT,
    book_title TEXT,
    book_page INT,
    instructions TEXT,

    CONSTRAINT valid_type CHECK (type IN (1, 2, 3)),

    CONSTRAINT one_source_field CHECK (
        (type = 1 AND url IS NOT NULL AND book_title IS NULL AND book_page IS NULL AND instructions IS NULL) OR
        (type = 2 AND url IS NULL AND book_title IS NOT NULL AND book_page IS NOT NULL AND instructions IS NULL) OR
        (type = 3 AND url IS NULL AND book_title IS NULL AND book_page IS NULL AND instructions IS NOT NULL)
    )
);