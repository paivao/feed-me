package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host string
	Port int
	Key  string
	Log  struct {
		Access string
		System string
		Level  string
	}
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

func (c *Config) ConnectDB() (*sql.DB, error) {
	switch c.Database.Driver {
	case "mysql":
		return sql.Open("mysql", c.getMysqlDSN())
	case "postgres":
		return sql.Open("postgres", c.getMysqlDSN())
	default:
		return nil, errors.New("unrecognized database backend")
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
