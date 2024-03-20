use base64::prelude::*;
use p384::pkcs8::{EncodePrivateKey, EncodePublicKey};
use pkcs8::LineEnding;
use serde::{Deserialize, Serialize};
use serde_wasm_bindgen;
use wasm_bindgen::prelude::*;

#[wasm_bindgen]
extern "C" {
    pub fn alert(s: &str);
}

#[wasm_bindgen]
pub fn greet(name: &str) {
    alert(&format!("Hello, {}!", name));
}

#[wasm_bindgen]
pub fn get_api_password(user_password: &[u8]) -> JsValue {
    // TODO: User specific salt
    let (_, api_password) = key_management::get_key_and_api_password(user_password, b"User Salt");
    let encoded_api_password = BASE64_STANDARD.encode(api_password);
    JsValue::from_str(&encoded_api_password)
}

#[wasm_bindgen]
pub fn get_encrypted_file() {}

#[wasm_bindgen]
pub fn get_decrypted_file() {}

#[derive(Serialize, Deserialize)]
pub struct KeysAndPassword {
    api_password: String,
    encrypted_key: String,
    pub_key: String,
}

#[wasm_bindgen]
pub fn get_new_keys_and_password(user_password: &[u8]) -> Result<JsValue, JsValue> {
    let (key_password, api_password) =
        key_management::get_key_and_api_password(user_password, b"User Salt");
    let (secret_key, pub_key) = key_management::generate_key();
    let secret_key_pem = secret_key
        .to_pkcs8_encrypted_pem(&mut rand::rngs::OsRng, &key_password, LineEnding::default())
        .unwrap();
    Ok(serde_wasm_bindgen::to_value(&KeysAndPassword {
        api_password: BASE64_STANDARD.encode(api_password),
        encrypted_key: secret_key_pem.to_string(),
        pub_key: pub_key.to_public_key_pem(LineEnding::default()).unwrap(),
    })?)
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

pub mod decryption {
    use p384::{ecdh, ecdh::SharedSecret, PublicKey, SecretKey};

    pub fn decrypt_file() {}

    pub fn generate_shared_secret(secret_key: SecretKey, public_key: PublicKey) -> SharedSecret {
        ecdh::diffie_hellman(secret_key.to_nonzero_scalar(), public_key.as_affine())
    }
}
