# Turbex Crypt 🦀 
The wasm crate used by turbex for cryptography on client-side

## Getting started
First you need to install `wasm-pack` in order to compile to WebAssembly

``` shell
cargo install wasm-pack
```

Then build the crate with

``` shell
wasm-pack build --target web
```

Finally you can import the wasm function in JS with `./pkg/turbex_crypt.js`

