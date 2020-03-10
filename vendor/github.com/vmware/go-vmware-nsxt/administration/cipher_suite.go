/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package administration

// HTTP cipher suite
type CipherSuite struct {

	// Enable status for this cipher suite
	Enabled bool `json:"enabled"`

	// Name of the TLS cipher suite
	Name string `json:"name"`
}