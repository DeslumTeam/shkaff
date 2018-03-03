package consts

const (
	CONFIG_FILE                   = "config.json"
	DEFAULT_HOST                  = "localhost"
	DEFAULT_DATABASE_PORT         = 5432
	DEFAULT_AMQP_PORT             = 5672
	DEFAULT_STATDB_PORT           = 9000
	DEFAULT_MONGO_RESTORE_PORT    = 27018
	DEFAULT_DATABASE_DB           = "postgres"
	DEFAULT_REFRESH_DATABASE_SCAN = 60

	INVALID_DATABASE_HOST     = "Database host in config file is empty. Shkaff set '%s'\n"
	INVALID_DATABASE_PORT     = "Database port %d in config file invalid. Shkaff set '%d'\n"
	INVALID_DATABASE_DB       = "Database name in config file is empty. Shkaff set '%s'\n"
	INVALID_DATABASE_USER     = "Database user name is empty"
	INVALID_DATABASE_PASSWORD = "DatP user name is empty"
	INVALID_AMQP_PASSWORD     = "AMQP password is empty"

	INVALID_AMQP_HOST = "AMQP host in config file is empty. Shkaff set '%s'\n"
	INVALID_AMQP_PORT = "AMPQ port %d in config file invalid. Shkaff set '%d'\n"
	INVALID_AMQP_USER = "AMQ user invalid"

	INVALID_STATDB_HOST = "STATDB host in config file is empty. Shkaff set '%s'\n"
	INVALID_STATDB_PORT = "STATDB port %d in config file invalid. Shkaff set '%d'\n"

	INVALID_MONGO_RESTORE_HOST = "MONGO_RESTORE host in config file is empty. Shkaff set '%s'\n"

	INVALID_MONGO_RESTORE_PORT = "MONGO_RESTORE port %d in config file invalid. Shkaff set '%d'\n"

	RMQ_URI_TEMPLATE  = "amqp://%s:%s@%s:%d/%s"
	PSQL_URI_TEMPLATE = "postgres://%s:%s@%s:%d/%s?sslmode=disable"

	REQUEST_GET_STARTTIME = `SELECT task_id, db.db_id, user_id, "verb", thread_count,
    ipv6, gzip, host, port, databases, db_user, db_password, dumpfolder, tp.type as db_type
	FROM shkaff.tasks t 
	INNER JOIN shkaff.db_settings db 
	ON t.db_id = db.db_id 
	INNER JOIN shkaff.types tp ON tp.type_id = db.type_id
	WHERE (months @> '{}' or months @> ARRAY[%d]) 
	AND (days @> '{}' or days @> ARRAY[%d])
	AND (hours @> '{}' or hours @> ARRAY[%d])
	AND (minutes <= %d) 
	AND t.is_active = true AND t.is_delete = false;`

	REQUESR_UPDATE_ACTIVE = "UPDATE shkaff.tasks SET is_active = $1 WHERE task_id = $2 and AND is_delete = false;"

	MONGO_CLIENT_COMMAND  = "mongo"
	MONGO_DUMP_COMMAND    = "mongodump"
	MONGO_RESTORE_COMMAND = "mongorestore"
	MONGO_HOST_KEY        = "--host"
	MONGO_PORT_KEY        = "--port"
	MONGO_LOGIN_KEY       = "--username"
	MONGO_PASS_KEY        = "--password"
	MONGO_IPV6_KEY        = "--ipv6"
	MONGO_DATABASE_KEY    = "--db"
	MONGO_COLLECTION_KEY  = "--collection"
	MONGO_GZIP_KEY        = "--gzip"
	MONGO_PARALLEL_KEY    = "-j"

	CACHEPATH = "cache/"
)
