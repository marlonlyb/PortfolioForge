package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Applicatition name del pool de conexiones
const AppName = "EDcommerce de MLB"

func NewDBConnection() (*pgxpool.Pool, error) {

	min := 3
	max := 100
	minConn := os.Getenv("DB_MIN_CONN")
	maxConn := os.Getenv("DB_MAX_CONN")
	validateMinMaxConn(minConn, maxConn, &min, &max)

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")

	connString := makeURL(user, pass, host, port, dbName, sslMode, min, max)
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		//return nil, model.NewError()
		return nil, fmt.Errorf("%s %w", "pgxpool.ParseConfig()", err)
	}

	config.ConnConfig.RuntimeParams["application_name"] = AppName

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s %w", "pgxpool.NewWithConfig()", err)
	}

	return pool, nil
}

// Verificación de numero de conexiones, validas y dentro del rango
func validateMinMaxConn(minConn string, maxConn string, min *int, max *int) {

	vMinConn := 0
	vMaxConn := 0
	vMinCond := false
	vMaxCond := false

	if minConn != "" {
		v, err := strconv.Atoi(minConn)
		if err != nil {
			log.Println("warning: DB_MIN_COONN has not a valid value, we will set min connections to", min)
			vMinCond = false
		} else {
			vMinConn = v
			vMinCond = true
		}
	}

	if maxConn != "" {
		v, err := strconv.Atoi(maxConn)
		if err != nil {
			log.Println("warning: DB_MAX_COONN has not a valid value, we will set max connections to", max)
			vMaxCond = false
		} else {
			vMaxConn = v
			vMaxCond = true
		}
	}

	if vMinCond && vMaxCond {
		if vMinConn <= vMaxConn {
			if vMinConn >= *min && vMinConn <= *max {
				*min = vMinConn
			} else {
				log.Printf("warning: DB_MIN_CONN has not a valid value: [%v, %v]", *min, *max)
			}
			if vMaxConn >= *min && vMaxConn <= *max {
				*max = vMaxConn
			} else {
				log.Printf("warning: DB_MAX_CONN has not a valid value: [%v, %v]", *min, *max)
			}
		} else {
			log.Printf("warning: DB_MAX_CONN: %v must be greater than DB_MIN_CONN: %v", vMaxConn, vMinConn)
		}
	}
}

// función que arma el string de conexiones
func makeURL(user, pass, host, port, dbName, sslMode string, minConn, maxConn int) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_min_conns=%d pool_max_conns=%d",
		user,
		pass,
		host,
		port,
		dbName,
		sslMode,
		minConn,
		maxConn)
}
