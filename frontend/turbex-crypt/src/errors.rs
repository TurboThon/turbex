pub mod turbex_errors {
    use std::fmt::{self};

    use wasm_bindgen::JsValue;


    #[derive(Debug)]
    pub enum TurbexError {
        PublicKeyDecodingError,
        PublicKeyEncodingError,
        PrivateKeyDecryptionError,
        PrivateKeyEncryptionError,
        DHSecretHKDFExpandError,
        IncorrectPFKLengthError,
        PFKDecryptionError,
        PFKEncryptionError,
        EphemeralKeyEncodingError,
        FileDecryptionError,
        FileEncryptionError,
        Base64DecodingError
    }


    impl fmt::Display for TurbexError {
        fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
            write!(f, "{:?}", self)
        }
    }

    impl Into<JsValue> for TurbexError {
        fn into(self) -> JsValue {
            JsValue::from_str(&self.to_string())
        }
    }

    impl From<base64::DecodeError> for TurbexError {
        fn from(_value: base64::DecodeError) -> Self {
            TurbexError::Base64DecodingError
        }
    }

}
