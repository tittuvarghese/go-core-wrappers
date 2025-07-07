package storage

import (
	"github.com/tittuvarghese/go-core-wrappers/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var log = logger.NewLogger("storage-engine")

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

const CreateCommand = "create"
const UpsertCommand = "upsert"
const UpdateCommand = "update"
const DeleteCommand = "delete"

type RelationalDB struct {
	Connection string
	Instance   *gorm.DB
}

type Operation struct {
	Command   string
	Model     interface{}
	Condition interface{}
	Expr      Expr // Required only for update operations
}

type Expr struct {
	Column string
	Value  interface{}
}

type AtomicTransaction struct {
	Operations []Operation
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
func (handler *RelationalDB) QueryByCondition(model interface{}, condition map[string]interface{}, preload ...string) ([]interface{}, error) {
	var results []interface{}
	var queryBuilder = handler.Instance

	for _, tableName := range preload {
		queryBuilder = queryBuilder.Preload(tableName)
	}

	// Perform the query with conditions
	if err := queryBuilder.Where(condition).Find(model).Error; err != nil {
		return nil, err
	}

	// Cast model to slice of interface{}
	results = append(results, model)
	return results, nil
}

func (handler *RelationalDB) BuildExpr(column string, args ...interface{}) clause.Expr {
	return gorm.Expr(column, args)
}

func (handler *RelationalDB) Transaction(ops AtomicTransaction) error {
	tx := handler.Instance.Begin()
	for _, record := range ops.Operations {
		switch record.Command {
		case CreateCommand:
			if err := tx.Create(record.Model).Error; err != nil {
				log.Error("Failed to perform create", err)
				tx.Rollback()
				return err
			}
		case UpsertCommand:
			if err := tx.Save(record).Error; err != nil {
				log.Error("Failed to perform upsert", err)
				tx.Rollback()
				return err
			}
		case UpdateCommand:
			if err := tx.Model(record.Model).
				Where(record.Condition).
				Update(record.Expr.Column, record.Expr.Value).Error; err != nil {
				log.Error("Failed to perform update operation", err)
				tx.Rollback()
				return err
			}
		case DeleteCommand:
			if err := tx.Delete(record.Model).Error; err != nil {
				log.Error("Failed to perform delete operation", err)
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}
