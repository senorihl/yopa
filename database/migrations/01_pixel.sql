CREATE TABLE IF NOT EXISTS site (
    id                INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY ,
    name              VARCHAR NOT NULL
);

INSERT INTO site(name) VALUES ('Default Test App');

CREATE TABLE IF NOT EXISTS pixel
(
    id                BIGINT GENERATED ALWAYS AS IDENTITY,
    site_id           INT REFERENCES site (id) NOT NULL,
    ts                TIMESTAMP WITH TIME ZONE NOT NULL,
    visitor           UUID                     NOT NULL,
    name              VARCHAR                  NOT NULL,

    page              VARCHAR,
    page_chapter1     VARCHAR,
    page_chapter2     VARCHAR,
    page_chapter3     VARCHAR,
    page_full         VARCHAR GENERATED ALWAYS AS (page_chapter1 || '::' || page_chapter2 || '::' || page_chapter3 || '::'|| page) STORED,

    action            VARCHAR,
    action_type       VARCHAR,
    action_chapter1   VARCHAR,
    action_chapter2   VARCHAR,
    action_chapter3   VARCHAR,
    action_full       VARCHAR GENERATED ALWAYS AS (action_chapter1 || '::' || action_chapter2 || '::' || action_chapter3 || '::' || action || '@' || action_type) STORED,

    custom_properties JSON,
    PRIMARY KEY (id, ts)
);

INSERT INTO migrations(name) VALUES ('01_pixel');