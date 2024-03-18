package postrgeSql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func New(storagePath string) (*sql.DB, error) {
	const op = "storage.postgre.New"

	db, err := sql.Open("postgres", storagePath)

	err = db.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		log.Fatalf("cant connect to db, err: %v\n", err)
	}

	createTableStmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS actor (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		gender VARCHAR(10) CHECK (gender IN ('man', 'woman')) NOT NULL,
		birth_date VARCHAR(20) NOT NULL,
		CONSTRAINT unique_actor_fields UNIQUE (name, gender, birth_date)
	)
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = createTableStmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createFilmStmt, err := db.Prepare(`	CREATE TABLE IF NOT EXISTS film(
    		id SERIAL PRIMARY KEY,
    		title varchar(150) NOT NULL,
    		description varchar(1000) NOT NULL,
    		release_date varchar(12) NOT NULL,
    		rating int NOT NULL,
    		CONSTRAINT rating CHECK (rating >= 1 AND rating <= 10),
	        CONSTRAINT unique_title_release_date UNIQUE (title, release_date)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = createFilmStmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createFilmActorStmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS film_actor(
    		id SERIAL PRIMARY KEY,
    		film_id int NOT NULL,
    		actor_id int NOT NULL,
    		CONSTRAINT film_actor_film_id_fkey FOREIGN KEY (film_id) REFERENCES film(id),
    		CONSTRAINT film_actor_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES actor(id),
    		CONSTRAINT unique_film_actor UNIQUE (film_id, actor_id));
    		`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = createFilmActorStmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, login VARCHAR(255) NOT NULL, password BYTEA NOT NULL, role bool, CONSTRAINT unique_login UNIQUE (login))")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sessions (id VARCHAR(32), user_id INT NOT NULL, CONSTRAINT unique_session_id UNIQUE (id), CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES users(id))")

	return db, nil
}
