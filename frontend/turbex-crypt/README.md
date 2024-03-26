# Turbex Crypt 🦀 
The wasm crate used by turbex for cryptography on client-side

## Getting started
First you need to install `wasm-pack` in order to compile to WebAssembly

```bash
cargo install wasm-pack
```

Then build the crate with

```bash
wasm-pack build --target web
```

Finally you can import the wasm function in JS with `./pkg/turbex_crypt.js`

## Development

To run unit tests:

```bash
cargo test
```

To run integration tests:
1. Compile the project using `wasm-pack`
2. Open a small webserver serving your local directories like `python3 -m http.server`
3. Open your web-browser console and paste the content of the `test/integration-test.js` file

