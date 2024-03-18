CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     login VARCHAR(255) NOT NULL,
                                     password BYTEA NOT NULL,
                                     role BOOL,
                                     CONSTRAINT unique_login UNIQUE (login)
);

CREATE TABLE IF NOT EXISTS actor (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     gender VARCHAR(10) CHECK (gender IN ('man', 'woman')) NOT NULL,
                                     birth_date VARCHAR(20) NOT NULL,
                                     CONSTRAINT unique_actor_fields UNIQUE (name, gender, birth_date)
);

CREATE TABLE IF NOT EXISTS film (
                                    id SERIAL PRIMARY KEY,
                                    title VARCHAR(150) NOT NULL,
                                    description VARCHAR(1000) NOT NULL,
                                    release_date VARCHAR(12) NOT NULL,
                                    rating INT NOT NULL,
                                    CONSTRAINT rating CHECK (rating >= 1 AND rating <= 10),
                                    CONSTRAINT unique_title_release_date UNIQUE (title, release_date)
);

CREATE TABLE IF NOT EXISTS film_actor (
                                          id SERIAL PRIMARY KEY,
                                          film_id INT NOT NULL,
                                          actor_id INT NOT NULL,
                                          CONSTRAINT film_actor_film_id_fkey FOREIGN KEY (film_id) REFERENCES film(id),
                                          CONSTRAINT film_actor_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES actor(id),
                                          CONSTRAINT unique_film_actor UNIQUE (film_id, actor_id)
);



CREATE TABLE IF NOT EXISTS sessions (
                                        id VARCHAR(32),
                                        user_id INT NOT NULL,
                                        CONSTRAINT unique_session_id UNIQUE (id),
                                        CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users(id)
);
