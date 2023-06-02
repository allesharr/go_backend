package api

type PropertiesStruct struct {
	APP_LISTEN_PORT           int  `json:"app_listen_port"`
	APP_AUTH_ENABLED          bool `json:"app_auth_enabled"`
	APP_SESSION_TIMEOUT_HOURS int  `json:"app_session_timeout_hours"`
	APP_DEBUG_MODE            bool `json:"app_debug_mode"`

	DB_HOST          string `json:"db_host"`
	DB_PORT          int    `json:"db_port"`
	DB_DATABASE_NAME string `json:"db_database_name"`
	DB_USERNAME      string `json:"db_username"`
	DB_PASSWORD      string `json:"db_password"`
}
