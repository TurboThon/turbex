package consts

import "time"

// Used in gin contexts to retreive a pointer to the database client
const CONTEXT_DB = "mongoClient"
const CONTEXT_ENV = "environmentVariables"
const CONTEXT_SESSION = "sessionData"

const SESSION_COOKIE_NAME = "session"
const HOST_HEADER = "Host"

const DATE_FORMAT = time.RFC3339

// In seconds
const SESSION_CLEAN_INTERVAL = 60
