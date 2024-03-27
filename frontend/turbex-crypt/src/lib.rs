use aes_gcm::{Aes256Gcm, Key};
use base64::prelude::*;
use p384::pkcs8::{EncodePrivateKey, EncodePublicKey};
use pkcs8::LineEnding;
use serde::{Deserialize, Serialize};
use wasm_bindgen::prelude::*;

// Following salts are choosen from the bytes returned by SHA512(b"Turbex file sharing")
// with 160 rounds. The following lines are the complete output
// \x19\x26\x15\x49\x36\xf2\x2c\x37
// \x0c\x63\x47\xee\x0d\x1c\x9f\xae
// \xa0\xff\xb7\x82\x40\x6d\x4f\xc5
// \xfe\x03\x95\x01\x88\x8b\xb7\xfc
// \xf1\x17\x39\x18\x60\x83\x4c\xd0
// \xa9\xe3\xf7\x89\x5e\x99\x1c\xcf
// \xdf\xfb\x4a\x56\x0e\x41\x9f\x3d
// \x98\x2a\x21\x7a\x9b\x8c\xf5\x9a
static USER_SALT: &[u8] = b"\x19\x26\x15\x49\x36\xf2\x2c\x37";
static KEY_SALT: &[u8] = b"\x0c\x63\x47\xee\x0d\x1c\x9f\xae";

// The crate we are using generates 96 bits nonces
// Ref: https://docs.rs/aes-gcm/latest/aes_gcm/type.Aes256Gcm.html
static NONCE_SIZE: usize = 96 / 8;

#[wasm_bindgen]
extern "C" {
    pub fn alert(s: &str);

    #[wasm_bindgen(js_namespace = console)]
    fn log(s: &str);
}
// macro_rules! console_log {
//     // Note that this is using the `log` function imported above during
//     // `bare_bones`
//     ($($t:tt)*) => (log(&format_args!($($t)*).to_string()))
// }

#[wasm_bindgen]
pub fn greet(name: &str) {
    alert(&format!("Hello, {}!", name));
}

#[derive(Serialize, Deserialize)]
#[wasm_bindgen(getter_with_clone)]
pub struct UserPasswords {
    // base64 encoded key password
    pub key_password: String,
    // base64 encoded api password
    pub api_password: String,
}

// Computes the user keys based on the user's password
#[wasm_bindgen]
pub fn get_api_password_and_key(user_password: String) -> UserPasswords {
    let user_password_bytes = user_password.as_bytes();
    // TODO: User specific salt
    let (key_password, api_password) =
        key_management::get_key_and_api_password(user_password_bytes, USER_SALT);

    UserPasswords {
        key_password: BASE64_STANDARD.encode(key_password),
        api_password: BASE64_STANDARD.encode(api_password),
    }
}

#[derive(Serialize, Deserialize)]
#[wasm_bindgen(getter_with_clone)]
pub struct UserKeys {
    // PEM encrypted (pkcs8) private key
    pub private_key: String,
    // PEM encoded (pkcs8) public key
    pub public_key: String,
}

// Creates a new key and protects it using the user_key_password
// user_key_password is base64 encoded
#[wasm_bindgen]
pub fn create_user_key(user_key_password: String) -> UserKeys {
    let key_password_bytes = BASE64_STANDARD.decode(user_key_password).unwrap();
    let (priv_key, pub_key) = key_management::generate_key();
    let protected_priv_key = key_management::encrypt_key(priv_key, &key_password_bytes);
    let encoded_pub_key = key_management::encode_public_key(pub_key);

    UserKeys {
        private_key: protected_priv_key.to_string(),
        public_key: encoded_pub_key,
    }
}

// Generates a new AES key, can be used as a PFK (Primary File Key)
// The returned key is base64 encoded
#[wasm_bindgen]
pub fn generate_aes_key() -> String {
    let aes_key = symetric_crypto::generate_aes_key();
    let b64key = BASE64_STANDARD.encode(aes_key);
    b64key
}

// Encrypts the provided file using the provided PFK
#[wasm_bindgen]
pub fn encrypt_file(file: &[u8], key: String) -> Vec<u8> {
    let key_bytes = BASE64_STANDARD.decode(key).unwrap();
    let aes_key = Key::<Aes256Gcm>::from_slice(&key_bytes);
    let encrypted_file = symetric_crypto::encrypt_file(file, aes_key);
    encrypted_file
}

#[derive(Serialize, Deserialize)]
#[wasm_bindgen(getter_with_clone)]
pub struct EncryptedPFKForRecipient {
    // encrypted_pfk is the base64 encoded form of the encrypted string
    pub encrypted_pfk: String,
    // ephemeral_pub_key is a PEM encoded Nist384
    pub ephemeral_pub_key: String,
}

// Encrypt the PFK using the recipient's public key
// pfk is a base64 encoded AES256 key
// public_key is a PEM encoded Nist384 public key
#[wasm_bindgen]
pub fn encrypt_pfk(pfk: String, public_key: String) -> EncryptedPFKForRecipient {
    let pfk_bytes = BASE64_STANDARD.decode(pfk).unwrap();
    let recipient_pubkey = key_management::decode_public_key(&public_key);
    let (shared_secret, sender_ephemeral_pub_key) =
        encryption::generate_shared_secret(recipient_pubkey);
    let encoded_public_key = sender_ephemeral_pub_key
        .to_public_key_pem(LineEnding::default())
        .unwrap();
    let aes_encrypted_key = symetric_crypto::encrypt_pfk(&pfk_bytes, &shared_secret);

    EncryptedPFKForRecipient {
        encrypted_pfk: BASE64_STANDARD.encode(aes_encrypted_key),
        ephemeral_pub_key: encoded_public_key,
    }
}

// Decrypts a file
// encrypted_file is an encrypted file
// encrypted_pfk is a base64 encoded pfk, which can be used to decrypt the previously mentioned
// file
// ephemeral_pubkey is the PEM encoded ephemeral public key that was used to compute the shared
// secret
// priv_key is the user's private key which is PEM encoded and protected
// priv_key_password is the user's priv_key password
#[wasm_bindgen]
pub fn decrypt_file(
    encrypted_file: &[u8],
    encrypted_pfk: String,
    ephemeral_pubkey: String,
    priv_key: String,
    priv_key_passwd: String,
) -> Vec<u8> {
    // Decrypt and load the user's private key
    let priv_key_passwd_bytes = BASE64_STANDARD.decode(priv_key_passwd).unwrap();
    let private_key = key_management::decrypt_private_key(priv_key.into(), &priv_key_passwd_bytes);
    // Decode the sender's ephemeral public key
    let public_key = key_management::decode_public_key(&ephemeral_pubkey);
    // Use the private key and public key to compute the shared secret
    let shared_secret = decryption::generate_shared_secret(private_key, public_key);
    // Decrypts the PFK using the previously computed shared secret
    let encrypted_pfk_bytes = BASE64_STANDARD.decode(encrypted_pfk).unwrap();
    let pfk = symetric_crypto::decrypt_pfk(&encrypted_pfk_bytes, &shared_secret);
    // Use the PFK to decrypt the file
    let clear_file = symetric_crypto::decrypt_file(encrypted_file, &pfk).unwrap();
    clear_file
}

#[derive(Serialize, Deserialize)]
#[wasm_bindgen(getter_with_clone)]
pub struct KeysAndPassword {
    pub api_password: String,
    pub encrypted_key: String,
    pub pub_key: String,
}

#[wasm_bindgen]
pub fn get_new_keys_and_password(user_password: String) -> KeysAndPassword {
    let (key_password, api_password) =
        key_management::get_key_and_api_password(user_password.as_bytes(), USER_SALT);
    let (secret_key, pub_key) = key_management::generate_key();
    let secret_key_pem = secret_key
        .to_pkcs8_encrypted_pem(&mut rand::rngs::OsRng, &key_password, LineEnding::default())
        .unwrap();
    KeysAndPassword {
        api_password: BASE64_STANDARD.encode(api_password),
        encrypted_key: secret_key_pem.to_string(),
        pub_key: pub_key.to_public_key_pem(LineEnding::default()).unwrap(),
    }
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

pub mod symetric_crypto {
    use aes_gcm::{aead::Aead, AeadCore, Aes256Gcm, KeyInit, Nonce};
    use p384::ecdh::SharedSecret;
    use rand::rngs::OsRng;
    use sha2::{
        digest::{generic_array::GenericArray, typenum, Key},
        Sha256,
    };

    use crate::{KEY_SALT, NONCE_SIZE};

    pub fn encrypt_file(file: &[u8], key: &Key<Aes256Gcm>) -> Vec<u8> {
        let nonce = Aes256Gcm::generate_nonce(&mut OsRng);
        let encrypted_file = encrypt_decrypt_message(file, &key, nonce, true).unwrap();

        let mut complete_ciphertext = nonce.to_vec();
        complete_ciphertext.extend(encrypted_file);
        complete_ciphertext
    }

    pub fn decrypt_file(file: &[u8], key: &Key<Aes256Gcm>) -> Result<Vec<u8>, aes_gcm::Error> {
        let nonce: GenericArray<u8, typenum::U12> = *GenericArray::from_slice(&file[..NONCE_SIZE]);
        let encrypted_part = &file[NONCE_SIZE..];
        encrypt_decrypt_message(encrypted_part, &key, nonce, false)
    }

    // Encrypts the PFK using the provided shared_secret
    pub fn encrypt_pfk(pfk: &[u8], shared_secret: &SharedSecret) -> Vec<u8> {
        let aes_key = compute_key_from_shared_secret(&shared_secret);
        let nonce = Aes256Gcm::generate_nonce(&mut OsRng);
        let ciphertext = encrypt_decrypt_message(pfk, &aes_key, nonce, true).unwrap();

        let mut complete_ciphertext = nonce.to_vec();
        complete_ciphertext.extend(ciphertext);
        complete_ciphertext
    }

    pub fn decrypt_pfk(pfk: &[u8], shared_secret: &SharedSecret) -> Key<Aes256Gcm> {
        let aes_key = compute_key_from_shared_secret(&shared_secret);
        let nonce: GenericArray<u8, typenum::U12> = *GenericArray::from_slice(&pfk[..NONCE_SIZE]);
        let encrypted_part = &pfk[NONCE_SIZE..];
        let clear_pfk: Vec<u8> =
            encrypt_decrypt_message(encrypted_part, &aes_key, nonce, false).unwrap();
        assert!(clear_pfk.len() == 32);
        *Key::<Aes256Gcm>::from_slice(&clear_pfk)
    }

    pub fn generate_aes_key() -> Key<Aes256Gcm> {
        Aes256Gcm::generate_key(&mut OsRng)
    }

    fn compute_key_from_shared_secret(shared_secret: &SharedSecret) -> Key<Aes256Gcm> {
        let hkdf = shared_secret.extract::<Sha256>(Some(KEY_SALT));
        let mut key = [0u8; 32];
        // TODO: Error handling
        hkdf.expand(&[], &mut key).unwrap();

        *Key::<Aes256Gcm>::from_slice(&key)
    }

    fn encrypt_decrypt_message(
        message: &[u8],
        key: &Key<Aes256Gcm>,
        nonce: Nonce<typenum::U12>,
        encrypt: bool,
    ) -> Result<Vec<u8>, aes_gcm::Error> {
        let cipher = Aes256Gcm::new(&key);
        if encrypt {
            cipher.encrypt(&nonce, message)
        } else {
            cipher.decrypt(&nonce, message)
        }
    }
    #[cfg(test)]
    mod tests {
        use p384::{elliptic_curve::ecdh, AffinePoint, NonZeroScalar, Scalar};
        // Note this useful idiom: importing names from outer (for mod tests) scope.
        use super::*;

        #[test]
        fn verify_encrypt_decrypt() {
            let message = "testmsg".as_bytes();
            let key = Key::<Aes256Gcm>::from_slice("abcddjguabcddjguabcddjguabcddjgu".as_bytes());
            let nonce = Aes256Gcm::generate_nonce(&mut OsRng);
            let encrypted_msg = encrypt_decrypt_message(message, key, nonce, true).unwrap();
            let clear_msg = encrypt_decrypt_message(&encrypted_msg, key, nonce, false).unwrap();
            assert_eq!(message, clear_msg);
        }

        #[test]
        fn verify_encrypt_decrypt_pfk() {
            let message = "testkey-testkey-testkey-testkey-".as_bytes();
            let pubkey = AffinePoint::GENERATOR;
            let secretkey = NonZeroScalar::new(Scalar::from_u64(12475)).unwrap();
            let shared_secret = ecdh::diffie_hellman(secretkey, pubkey);
            let shared_secret2 = ecdh::diffie_hellman(secretkey, pubkey);
            let encrypted_msg = encrypt_pfk(&message, &shared_secret);
            let clear_msg = decrypt_pfk(&encrypted_msg, &shared_secret2);

            assert_eq!(message, clear_msg.to_vec());
        }

        #[test]
        fn verify_encrypt_decrypt_file() {
            let key = Key::<Aes256Gcm>::from_slice("abcddjguabcddjguabcddjguabcddjgu".as_bytes());
            let file = "test-file_content".as_bytes();

            let encrypted_file = encrypt_file(file, key);
            let clear_file = decrypt_file(&encrypted_file, key).unwrap();
            assert_eq!(file, clear_file);
        }
    }
}

pub mod key_management {
    use p384::{
        pkcs8::{DecodePrivateKey, EncodePrivateKey},
        PublicKey, SecretKey,
    };
    use pbkdf2;
    use pkcs8::{der::zeroize::Zeroizing, DecodePublicKey, EncodePublicKey, LineEnding};
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

    /// Return a new pair of Nist P384 keys
    pub fn generate_key() -> (SecretKey, PublicKey) {
        let secret_key = SecretKey::random(&mut OsRng);
        let pub_key = secret_key.public_key();
        (secret_key, pub_key)
    }

    pub fn encrypt_key(secret_key: SecretKey, key_password: &[u8]) -> Zeroizing<String> {
        // TODO: Remove unwrap and handle errors
        secret_key
            .to_pkcs8_encrypted_pem(&mut OsRng, key_password, LineEnding::default())
            .unwrap()
    }

    pub fn decrypt_private_key(
        encrypted_secret_key: Zeroizing<String>,
        key_password: &[u8],
    ) -> SecretKey {
        let secret_key =
            SecretKey::from_pkcs8_encrypted_pem(&encrypted_secret_key, key_password).unwrap();
        secret_key
    }

    pub fn encode_public_key(public_key: PublicKey) -> String {
        public_key.to_public_key_pem(LineEnding::default()).unwrap()
    }

    pub fn decode_public_key(encoded_public_key: &String) -> PublicKey {
        PublicKey::from_public_key_pem(&encoded_public_key).unwrap()
    }
}

pub mod decryption {
    use p384::{ecdh, ecdh::SharedSecret, PublicKey, SecretKey};

    pub fn generate_shared_secret(secret_key: SecretKey, public_key: PublicKey) -> SharedSecret {
        ecdh::diffie_hellman(secret_key.to_nonzero_scalar(), public_key.as_affine())
    }
}
