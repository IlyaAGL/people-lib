CREATE TABLE genders (
    id SERIAL PRIMARY KEY,
    gender VARCHAR(10) UNIQUE NOT NULL
);

CREATE TABLE nationalities (
    id SERIAL PRIMARY KEY,
    nationality VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE people (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    patronymic VARCHAR(100),
    age INTEGER NOT NULL,
    gender_id INTEGER NOT NULL REFERENCES genders(id) ON DELETE CASCADE,
    nationality_id INTEGER NOT NULL REFERENCES nationalities(id) ON DELETE CASCADE
);

CREATE INDEX idx_people_name ON people(name);
CREATE INDEX idx_people_surname ON people(surname);
CREATE INDEX idx_people_age ON people(age);
CREATE INDEX idx_people_gender_id ON people(gender_id);
CREATE INDEX idx_people_nationality_id ON people(nationality_id);
CREATE INDEX idx_genders_gender ON genders(gender);
CREATE INDEX idx_nationalities_nationality ON nationalities(nationality);