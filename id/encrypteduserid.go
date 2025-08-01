// Copyright (c) 2021 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package id

import (
	"fmt"
	"regexp"
)

// EncryptedUserID represents a Matrix user ID encrypted with RSA-OAEP
type EncryptedUserID string

// TODO: modify this number was only chosen because it was greater than the RSA key length
const EncryptedUserIDMaxLength = 750

func NewEncryptedUserID(localpart, homeserver string) EncryptedUserID {
	return EncryptedUserID(fmt.Sprintf("@%s:%s", localpart, homeserver))
}

// Parse parses the user ID into the localpart and server name.
//
// Note that this only enforces very basic user ID formatting requirements: user IDs start with
// a @, and contain a : after the @. If you want to enforce localpart validity, see the
// ParseAndValidate and ValidateUserLocalpart functions.
func (encryptedUserID EncryptedUserID) Parse() (localpart, homeserver string, err error) {
	var sigil byte
	sigil, localpart, homeserver = ParseCommonIdentifier(encryptedUserID)
	if sigil != '@' || homeserver == "" {
		err = fmt.Errorf("'%s' %w", encryptedUserID, ErrInvalidUserID)
	}
	return
}

func (encryptedUserID EncryptedUserID) Localpart() string {
	localpart, _, _ := encryptedUserID.Parse()
	return localpart
}

func (encryptedUserID EncryptedUserID) Homeserver() string {
	_, homeserver, _ := encryptedUserID.Parse()
	return homeserver
}

// URI returns the user ID as a MatrixURI struct, which can then be stringified into a matrix: URI or a matrix.to URL.
//
// This does not parse or validate the user ID. Use the ParseAndValidate method if you want to ensure the user ID is valid first.
func (encryptedUserID EncryptedUserID) URI() *MatrixURI {
	if encryptedUserID == "" {
		return nil
	}
	return &MatrixURI{
		Sigil1: '@',
		MXID1:  string(encryptedUserID)[1:],
	}
}

var ValidEncryptedLocalpartRegex = regexp.MustCompile("^[0-9a-zA-Z-.=_/+]+$")

// ValidateUserLocalpart validates that an Encrypted Matrix user ID localpart only contains base64 characters
func ValidateEncryptedUserLocalpart(localpart string) error {
	if len(localpart) == 0 {
		return ErrEmptyLocalpart
	} else if !ValidEncryptedLocalpartRegex.MatchString(localpart) {
		return fmt.Errorf("'%s' %w", localpart, ErrNoncompliantLocalpart)
	}
	return nil
}

// ParseAndValidate parses the encrypted user ID into the localpart and server name like Parse,
// and also validates that the localpart is allowed according to the user identifiers spec.
func (encryptedUserID EncryptedUserID) ParseAndValidate() (localpart, homeserver string, err error) {
	localpart, homeserver, err = encryptedUserID.Parse()
	if err == nil {
		err = ValidateEncryptedUserLocalpart(localpart)
	}
	if err == nil && len(encryptedUserID) > EncryptedUserIDMaxLength {
		err = ErrUserIDTooLong
	}
	if err == nil && !ValidateServerName(homeserver) {
		err = fmt.Errorf("%q %q", homeserver, ErrNoncompliantServerPart)
	}
	return
}

func (encryptedUserID EncryptedUserID) String() string {
	return string(encryptedUserID)
}
