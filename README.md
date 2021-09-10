# Real-Blockchain-In-Golang

## c1_genesis_json
The blockchain is an **IMMUTABLE** database. The token supply, initial user balances and global blockchain settings are defined in a Genesis file.

## c2_db_changes.txt
The Genesis balances indicate what was the original blockchain state and are never updated afterwards.
The database state changes are called **Transactions (TX)**.

## c3_state_blockchain_component
Transactions are old fashion Events representing actions within the system.

## c4_caesar_transfer
Closed software with centralized access to private data and rules puts only a few people to the position fo power.
User don't have a choice, and shareholders are in business to make money.

## c5_broken_trust
Blockchain developers aim to develop protocols where applications' entrepreneurs and users synergize in a transparent, auditable 
relationship. Specifications fo the blockchain system should be well-defined from the beginning and only change if its users
support it.

## c6_immutable_hash
The database content is hashed by secure cryptographic hash function. The blockchain participants use the resulted hash
to reference a specific database state.

## c7_blockchain_programming_model
Transactions are grouped into batches for performance reasons. A batch of transactions make a Block. Each block is encoded and hashed
using a secure, cryptographic hash function.
Block contains *Header* and *Payload*. The Header stores various metadata such a time and a reference to the *Parent Block* (the previous
immutable database state). The Payload carries the new database transactions.
