package core

import "github.com/sirupsen/logrus"

type Core struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
	Log *logrus.Logger
}