// Copyright (c) 2022 The Jaeger Authors.
// SPDX-License-Identifier: Apache-2.0

package tlscfg

import (
	"crypto/tls"
	"fmt"
)

// https://pkg.go.dev/crypto/tls#pkg-constants
var versions = map[string]uint16{
	"1.0": tls.VersionTLS10,
	"1.1": tls.VersionTLS11,
	"1.2": tls.VersionTLS12,
	"1.3": tls.VersionTLS13,
}

func allCiphers() map[string]uint16 {
	acceptedCiphers := make(map[string]uint16)
	for _, suite := range tls.CipherSuites() {
		acceptedCiphers[suite.Name] = suite.ID
	}
	return acceptedCiphers
}

// CipherSuiteNamesToIDs returns a list of cipher suite IDs from the cipher suite names passed.
func CipherSuiteNamesToIDs(cipherNames []string) ([]uint16, error) {
	var ciphersIDs []uint16
	possibleCiphers := allCiphers()
	for _, cipher := range cipherNames {
		intValue, ok := possibleCiphers[cipher]
		if !ok {
			return nil, fmt.Errorf("cipher suite %s not supported or doesn't exist", cipher)
		}
		ciphersIDs = append(ciphersIDs, intValue)
	}
	return ciphersIDs, nil
}

// VersionNameToID returns the version ID from version name
func VersionNameToID(versionName string) (uint16, error) {
	if version, ok := versions[versionName]; ok {
		return version, nil
	}
	return 0, fmt.Errorf("unknown tls version %q", versionName)
}
