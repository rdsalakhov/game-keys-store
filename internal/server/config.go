package server

type Config struct {
	DbConnection    string `yml:"db_connection"`
	RedisConnection string `yml:"redis_connection"`
	Port            string `yml:"port"`
	AccessSecret    string `yml:"access_secret"`
	RefreshSecret   string `yml:"refresh_secret"`
}
