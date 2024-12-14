package repository

import (
	sq "github.com/jmoiron/sqlx"
	"gitlab.crja72.ru/gospec/go9/netevent/auth_service/internal/db/postgres"
	"gitlab.crja72.ru/gospec/go9/netevent/auth_service/internal/db/postgres/models"
)

/*
  rpc Register(RegisterRequest) returns (RegisterResponse);
  rpc Authenticate(AuthenticateRequest) returns (AuthenticateResponse);
  rpc Authorise(AuthoriseRequest) returns (AuthoriseResponse);
  rpc GetInterests(GetInterestsRequest) returns (GetInterestsResponse);
*/
type UserRepository struct {
	db *sqlx.DB
}
func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db.Db}
}


/*
	id INT PRIMARY KEY NOT NULL,
	name VARCHAR(30),
	password VARCHAR(30),
	email   VARCHAR(30),
	role    int,
	interest TEXT,
	accesstkn TEXT,
	accessttl INT,
	refreshtkn TEXT,
	refreshttl INT
*/
func (u *UserRepository) NewUser(users *models.User) (error){
	_, err := sq.Insert("tuser").Columns("id", "name", "password", "email", "interests").
				Values(users.Id, users.Name, users.Password, users.Email, users.Interests).
				RunWith(u.db)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) AuthUser(users *models.User) (err){
	
	res, err := sq.Select("name", "password").
				From("tuser").
				Where(sq.And{sq.Eq{"name", users.Name}, sq.Eq{"password", users.Password}}).
				RunWith(u.db)

	if err != nil {
		return err
	}
	us := models.User{}
	res.QueryRow().Scan(&us.Name, &us.Password)
	log.Println(us)


	return nil
}


func (u *UserRepository) GetId(users *models.User) (int, error){
	res, err := sq.Select("id").
	From("tuser").
	Where(sq.Eq{"name", users.Name}).
	RunWith(u.db)

	if err != nil {
	return err
	}
	us := models.User{}
	res.QueryRow().Scan(&us.Name, &us.Password)
	log.Println(us)

	if err != nil {
		return err
	}
	return nil
}


func (u *UserRepository) (users *models.User) (error){


	if err != nil {
		return err
	}
	return nil
}

