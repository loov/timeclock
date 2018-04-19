package pgdb

import (
	"fmt"
	"math/rand"
	"time"
	"unicode"

	"github.com/icrowley/fake"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var rng = rand.New(rand.NewSource(1))

func rIntBetween(low, high int) int {
	return low + rng.Intn(high-low)
}

func (db *Database) FakeDatabase() error {
	return (&faker{Database: db}).Run()
}

type faker struct {
	*Database
	Errors []error

	Admin      int64
	Accountant int64
	Workers    []int64
}

func (db *faker) Run() error {
	fake.Seed(1)

	const NWorkers = 10
	const NCustomers = 5

	db.Admin = db.AddUser("Admin", "Admin", true, true, true)
	db.Accountant = db.AddUser("Accountant", "Accountant", false, true, false)

	for range [NWorkers]struct{}{} {
		firstName := fake.FirstName()
		lastName := fake.LastName()
		userId := db.AddUser(firstName, lastName, false, false, true)
		db.Workers = append(db.Workers, userId)

		if db.ShouldAbort() {
			return db.Error()
		}
	}

	for range [NCustomers]struct{}{} {
		customerId := db.AddCustomer(fake.Company())

		if db.ShouldAbort() {
			return db.Error()
		}

		for nprojects := rIntBetween(3, 7); nprojects > 0; nprojects-- {
			projectId, activities := db.AddProject(customerId, fake.ProductName())
			if db.ShouldAbort() {
				return db.Error()
			}

			db.AddActivities(customerId, projectId, activities)
			if db.ShouldAbort() {
				return db.Error()
			}
		}
	}

	return db.Error()
}

func (db *faker) Error() error {
	if len(db.Errors) > 0 {
		return fmt.Errorf("%v", db.Errors)
	}
	return nil
}

func (db *faker) ShouldAbort() bool { return len(db.Errors) > 0 }
func (db *faker) Failed(err error) bool {
	if err == nil {
		return false
	}
	db.Errors = append(db.Errors, err)
	return true
}

func (db *faker) AddUser(firstName, lastName string, admin, accountant, worker bool) int64 {
	name := firstName + " " + lastName
	email := firstName + lastName + "@example.com"

	result := db.QueryRow(`
			INSERT INTO Users
			(Alias, Name, Email)
			VALUES ($1, $2, $3)
			RETURNING ID
		`, string(firstName[0])+". "+lastName, name, email)

	var userId int64
	if err := result.Scan(&userId); db.Failed(err) {
		return -1
	}

	key, err := bcrypt.GenerateFromPassword([]byte(firstName), bcrypt.DefaultCost)
	if db.Failed(err) {
		return -1
	}

	_, err = db.Exec(`
			INSERT INTO Credentials 
			(UserID, Provider, Key)
			VALUES ($1, $2, $3)
		`, userId, "password", key)
	if db.Failed(err) {
		return -1
	}

	_, err = db.Exec(`
			INSERT INTO Roles 
			(UserID, Admin, Accountant, Worker)
			VALUES ($1, $2, $3, $4)
		`, userId, admin, accountant, worker)
	if db.Failed(err) {
		return -1
	}

	return userId
}

func (db *faker) AddCustomer(customerName string) int64 {
	result := db.QueryRow(`
			INSERT INTO Customers
			(Slug, Name)
			VALUES ($1, $2)
			RETURNING ID
		`, slugify(customerName), customerName)

	var customerId int64
	if err := result.Scan(&customerId); db.Failed(err) {
		return -1
	}

	return customerId
}

func fakeActivities() []string {
	switch rng.Intn(3) {
	case 0:
		return []string{"Plumbing", "Driving", "Welding", "Assembly", "Delta"}
	case 1:
		return []string{"Alternate Rail", "Lower Rail", "Upper Rail", "Sigma", "Alpha"}
	case 2:
		return []string{"Driving", "Construction", "Painting"}
	}
	return []string{"Delta"}
}

func (db *faker) AddProject(customerId int64, projectName string) (int64, []string) {
	activities := fakeActivities()
	result := db.QueryRow(`
			INSERT INTO Projects
			(CustomerID, Slug, Name, Completed, Activities, Description)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING ID
		`, customerId, fake.DigitsN(rIntBetween(3, 5)), projectName, rIntBetween(0, 2) == 0, pq.StringArray(activities), fake.Paragraphs())

	var projectId int64
	if err := result.Scan(&projectId); db.Failed(err) {
		return -1, nil
	}

	return projectId, fakeActivities()
}

func (db *faker) AddActivities(customerId int64, projectId int64, activities []string) {
	nworkers := rIntBetween(2, 4)
	var workers []int64
	for _, k := range rng.Perm(len(db.Workers))[:nworkers] {
		workers = append(workers, db.Workers[k])
	}

	ndays := rIntBetween(15, 30)
	simulatedTime := time.Now().AddDate(0, 0, -rIntBetween(ndays, ndays+30))
	for ; ndays > 0; ndays-- {
		weekday := simulatedTime.Weekday()
		if weekday == time.Sunday || weekday == time.Saturday {
			if rIntBetween(0, 6) > 0 {
				continue
			}
		} else {
			if rIntBetween(0, 20) > 0 {
				continue
			}
		}

		for nactivities := rIntBetween(1, 5); nactivities > 0; nactivities-- {
			workerId := workers[rIntBetween(0, len(workers)-1)]
			activity := activities[rIntBetween(0, len(activities)-1)]

			_, err := db.Exec(`
				INSERT INTO Activities
				(WorkerID, ProjectID, Time, Name, Amount)
				VALUES ($1, $2, $3, $4, $5)
			`, workerId, projectId, simulatedTime, activity, float64(rIntBetween(1, 8))/2)
			if db.Failed(err) {
				return
			}
		}
		simulatedTime = simulatedTime.AddDate(0, 0, 1)
	}
}

func slugify(s string) string {
	slug := ""
	for _, r := range s {
		if unicode.IsNumber(r) || unicode.IsLetter(r) {
			r = unicode.ToLower(r)
			slug += string(r)
		}
	}
	if len(slug) > 5 {
		slug = slug[:5]
	}
	return slug
}
