package actor

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type ActorRepository struct {
	db *sql.DB
}

func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{
		db: db,
	}
}

func (repo *ActorRepository) Add(actor *Actor) error {
	op := "actor_repo.Add"
	_, err := repo.db.Exec(`INSERT INTO actor(name, gender, birth_date) VALUES($1, $2, $3) `, actor.Name, actor.Gender, actor.BirthDate)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (repo *ActorRepository) GetActorId(actor *Actor) (int64, error) {
	op := "actor_repo.GetByID"
	row := repo.db.QueryRow("SELECT id FROM actor where name = $1 and  gender = $2 and birth_date = $3",
		actor.Name, actor.Gender, actor.BirthDate)
	var actor_id int64
	err := row.Scan(&actor_id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return actor_id, nil
}

func (repo *ActorRepository) Update(actor_id int64, newActor *Actor) error {
	op := "actor_repo.UpdateActor"
	_, err := repo.db.Query("UPDATE actor SET name = $1, gender = $2, birth_date = $3 WHERE id = $4",
		newActor.Name, newActor.Gender, newActor.BirthDate, actor_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (repo *ActorRepository) Delete(actor_id int64) error {
	op := "actor_repo.DeleteActor"
	_, err := repo.db.Query("DELETE FROM actor WHERE id = $1", actor_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
