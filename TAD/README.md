# Technical Architecture Document

This document and all resources aims at describing the objectives of the project
and its technical architecture.

### Authentication flow

```mermaid
graph TD
	subgraph CLIENT
	A[MDP Utilisateur] -->|PBKDF2| B[APIPassword + ECDSAKeyPassword]
	B -->|ECDSAKeyPassword| CKEY
	F[PrivKey ECDSA chiffrée] -->|Déchiffrement côté client| CKEY[PrivKey ECDSA clair]
	end
	subgraph SRV
	B -->|Envoi à l'API web| D[APIPassword]
	D[APIPassword] -->|PBKDF2| E[Comparaison DB]
	end
	E -->|Envoi de la clé privée| F
```

### Encryption flow

The intended flow is a Diffie-Hellman Key Exchange using static-ephemeral keys.
The sender computes an ephemeral key whereas the receiver has a static key (prior to the exchange).

This is the basis of IES (Integrated Encryption Scheme). Since we are using ECC, we implement ECIES with p384, which is recommended by ANSSI and NIST.

```mermaid
graph TD
    SEK[Sender Ephemeral key] --> PFKCompute
    SEK --> SEPUBK[Sender Ephemeral public key]
    SEPUBK -->|Envoi| SRV
    BPUBK[Clé publique Bob] --> PFKCompute
    PFKCompute(Dérivation ECDH) -->|Dérivation ECDH| PFK
	PFK[Shared Secret] -->|AES-GCM| EF[Fichier chiffré]
	EF -->|Envoi| SRV
```

### Decryption flow

```mermaid
graph TD
	SEPUBK[Sender Ephemeral public key] --> PFKCompute
	BPUBK[Clé privée Bob] --> PFKCompute
	PFKCompute(Dérivation ECDH) -->|Dérivation ECDH| PFK
	PFK[Shared Secret] -->|AES-GCM| EF[Fichier déchiffré]
```

