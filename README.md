# Turbex 🔐
A E2E encrypted filesharing in Rust 🦀, Svelte and Go
by [djex](https://github.com/djexvr) & [TurboThon](https://github.com/TurboThon/) 🐟 

This project is an example of a secure file sharing system. It uses client-side encryption
with state-of-the-art cryptography (ANSSI / NIST / OWASP).

## To the attention of all users

**Turbex's security is based on the strength of your password!**

Therefore we recommend you to follow the following guidelines:
- Your password must contain at least 15 characters (you can use any Unicode codepoints)
- You must not use it anywhere else
- You should consider generating a fully random password, for example using Keepass
- You should consider using a passphrase of at least 10 words, to that extend, pick 10 random words from
  a dictionnary, and concatenate them

If you loose your password, no one will be able to recover it. You will need to create another account.

## Deploy the app

This repository provides a `docker-compose.yml` almost ready for production.

For the moment, it is your responsibility, as administrators, to change default credentials
provided in the `docker-compose.yml` file.

## Developping

### Understanding the architecture

Several architecture schema are available in the TAD directory. We heavily recommend to take a look at
those documents to understand the overall architecture, API calls, protocols...

### Launching in dev mode

The project is devided in 3 parts you need to build separately:
- backend
- frontend/turbex-crypt
- frontend/frontend-svelte

For the sake of dev experience, a `docker-compose-dev.yml` is provided. It launches two containers:
a mongo and a mongo-express, so that you can easily see what is happening in the database.

Once those containers are ready, you can launch the backend by following instructions in the corresponding
README.md file (which act as a source of truth).

You also need to build the web-assembly rust crate used to handle cryptographic operations. Just like
you did with the backend, you need to follow the instructions provided in the corresponding README.md file.

Finaly, you can launch the frontend in `dev` mode by following the instructions in the corresponding README.md file.

For the sake of simplicity, all the commands are provided below. Make sure to launch them in separate terminals.

```bash
# Launch containers and backend
sudo docker compose -f docker-compose-dev.yml start && cd backend && swag init && go run .
# Build WASM and start frontend
cd frontend/turbex-crypt && wasm-pack build --target web && cd ../frontend-svelte && npm run dev
```

