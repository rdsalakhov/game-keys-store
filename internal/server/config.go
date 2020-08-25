package server

type Config struct {
	DbConnection string `yml:"db_connection"`
	Port         string `yml:"port"`
}
