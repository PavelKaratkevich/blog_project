CREATE TABLE IF NOT EXISTS articles (
    id serial PRIMARY key not null,
    title varchar(100),
    anons varchar(255),
    full_text text
)