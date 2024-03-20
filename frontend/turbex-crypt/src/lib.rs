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

pub mod key_management {
    use p384::{
        pkcs8::{DecodePrivateKey, EncodePrivateKey, SecretDocument},
        PublicKey, SecretKey,
    };
    use pbkdf2;
    use rand::rngs::OsRng;
    use sha2::Sha256;

    /// Return the key password and api password from a user password
    ///
    /// TODO: Consider using Zeroize to protect key and api password
    pub fn get_key_and_api_password(user_password: &[u8], user_salt: &[u8]) -> (Vec<u8>, Vec<u8>) {
        let mut raw_key = [0u8; 64];
        pbkdf2::pbkdf2_hmac::<Sha256>(user_password, user_salt, 600000, &mut raw_key);
        let key_password = raw_key[0..32].to_owned();
        let api_password = raw_key[32..64].to_owned();
        (key_password, api_password)
    }

    /// Return a new pair of Nist P384 keys, the private key is encrypted with the user password.
    pub fn generate_key() -> (SecretKey, PublicKey) {
        let secret_key = SecretKey::random(&mut OsRng);
        let pub_key = secret_key.public_key();
        (secret_key, pub_key)
    }

    pub fn encrypt_key(secret_key: SecretKey, key_password: &[u8]) -> SecretDocument {
        // TODO: Remove unwrap and handle errors
        secret_key
            .to_pkcs8_encrypted_der(&mut OsRng, key_password)
            .unwrap()
    }

    pub fn decrypt_private_key(
        encrypted_secret_key: SecretDocument,
        key_password: &[u8],
    ) -> SecretKey {
        let secret_key =
            SecretKey::from_pkcs8_encrypted_der(encrypted_secret_key.as_bytes(), key_password)
                .unwrap();
        secret_key
    }
}
