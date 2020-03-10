/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package manager

// MAC learning configuration
type MacLearningSpec struct {

	// Aging time in sec for learned MAC address
	AgingTime int32 `json:"aging_time,omitempty"`

	// Allowing source MAC address learning
	Enabled bool `json:"enabled"`

	// The maximum number of MAC addresses that can be learned on this port
	Limit int32 `json:"limit,omitempty"`

	// The policy after MAC Limit is exceeded
	LimitPolicy string `json:"limit_policy,omitempty"`

	// Allowing flooding for unlearned MAC for ingress traffic
	UnicastFloodingAllowed bool `json:"unicast_flooding_allowed"`
}