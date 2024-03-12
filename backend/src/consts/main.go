package consts

import "time"

// Used as context keys in gin to pass pointers to handlers
const CONTEXT_DB = "mongoClient"
const CONTEXT_ENV = "environmentVariables"
const CONTEXT_SESSION = "sessionData"

// Cliend side variables
const SESSION_COOKIE_NAME = "session"
const HOST_HEADER = "Host"

// The date format should satisfy the lexicographical order
// to optimise research in the mongodb collection
const DATE_FORMAT = time.RFC3339

const SESSION_CLEAN_INTERVAL = 60 // In seconds
