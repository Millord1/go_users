package services

import (
	"microservices/models"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

// Satisfying repository.UserRepository
type MockMySQLRepo struct {
}

func checkisHashed(pw string) bool {
	return pw[0:4] == "$2a$" && len(pw) > 50
}

func (mock MockMySQLRepo) Update(user *models.User) (*models.User, error) {
	return user, nil
}

func (mock MockMySQLRepo) FindAllNames() (*[]models.User, error) {
	var users []models.User
	var user models.User
	gofakeit.Struct(&user)
	users = append(users, user)

	return &users, nil
}

func (mock MockMySQLRepo) FindByMail(email string) (*models.User, error) {
	return &models.User{
		Username: gofakeit.Username(),
		Email:    email,
		Password: gofakeit.Password(true, true, true, true, false, 10),
	}, nil
}

func (mock MockMySQLRepo) Save(user models.User) (*models.User, error) {
	return &user, nil
}

func TestCreateNewUser(t *testing.T) {
	repo := MockMySQLRepo{}
	var user models.User
	gofakeit.Struct(&user)

	newUser, err := CreateNewUser(repo, &user)
	assert.NoError(t, err)
	assert.True(t, checkisHashed(newUser.Password))
	assert.NotEqual(t, user.Password, newUser.Password)
}

func TestLoginUser(t *testing.T) {
	repo := MockMySQLRepo{}
	user, _ := repo.FindByMail(gofakeit.Email())

	hashedUser, err := user.HashPassword()
	assert.NoError(t, err)

	jwt, err := LoginUser(repo, hashedUser.Email, user.Password)
	assert.NoError(t, err)
	assert.Equal(t, "string", reflect.TypeOf(jwt))

	// TODO check jwt login
}

/* func TestUpdatePassword(t *testing.T) {
	repo := MockMySQLRepo{}
	var user models.User
	gofakeit.Struct(&user)
	pw := user.Password

	hashedUser, err := user.HashPassword()
	assert.NoError(t, err)
	assert.NotEqual(t, user.Password, hashedUser.Password)


}
*/
