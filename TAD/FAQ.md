# FAQ

#### What exactly is stored on the server ?

Encrypted files, public keys, encrypted private keys, hashed passwords and ephemeral keys.
Turbex uses your password to decrypt your private key which grant you access to all files shared with you.
Your password is never sent to the server, therefore the server does not have enough information to decrypt any file.

#### Why do I need to create a strong password ?

Cryptographic operations are only as strong as your password is, and your password is the only private information needed
to access to files.

#### I lost my password, can I recover my account ?

Unfortunately, this is over for your account. There are no means to recover it from the application.
You can ask the administrator to delete your account so that cleaning routine can delete your public key and files
shared with your account since they must remain confidential.

#### Can I change my username ?

Currently we cannot allow people to change username because of technical architecture.
As of now we can only suggest you to use a random username if you need to remain anonymous.
