package main

import (
	"fmt"
	"os"
)

var ErrNoEnvVar = fmt.Errorf("no environment variable")

func main() {
	fmt.Println("hello world")
}

func Init() (*Config, error) {
	var password, user, database, host string

	if pwd, ok := os.LookupEnv("PGPASSWORD"); !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoEnvVar, "PGPASSWORD")
	} else {
		password = pwd
	}

	if usr, ok := os.LookupEnv("PGUSER"); !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoEnvVar, "PGUSER")
	} else {
		user = usr
	}

	if db, ok := os.LookupEnv("PGDATABASE"); !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoEnvVar, "PGDATABASE")
	} else {
		database = db
	}

	if h, ok := os.LookupEnv("PGHOST"); !ok {
		return nil, fmt.Errorf("%w: %s", ErrNoEnvVar, "PGHOST")
	} else {
		host = h
	}

	fmt.Println(password, user, database, host)

	return nil, fmt.Errorf("not implemented")
}
