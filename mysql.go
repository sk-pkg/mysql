// Package mysql provides functionality for creating and managing MySQL database connections
// using the GORM library. It offers options for configuring single and multiple database
// connections with customizable connection pool settings.
package mysql

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

const (
	// defaultMaxIdleConn is the default maximum number of connections in the idle connection pool.
	defaultMaxIdleConn = 10
	// defaultMaxOpenConn is the default maximum number of open connections to the database.
	defaultMaxOpenConn = 50
	// defaultConnMaxLifetime is the default maximum amount of time a connection may be reused.
	defaultConnMaxLifetime = 3 * time.Hour
)

// Config represents the configuration for a MySQL database connection.
type Config struct {
	User     string // Database user
	Password string // Database password
	Host     string // Database host
	DBName   string // Database name
}

// Option is a function type used to apply configuration options.
type Option func(*option)

// option holds all the configurable options for database connections.
type option struct {
	dbConfigs       []Config      // Slice of database configurations
	gormConfig      gorm.Config   // GORM configuration
	maxIdleConn     int           // Maximum number of connections in the idle connection pool
	maxOpenConn     int           // Maximum number of open connections to the database
	connMaxLifetime time.Duration // Maximum amount of time a connection may be reused
}

// WithConfigs returns an Option that sets the database configurations.
//
// Parameters:
//   - cfg: One or more Config structs representing database configurations.
//
// Returns:
//   - An Option function that sets the database configurations when applied.
//
// Example:
//
//	cfg1 := Config{User: "user1", Password: "pass1", Host: "host1", DBName: "db1"}
//	cfg2 := Config{User: "user2", Password: "pass2", Host: "host2", DBName: "db2"}
//	db, err := New(WithConfigs(cfg1, cfg2))
func WithConfigs(cfg ...Config) Option {
	return func(o *option) {
		o.dbConfigs = cfg
	}
}

// WithGormConfig returns an Option that sets the GORM configuration.
//
// Parameters:
//   - cfg: A gorm.Config struct with desired GORM settings.
//
// Returns:
//   - An Option function that sets the GORM configuration when applied.
//
// Example:
//
//	gormCfg := gorm.Config{SkipDefaultTransaction: true}
//	db, err := New(WithGormConfig(gormCfg))
func WithGormConfig(cfg gorm.Config) Option {
	return func(o *option) {
		o.gormConfig = cfg
	}
}

// WithMaxIdleConn returns an Option that sets the maximum number of idle connections.
//
// Parameters:
//   - maxIdleConn: An integer representing the maximum number of idle connections.
//
// Returns:
//   - An Option function that sets the maximum idle connections when applied.
//
// Example:
//
//	db, err := New(WithMaxIdleConn(20))
func WithMaxIdleConn(maxIdleConn int) Option {
	return func(o *option) {
		o.maxIdleConn = maxIdleConn
	}
}

// WithConnMaxLifetime returns an Option that sets the maximum lifetime of connections.
//
// Parameters:
//   - connMaxLifetime: A time.Duration representing the maximum lifetime of a connection.
//
// Returns:
//   - An Option function that sets the maximum connection lifetime when applied.
//
// Example:
//
//	db, err := New(WithConnMaxLifetime(4 * time.Hour))
func WithConnMaxLifetime(connMaxLifetime time.Duration) Option {
	return func(o *option) {
		o.connMaxLifetime = connMaxLifetime
	}
}

// WithMaxOpenConn returns an Option that sets the maximum number of open connections.
//
// Parameters:
//   - maxOpenConn: An integer representing the maximum number of open connections.
//
// Returns:
//   - An Option function that sets the maximum open connections when applied.
//
// Example:
//
//	db, err := New(WithMaxOpenConn(100))
func WithMaxOpenConn(maxOpenConn int) Option {
	return func(o *option) {
		o.maxOpenConn = maxOpenConn
	}
}

// New initializes and returns a single database connection instance.
//
// Parameters:
//   - opts: A variadic list of Option functions to configure the database connection.
//
// Returns:
//   - A pointer to a gorm.DB instance representing the database connection.
//   - An error if the initialization fails.
//
// Example:
//
//	db, err := New(
//	    WithConfigs(Config{User: "user", Password: "pass", Host: "host", DBName: "db"}),
//	    WithMaxIdleConn(15),
//	    WithMaxOpenConn(75),
//	)
func New(opts ...Option) (*gorm.DB, error) {
	opt := setOption(opts...)
	if len(opt.dbConfigs) != 1 {
		return nil, errors.New("this method can only initialize one database connection instance")
	}

	return newConnect(&opt.dbConfigs[0], opt)
}

// NewMulti initializes and returns multiple database connection instances.
//
// Parameters:
//   - opts: A variadic list of Option functions to configure the database connections.
//
// Returns:
//   - A map with database names as keys and corresponding gorm.DB instances as values.
//   - An error if the initialization fails.
//
// Example:
//
//	cfg1 := Config{User: "user1", Password: "pass1", Host: "host1", DBName: "db1"}
//	cfg2 := Config{User: "user2", Password: "pass2", Host: "host2", DBName: "db2"}
//	dbs, err := NewMulti(WithConfigs(cfg1, cfg2), WithMaxIdleConn(15))
func NewMulti(opts ...Option) (map[string]*gorm.DB, error) {
	opt := setOption(opts...)

	if len(opt.dbConfigs) < 1 {
		return nil, errors.New("the number of database configurations to initialize cannot be 0")
	}

	dbs := make(map[string]*gorm.DB)
	for _, cfg := range opt.dbConfigs {
		conn, err := newConnect(&cfg, opt)
		if err != nil {
			return nil, err
		}

		dbs[cfg.DBName] = conn
	}

	return dbs, nil
}

// setOption applies the provided options and returns the resulting option struct.
//
// Parameters:
//   - opts: A variadic list of Option functions to apply.
//
// Returns:
//   - A pointer to the option struct with all options applied.
func setOption(opts ...Option) *option {
	opt := &option{
		maxIdleConn:     defaultMaxIdleConn,
		maxOpenConn:     defaultMaxOpenConn,
		connMaxLifetime: defaultConnMaxLifetime,
	}

	for _, f := range opts {
		f(opt)
	}

	return opt
}

// newConnect creates a new database connection with the given configuration and options.
//
// Parameters:
//   - cfg: A pointer to a Config struct containing database connection details.
//   - opt: A pointer to an option struct containing additional configuration options.
//
// Returns:
//   - A pointer to a gorm.DB instance representing the database connection.
//   - An error if the connection fails.
func newConnect(cfg *Config, opt *option) (*gorm.DB, error) {
	// Construct the DSN (Data Source Name) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.DBName)

	// Open the database connection
	db, err := gorm.Open(mysql.Open(dsn), &opt.gormConfig)
	if err != nil {
		return nil, err
	}

	// Get the underlying database/sql DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Configure the connection pool
	sqlDB.SetMaxIdleConns(opt.maxIdleConn)        // Set the maximum number of connections in the idle connection pool
	sqlDB.SetMaxOpenConns(opt.maxOpenConn)        // Set the maximum number of open connections to the database
	sqlDB.SetConnMaxLifetime(opt.connMaxLifetime) // Set the maximum amount of time a connection may be reused

	return db, nil
}
