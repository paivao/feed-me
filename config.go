package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	Port     int
	Database struct {
		Driver string
		Host   string
		Port   int
		Name   string
		User   string
		Pass   string
	}
}

func LoadConfiguration(filename string) (*Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return &config, err
	}
	err = json.NewDecoder(file).Decode(&config)
	return &config, err
}

func (c *Config) ConnectDB() (*gorm.DB, error) {
	switch c.Database.Driver {
	case "mysql":
		return gorm.Open(mysql.Open(c.getMysqlDSN()), &gorm.Config{})
	case "postgres":
		return gorm.Open(postgres.Open(c.getPostgresDSN()), &gorm.Config{})
	case "sqlite":
		return gorm.Open(sqlite.Open(c.Database.Host), &gorm.Config{})
	default:
		return nil, errors.New("Unrecognized database backend")
	}
}

func (c *Config) getMysqlDSN() string {
	var proto string
	host := c.Database.Host
	if host[0] == '.' || host[0] == '/' {
		proto = "unix"
	} else {
		proto = "tcp"
		if c.Database.Port > 0 {
			host = net.JoinHostPort(c.Database.Host, strconv.Itoa(c.Database.Port))
		}
	}

	return fmt.Sprintf("%s:%s@%s(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Database.User, c.Database.Pass, proto, host, c.Database.Name)
}

func (c *Config) getPostgresDSN() string {
	port := ""
	if c.Database.Port > 0 {
		port = fmt.Sprintf(" port=%d", c.Database.Port)
	}
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s%s", c.Database.Host, c.Database.User, c.Database.Pass, c.Database.Name, port)
}
