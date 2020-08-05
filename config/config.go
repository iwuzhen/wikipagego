/*
 * @Version: 0.0.1
 * @Author: ider
 * @Date: 2020-06-02 13:57:57
 * @LastEditors: ider
 * @LastEditTime: 2020-08-05 11:07:53
 * @Description:
 */
package config

import (
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	PgConn string `env:"PGCONN" envDefault:"host=192.168.1.220 port=5432 user=postgres dbname=wiki_knogen password=postgres sslmode=disable"`
}

func GetConfig() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Println("cfg error ", err)
	}
	return cfg
}
