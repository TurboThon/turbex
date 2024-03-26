// Convert a string to a buffer (Uint8Array)
function strToBuffer (string) {
  let arrayBuffer = new ArrayBuffer(string.length * 1);
  let newUint = new Uint8Array(arrayBuffer);
  newUint.forEach((_, i) => {
    newUint[i] = string.charCodeAt(i);
  });
  return newUint;
}

// Returns true if elements of arr1 are equals to elements of arr2
// false otherwise
function isEqual(arr1, arr2) {
  if (arr1.length !== arr2.length) {
    return false
  }
  return arr1.every((value, index) => value === arr2[index])
}

// The following user stories should always be possible

// The users encrypts a file and send it to itself
passwords = turbex.get_api_password_and_key("abcd");
userkeys = turbex.create_user_key(passwords.key_password);
file = strToBuffer("abcdenfg");
aes_key = turbex.generate_aes_key();
encryptedfile = turbex.encrypt_file(file, aes_key);
sfk = turbex.encrypt_pfk(aes_key, userkeys.public_key);
clearfile = turbex.decrypt_file(encryptedfile, sfk.encrypted_pfk, sfk.ephemeral_pub_key, userkeys.private_key, passwords.key_password);

console.assert(isEqual(file, clearfile), "User is not able to encrypt a file for itself.")
console.log("TEST1: Passed (User is able to encrypt file to itself)")

// UserA encrypts a file and send it to UserB and UserC
passwordsB = turbex.get_api_password_and_key("userB password");
userkeysB = turbex.create_user_key(passwordsB.key_password);

passwordsC = turbex.get_api_password_and_key("userC password");
userkeysC = turbex.create_user_key(passwordsB.key_password);

file = strToBuffer("Content of fileA sent to userB");
aes_key = turbex.generate_aes_key();
encryptedfile = turbex.encrypt_file(file, aes_key);
sfkB = turbex.encrypt_pfk(aes_key, userkeysB.public_key);
sfkC = turbex.encrypt_pfk(aes_key, userkeysC.public_key);

clearfileB = turbex.decrypt_file(encryptedfile, sfkB.encrypted_pfk, sfkB.ephemeral_pub_key, userkeysB.private_key, passwordsB.key_password);
clearfileC = turbex.decrypt_file(encryptedfile, sfkB.encrypted_pfk, sfkB.ephemeral_pub_key, userkeysB.private_key, passwordsB.key_password);
console.assert(isEqual(file, clearfileB), "UserB does not get the correct file from UserA.")
console.assert(isEqual(file, clearfileC), "UserC does not get the correct file from UserA.")
console.log("TEST2: Passed (UserA is able to send a file to 2 Users)")

