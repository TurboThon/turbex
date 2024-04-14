# Technical Architecture Document

This document and all resources aims at describing the objectives of the project
and its technical architecture.

### Runtime architecture

The runtime architecture is described by the following communication diagram.

```mermaid
stateDiagram-v2
    state Browser {
        Svelte: Svelte static JS
        Svelte --> WASM : Function calls
        WASM: Rust WASM
        note right of WASM: Cryptographic functions
        note left of Svelte: Presentation layer
    }
    Svelte --> Nginx: HTTP requests
    state Server {
        Nginx: Nginx frontend files and proxy
        note left of Nginx
            Acts as proxy for backend and serve frontend
            compiled files (Svelte and Rust)
        end note
        Nginx --> Nginx: Frontend files
        Nginx --> Back: HTTP requests to /api
        Back: Go Backend
        Back --> Mongo: Users, keys, files...
        Mongo: Mongo database
        note right of Mongo
            Files are stored in a gridfs bucket, optimised for object storage
            Other data are stored in collections (~ tables)
        end note
    }
```

### Authentication flow

```mermaid
flowchart TD
	subgraph CLIENT
	A[User password] -->|PBKDF2| B[APIPassword + KeyPassword]
	B -->|KeyPassword stays in the client| CKEY
	F[Encrypted ECDH key] -->|Client side decryption| CKEY[Clear private ECDH key]
	end
	subgraph SRV
	B -->|APIPassword sent to web API| D[APIPassword]
	D[APIPassword] -->|PBKDF2| E[Compared with DB]
	end
	E -->|Sends the user's private key| F
```

### Encryption flow

The intended flow is a Diffie-Hellman Key Exchange using static-ephemeral keys.
The sender computes an ephemeral key whereas the receiver has a static key (prior to the exchange).

This is the basis of IES (Integrated Encryption Scheme). Since we are using ECC, we implement ECIES with p384, which is recommended by ANSSI and NIST.

```mermaid
flowchart TD
    SEK[Sender Ephemeral key] --> PFKCompute
    SEK --> SEPUBK[Sender Ephemeral public key]
    SEPUBK -->|Sent| SRV[web server]
    PFK[File encryption key] -->|AES-GCM| EF[Encrypted file]
    EF -->|Sent| SRV
    BPUBK[Bob public key] --> PFKCompute
    PFKCompute(ECDH Key agreement) -->|First operation of ECDH| SS
	SS[Shared Secret] -->|AES-GCM| EPFK[Encrypted file key]
    EPFK -->|Sent| SRV
```

### Decryption flow

```mermaid
flowchart TD
	SEPUBK[Sender Ephemeral public key] --> SSCompute
	BPUBK[Bob private key] --> SSCompute
	SSCompute(ECDH Key agreement) -->|Second operation of ECDH| SS
	SS[Shared Secret] --> PFKCompute(File key decryption)

    EPFK[Encrypted file key] --> PFKCompute
    PFKCompute -->|AES-GCM| PFK
    PFK --> FCompute(File decryption)
    EF[Encrypted file] --> FCompute
    FCompute -->|AES-GCM| F[File]
```

