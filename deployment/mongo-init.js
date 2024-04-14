// This mongo shell script will be executed when the container
// is created. It aims at creating and granding default permissions for turbex user

const user = _getEnv('TURBEX_DB_USER');
const pass = _getEnv('TURBEX_DB_PASS');
const db_name = _getEnv('MONGO_INITDB_DATABASE');
db.createUser({
  user: user,
  pwd: pass,
  roles: [{
    role: 'readWrite',
    db: db_name
  }]
})

