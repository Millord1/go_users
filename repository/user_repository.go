package repository

import (
	"errors"
	"log"
	"microservices/models"
	"microservices/utils"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLRepository struct {
	User     string
	Password string
	Protocol string
	Address  string
	Port     string
	Name     string
	db       *gorm.DB
}

type UserRepository interface {
	CreateNewUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	UpdatePassword(user *models.User, oldPw string) (*models.User, error)
	GetAllUserNames() (*[]models.User, error)
	GetUserByMail(email string) (*models.User, error)
}

func getMySQLRepo() MySQLRepository {
	return MySQLRepository{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Protocol: os.Getenv("DB_PROTOCOL"),
		Address:  os.Getenv("DB_ADDRESS"),
		Name:     os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}
}

func init() {

	// TODO test migration if table already exists
	repo := DbConnect(utils.GetEnvFile().Name)
	migrErr := repo.db.AutoMigrate(&models.User{})
	if migrErr != nil {
		log.Fatalln(migrErr)
	}

	// TODO test it
	// Create default user admin with pw specified in .env or .env.local file
	user := models.User{
		Username: "admin", Email: "admin@test.com", Password: os.Getenv("ADMIN_PW"),
	}
	userToPush, hashErr := user.HashPassword()
	if hashErr != nil {
		log.Fatalln(hashErr)
	}

	_, err := repo.CreateNewUser(userToPush)
	if err != nil {
		log.Fatalln(err)
	}
}

func DbConnect(envFile string) *MySQLRepository {
	// Init connection to database from specified env file variables
	err := godotenv.Load(envFile)

	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}

	repo := getMySQLRepo()

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: repo.User + ":" + repo.Password + "@" + repo.Protocol +
			"(" + repo.Address + ":" + repo.Port + ")/" + repo.Name +
			"?charset=utf8mb4&parseTime=True&loc=Local",
	}), &gorm.Config{})

	repo.db = db

	if err != nil {
		panic("Failed to connect database")
	}

	return &repo
}

func (repo *MySQLRepository) CreateNewUser(user *models.User) (*models.User, error) {

	userToPush, hashErr := user.CheckPasswordIsHashed()
	if hashErr != nil {
		log.Fatalln(hashErr)
	}

	err := repo.db.Create(&userToPush).Error
	if err != nil {
		log.Fatalln(err)
	}
	return user, err
}

func (repo *MySQLRepository) UpdateUser(user *models.User) (*models.User, error) {

	// building a new User to avoid pw modification without pw check
	err := repo.db.Model(&user).Updates(models.User{
		Username: user.Username,
		Email:    user.Email,
	}).Error
	if err != nil {
		log.Fatalln(err)
	}
	return user, err
}

func (repo *MySQLRepository) UpdatePassword(user *models.User, oldPw string) (*models.User, error) {
	if !user.VerifyPassword(oldPw) {
		err := errors.New("wrong password")
		log.Fatalln(err)
		return user, err
	}

	err := repo.db.Model(&user).Updates(user).Error
	if err != nil {
		log.Fatalln(err)
	}
	return user, err
}

func (repo *MySQLRepository) GetAllUserNames() (*[]models.User, error) {
	var dbUsers []models.User
	err := repo.db.Select("username").Find(&dbUsers).Error
	if err != nil {
		log.Fatalln(err)
	}
	return &dbUsers, err
}

func (repo *MySQLRepository) GetUserByMail(email string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		log.Fatalln(err)
	}
	return &user, err
}
