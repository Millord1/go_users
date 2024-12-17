package repository

import (
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
	Save(user models.User) error
	Update(user *models.User) error
	FindAllNames() (*[]models.User, error)
	FindByMail(email string) (*models.User, error)
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
	repo := DbConnect(utils.GetEnvFile().Name)

	if !repo.db.Migrator().HasTable("users") {
		migrate(*repo)
	}
}

func migrate(repo MySQLRepository) {
	migrErr := repo.db.AutoMigrate(&models.User{})
	if migrErr != nil {
		log.Fatalln(migrErr)
	}

	user := models.User{
		Username: "admin", Email: "admin@test.com", Password: os.Getenv("ADMIN_PW"),
	}
	hashErr := user.HashPassword()
	if hashErr != nil {
		log.Fatalln(hashErr)
	}

	err := repo.Save(user)
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

func (repo MySQLRepository) Update(user *models.User) error {
	err := repo.db.Model(&user).Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo MySQLRepository) FindAllNames() (*[]models.User, error) {
	var dbUsers []models.User
	err := repo.db.Select("username").Find(&dbUsers).Error
	if err != nil {
		log.Fatalln(err)
	}
	return &dbUsers, err
}

func (repo MySQLRepository) FindByMail(email string) (*models.User, error) {
	var user models.User
	err := repo.db.Where("email = ?", email).Find(&user).Error
	if err != nil {
		log.Fatalln(err)
	}

	return &user, err
}

func (repo MySQLRepository) Save(user models.User) error {
	err := repo.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
