# Distributed Systems Coursework
* By: Bekir Tekmen, Sigurd K. Würtz, Jonas Nielsen and Martin Frederiksen
* Course: Distributed Systems
* Aarhus University fall 2021

This project is a part of the Distributed Systems coursework.
It implements a peer-to-peer network with a blockchain ledger and RSA encryption for secure transactions.
The main functionalities of the project include:

Peer-to-Peer Network: The network is created by peers connecting to each other. Each peer maintains a list of other peers in the network.
Account Ledger: Manages account creation, balance management, and transaction handling.
Blockchain: The blockchain is used to store transactions securely. Each block contains a list of transaction IDs and is signed using RSA.
RSA Encryption: RSA is used to encrypt and decrypt messages, as well as to sign and verify transactions and blocks.

## Project Files

- **peer**: Contains the main implementation of the peer-to-peer network.
  - [`peer.go`](src\peer\peer.go): Main file implementing the Peer struct and its methods.
- **ownRSA**: Contains the implementation of RSA encryption and decryption.
  - [`ownRSA.go`](src\peer\ownRSA/ownRSA.go): Implements RSA key generation, encryption, decryption, and signature verification.
- **utility**: Contains utility functions used across the project.
  - [`utility.go`](src\peer\utility\utility.go): Implements helper functions for the project.

### Test Files

- **peer**: Contains tests for the peer-to-peer network.
  - [`peer_test.go`](src\peer\peer_test.go): Contains tests for the peer-to-peer network.
- **ownRSA**: Contains tests for the RSA encryption and decryption.
  - [`ownRSA_test.go`](src\peer\ownRSA\ownRSA_test.go): Contains tests for the RSA encryption and decryption.
- **utility**: Contains tests for the utility functions.
  - [`utility_test.go`](src\utility\utility_test.go): Contains tests for the utility functions.

### Project Details

- **Peer-to-Peer Network**: The network is created by peers connecting to each other. Each peer maintains a list of other peers in the network.
- **Account Ledger**: Manages account creation, balance management, and transaction handling.
- **Blockchain**: The blockchain is used to store transactions securely. Each block contains a list of transaction IDs and is signed using RSA.
- **RSA Encryption**: RSA is used to encrypt and decrypt messages, as well as to sign and verify transactions and blocks.
