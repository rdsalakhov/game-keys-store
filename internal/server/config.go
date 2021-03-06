package server

type Config struct {
	DbConnection          string  `yml:"dbConnection"`
	RedisConnection       string  `yml:"redisConnection"`
	RedisPassword         string  `yml:"redisPassword"`
	Port                  string  `yml:"port"`
	AccessSecret          string  `yml:"accessSecret"`
	RefreshSecret         string  `yml:"refreshSecret"`
	PlatformFeeShare      float64 `yml:"PlatformFeeShare"`
	PlatformAccount       string  `yml:"PlatformAccount"`
	PlatformEmail         string  `yml:"PlatformEmail"`
	PlatformEmailPassword string  `yml:"PlatformEmailPassword"`
	NotificationSalt      string  `yml:"notificationSalt"`
}
