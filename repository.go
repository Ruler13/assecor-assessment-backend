package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var colorMap = map[int]string{
	1: "blau",
	2: "grün",
	3: "violett",
	4: "rot",
	5: "gelb",
	6: "türkis",
	7: "weiß",
}

type PersonRepository interface {
	GetAll() ([]Person, error)
	GetByID(id int) (Person, error)
	GetByColor(color string) ([]Person, error)
	Add(person Person) (Person, error)
}

type CSVRepository struct {
	persons []Person
}

func NewCSVRepository() *CSVRepository {
	repo := &CSVRepository{}
	err := repo.loadPersons()
	if err != nil {
		panic(err)
	}
	return repo
}

func (r *CSVRepository) loadPersons() error {
	file, err := os.Open("sample-input.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for i, record := range records {
		if len(record) < 4 {
			continue
		}
		lastname := strings.TrimSpace(record[0])
		name := strings.TrimSpace(record[1])
		address := strings.TrimSpace(record[2])
		colorIDStr := strings.TrimSpace(record[3])

		parts := strings.Fields(address)
		if len(parts) < 2 {
			continue
		}
		zipcode := parts[0]
		city := strings.Join(parts[1:], " ")

		colorID, err := strconv.Atoi(colorIDStr)
		if err != nil {
			continue
		}
		color, ok := colorMap[colorID]
		if !ok {
			color = "unknown"
		}

		person := Person{
			ID:       i + 1,
			Name:     name,
			Lastname: lastname,
			Zipcode:  zipcode,
			City:     city,
			Color:    color,
		}
		r.persons = append(r.persons, person)
	}
	return nil
}

func (r *CSVRepository) GetAll() ([]Person, error) {
	return r.persons, nil
}

func (r *CSVRepository) GetByID(id int) (Person, error) {
	for _, p := range r.persons {
		if p.ID == id {
			return p, nil
		}
	}
	return Person{}, fmt.Errorf("person not found")
}

func (r *CSVRepository) GetByColor(color string) ([]Person, error) {
	var result []Person
	for _, p := range r.persons {
		if p.Color == color {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *CSVRepository) Add(person Person) (Person, error) {
	person.ID = len(r.persons) + 1
	r.persons = append(r.persons, person)
	return person, nil
}

type DBConfig struct {
	Driver   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Path     string
}

func (c DBConfig) ConnectionString() string {
	switch c.Driver {
	case "sqlite":
		return c.Path
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Password, c.Database)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			c.User, c.Password, c.Host, c.Port, c.Database)
	default:
		return c.Path
	}
}

type GORMRepository struct {
	db *gorm.DB
}

func NewGORMRepository(config DBConfig) (*GORMRepository, error) {
	db, err := gorm.Open(sqlite.Open(config.ConnectionString()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Person{})
	return &GORMRepository{db: db}, nil
}

func (r *GORMRepository) GetAll() ([]Person, error) {
	var persons []Person
	err := r.db.Find(&persons).Error
	return persons, err
}

func (r *GORMRepository) GetByID(id int) (Person, error) {
	var person Person
	err := r.db.First(&person, id).Error
	return person, err
}

func (r *GORMRepository) GetByColor(color string) ([]Person, error) {
	var persons []Person
	err := r.db.Where("color = ?", color).Find(&persons).Error
	return persons, err
}

func (r *GORMRepository) Add(person Person) (Person, error) {
	err := r.db.Create(&person).Error
	return person, err
}
