package consts

import "time"

// Used as context keys in gin to pass pointers to handlers
const CONTEXT_DB = "mongoClient"
const CONTEXT_ENV = "environmentVariables"
const CONTEXT_SESSION = "sessionData"
const CONTEXT_GRIDFS = "mongoGridFS"

// Cliend side variables
const SESSION_COOKIE_NAME = "session"
const SESSION_TOKEN_LENGTH_BYTES = 39 // In bytes (base64 encoded: divisible by 3 is better)
const HOST_HEADER = "Host"

// Password storage (NIST recommendation: 128 bits)
const SALT_LENGTH = 128 / 8

// The date format should satisfy the lexicographical order
// to optimise research in the mongodb collection
const DATE_FORMAT = time.RFC3339

const SESSION_CLEAN_INTERVAL = 60 // In seconds
