use wasm_bindgen::prelude::*;

#[wasm_bindgen]
extern "C" {
    pub fn alert(s: &str);
}

#[wasm_bindgen]
pub fn greet(name: &str) {
    alert(&format!("Hello, {}!", name));
}
pub mod encryption {
    use p384::{
        ecdh::{EphemeralSecret, SharedSecret},
        PublicKey,
    };
    use rand::rngs::OsRng;

    /// Generate a shared secret with the receiver public key and an ephemeral secret
    ///
    /// Return a tuple with the shared secret and the ephemeral secret public key.
    pub fn generate_shared_secret(receiver_pub_key: PublicKey) -> (SharedSecret, PublicKey) {
        let sender_ephemeral_secret = EphemeralSecret::random(&mut OsRng);
        let sender_ephemeral_pub = sender_ephemeral_secret.public_key();
        let shared_secret = sender_ephemeral_secret.diffie_hellman(&receiver_pub_key);
        (shared_secret, sender_ephemeral_pub)
    }
}
