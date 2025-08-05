# mautrixgo - Research Project Modifications

This repository is a modified fork of the original mautrixgo that can be found [here](https://github.com/mautrix/go). The code here has been modified to serve as a component in the "Achieving User Pseudonymity in Federated End-to-End Encrypted Messaging Platforms" research project. For an overview of the project, see the [main project README](https://github.com/Ijnek117/anonymous-matrix). 

---

## 1. Introduction

mautrixgo is a Go Matrix framework used by gomuks. 

## 2. Summary of Changes

The following changes were made to this repository for the research project:

* Introduced the `EncryptedUserID` type for an RSA-OAEP encrypted User ID and modifying/created structs and functions to support it.
* Added support to fetch a remote server's TLS certificate.
* Modified the User Invite sequenceto store the mapping between the invited User ID and their SenderID. (May no longer be necessary)

## 3. Rationale for Modifications

The new EncryptedUserID type is required as once encrypted the localpart of the User ID exceeds the length limit and doesn't follow the UserID's format.

## 4. Key Files Modified

This table provides a more detailed reference to the files and code that were changed or added.

| File Path | Change Description |
| :--- | :--- |
| `client.go` | Added `InviteUserWithResp` to invite a user and receive a `sender_key` (Sender ID) in the response. Also added `EncryptUser`, which fetches a homeserver's TLS public key to encrypt a User ID's localpart using RSA-OAEP. Certificate verification must be added. |
| `id/encrypteduserid.go` | **(New File)** Defines the `EncryptedUserID` type for an RSA-OAEP encrypted User ID, along with parsing and validation functions. |
| `requests.go` | Added the `ReqInviteEncryptedUser` struct to support user invites with `EncryptedUserID`. |
| `responses.go` | Added `InviteUserResp` to handle responses containing a `sender_key` and `ServerTLSCertResponse` for the new TLS key fetching endpoint. |
| `statestore.go` | Added the `SetPseudoMapping` method to the `StateStore` interface and `MemoryStateStore` implementation to map a `UserID` to a room-specific `SenderID`. |
