package storage

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Storage interface {
	Open() error
	Close() error
	AutoMigrate(schema ...interface{}) error
	Insert(record interface{}) error
	Update(record interface{}) error
	Delete(record interface{}) error
	QueryAll(model interface{}) ([]interface{}, error)
	QueryByCondition(model interface{}, condition map[string]interface{}) ([]interface{}, error)
}

type RelationalDB struct {
	Connection string
	Instance   *gorm.DB
}

// Initialize the database connection
func NewRelationalDbHandler(dsn string) (*RelationalDB, error) {
	return &RelationalDB{Connection: dsn}, nil
}

func (handler *RelationalDB) Open() error {
	// Open the MariaDB or MySQL connection with GORM
	db, err := gorm.Open(mysql.Open(handler.Connection), &gorm.Config{})
	if err != nil {
		return err
	}
	handler.Instance = db
	return nil
}

func (handler *RelationalDB) Close() error {
	// GORM does not have an explicit Close method as it manages connections automatically
	// But, if you have a custom connection pool or management system, you could add logic here
	sqlDB, err := handler.Instance.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close() // Close the SQL connection
}

func (handler *RelationalDB) AutoMigrate(schema ...interface{}) error {
	// AutoMigrate requires the models to be passed as arguments
	err := handler.Instance.AutoMigrate(schema...)
	if err != nil {
		return err
	}
	return nil
}

func (handler *RelationalDB) Insert(record interface{}) error {
	// Insert the record into the database
	if err := handler.Instance.Create(record).Error; err != nil {
		return err
	}
	return nil
}

// Update updates an existing record in the database
func (handler *RelationalDB) Update(record interface{}) error {
	if err := handler.Instance.Save(record).Error; err != nil {
		return err
	}
	return nil
}

// Delete deletes a record from the database by its ID
func (handler *RelationalDB) Delete(record interface{}) error {
	if err := handler.Instance.Delete(record).Error; err != nil {
		return err
	}
	return nil
}

// QueryAll retrieves all records of a specific model
func (handler *RelationalDB) QueryAll(model interface{}) ([]interface{}, error) {
	var results []interface{}

	// Perform the query to fetch all records
	if err := handler.Instance.Find(&results).Error; err != nil {
		return nil, err
	}

	// Return the results (already populated by Find)
	return results, nil
}

// QueryByCondition retrieves records based on dynamic conditions (like a WHERE clause)
func (handler *RelationalDB) QueryByCondition(model interface{}, condition map[string]interface{}) ([]interface{}, error) {
	var results []interface{}

	// Perform the query with conditions
	if err := handler.Instance.Where(condition).Find(model).Error; err != nil {
		return nil, err
	}

	// Cast model to slice of interface{}
	results = append(results, model)
	return results, nil
}
