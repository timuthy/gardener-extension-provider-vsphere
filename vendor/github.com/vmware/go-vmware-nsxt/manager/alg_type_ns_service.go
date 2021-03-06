/* Copyright © 2017 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause

   Generated by: https://github.com/swagger-api/swagger-codegen.git */

package manager

type AlgTypeNsServiceEntry struct {
	ResourceType string `json:"resource_type"`

	Alg string `json:"alg"`

	// The destination_port cannot be empty and must be a single value.
	DestinationPorts []string `json:"destination_ports,omitempty"`

	SourcePorts []string `json:"source_ports,omitempty"`
}

type AlgTypeNsService struct {
	NsService

	NsserviceElement AlgTypeNsServiceEntry `json:"nsservice_element"`
}
