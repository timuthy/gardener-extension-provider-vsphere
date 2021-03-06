/*
 * Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/ip_pools"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/sites/enforcement_points"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s"
	t1nat "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s/nat"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"

	api "github.com/gardener/gardener-extension-provider-vsphere/pkg/apis/vsphere"
	vinfra "github.com/gardener/gardener-extension-provider-vsphere/pkg/vsphere/infrastructure"
)

type lookupTier0GatewayTask struct{ baseTask }

func NewLookupTier0GatewayTask() Task {
	return &lookupTier0GatewayTask{baseTask{label: "tier-0 gateway lookup"}}
}

func (t *lookupTier0GatewayTask) NameToLog(spec vinfra.NSXTInfraSpec) *string {
	return &spec.Tier0GatewayName
}

func (t *lookupTier0GatewayTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.Tier0GatewayRef
}

func (t *lookupTier0GatewayTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	name := spec.Tier0GatewayName
	client := infra.NewDefaultTier0sClient(ctx.Connector())
	var cursor *string
	total := 0
	count := 0
	for {
		result, err := client.List(cursor, nil, nil, nil, nil, nil)
		if err != nil {
			return "", nicerVAPIError(err)
		}
		for _, item := range result.Results {
			if *item.DisplayName == name {
				// found
				state.Tier0GatewayRef = &api.Reference{ID: *item.Id, Path: *item.Path}
				return actionFound, nil
			}
		}
		if cursor == nil {
			total = int(*result.ResultCount)
		}
		count += len(result.Results)
		if count >= total {
			return "", fmt.Errorf("not found: %s", name)
		}
		cursor = result.Cursor
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type lookupEdgeClusterTask struct{ baseTask }

func NewLookupEdgeClusterTask() Task {
	return &lookupEdgeClusterTask{baseTask{label: "edge cluster lookup"}}
}

func (t *lookupEdgeClusterTask) NameToLog(spec vinfra.NSXTInfraSpec) *string {
	return &spec.EdgeClusterName
}

func (t *lookupEdgeClusterTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.EdgeClusterRef
}

func (t *lookupEdgeClusterTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	name := spec.EdgeClusterName
	client := enforcement_points.NewDefaultEdgeClustersClient(ctx.Connector())
	result, err := client.List(defaultSite, policyEnforcementPoint, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return "", nicerVAPIError(err)
	}
	for _, item := range result.Results {
		if *item.DisplayName == name {
			state.EdgeClusterRef = &api.Reference{ID: *item.Id, Path: *item.Path}
			return actionFound, nil
		}
	}
	return "", fmt.Errorf("not found: %s", name)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type lookupTransportZoneTask struct{ baseTask }

func NewLookupTransportZoneTask() Task {
	return &lookupTransportZoneTask{baseTask{label: "transport zone lookup"}}
}

func (t *lookupTransportZoneTask) NameToLog(spec vinfra.NSXTInfraSpec) *string {
	return &spec.TransportZoneName
}

func (t *lookupTransportZoneTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.TransportZoneRef
}

func (t *lookupTransportZoneTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	name := spec.TransportZoneName
	client := enforcement_points.NewDefaultTransportZonesClient(ctx.Connector())
	result, err := client.List(defaultSite, policyEnforcementPoint, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return "", nicerVAPIError(err)
	}
	for _, item := range result.Results {
		if *item.DisplayName == name {
			state.TransportZoneRef = &api.Reference{ID: *item.Id, Path: *item.Path}
			return actionFound, nil
		}
	}
	return "", fmt.Errorf("not found: %s", name)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type lookupSNATIPPoolTask struct{ baseTask }

func NewLookupSNATIPPoolTask() Task {
	return &lookupSNATIPPoolTask{baseTask{label: "SNAT IP pool lookup"}}
}

func (t *lookupSNATIPPoolTask) NameToLog(spec vinfra.NSXTInfraSpec) *string {
	return &spec.SNATIPPoolName
}

func (t *lookupSNATIPPoolTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.SNATIPPoolRef
}

func (t *lookupSNATIPPoolTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	name := spec.SNATIPPoolName
	client := infra.NewDefaultIpPoolsClient(ctx.Connector())
	var cursor *string
	total := 0
	count := 0
	for {
		result, err := client.List(cursor, nil, nil, nil, nil, nil)
		if err != nil {
			return "", nicerVAPIError(err)
		}
		for _, item := range result.Results {
			if *item.DisplayName == name {
				// found
				state.SNATIPPoolRef = &api.Reference{ID: *item.Id, Path: *item.Path}
				return actionFound, nil
			}
		}
		if cursor == nil {
			total = int(*result.ResultCount)
		}
		count += len(result.Results)
		if count >= total {
			return "", fmt.Errorf("not found: %s", name)
		}
		cursor = result.Cursor
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type tier1GatewayTask struct{ baseTask }

func NewTier1GatewayTask() Task {
	return &tier1GatewayTask{baseTask{label: "tier-1 gateway"}}
}

func (t *tier1GatewayTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.Tier1GatewayRef
}

func (t *tier1GatewayTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	client := infra.NewDefaultTier1sClient(ctx.Connector())

	tier1 := model.Tier1{
		DisplayName:  strptr(spec.FullClusterName()),
		Description:  strptr(description),
		FailoverMode: strptr(model.Tier1_FAILOVER_MODE_PREEMPTIVE),
		Tags:         spec.CreateTags(),
		RouteAdvertisementTypes: []string{
			model.Tier1_ROUTE_ADVERTISEMENT_TYPES_STATIC_ROUTES,
			model.Tier1_ROUTE_ADVERTISEMENT_TYPES_NAT,
			model.Tier1_ROUTE_ADVERTISEMENT_TYPES_LB_VIP,
			model.Tier1_ROUTE_ADVERTISEMENT_TYPES_LB_SNAT,
		},
		Tier0Path: &state.Tier0GatewayRef.Path,
	}

	if state.Tier1GatewayRef != nil {
		oldTier1, err := client.Get(state.Tier1GatewayRef.ID)
		if isNotFoundError(err) {
			state.Tier1GatewayRef = nil
			return t.Ensure(ctx, spec, state)
		}
		if err != nil {
			return readingErr(err)
		}
		if *oldTier1.DisplayName != *tier1.DisplayName ||
			oldTier1.FailoverMode == nil ||
			*oldTier1.FailoverMode != *tier1.FailoverMode ||
			oldTier1.Tier0Path == nil ||
			*oldTier1.Tier0Path != *tier1.Tier0Path ||
			!equalStrings(oldTier1.RouteAdvertisementTypes, tier1.RouteAdvertisementTypes) ||
			!equalTags(oldTier1.Tags, tier1.Tags) {
			err := client.Patch(state.Tier1GatewayRef.ID, tier1)
			if err != nil {
				return updatingErr(err)
			}
			return actionUpdated, nil
		}
		return actionUnchanged, nil
	}

	id := generateID("tier1gw")
	createdObj, err := client.Update(id, tier1)
	if err != nil {
		return creatingErr(err)
	}
	state.Tier1GatewayRef = &api.Reference{ID: *createdObj.Id, Path: *createdObj.Path}
	return actionCreated, nil
}

func (t *tier1GatewayTask) SetRecoveredReference(state *api.NSXTInfraState, ref *api.Reference, _ *string) {
	state.Tier1GatewayRef = ref
}

func (t *tier1GatewayTask) ListAll(ctx EnsurerContext, _ *api.NSXTInfraState, cursor *string) (interface{}, error) {
	client := infra.NewDefaultTier1sClient(ctx.Connector())
	return client.List(cursor, nil, nil, nil, nil, nil)
}

func (t *tier1GatewayTask) EnsureDeleted(ctx EnsurerContext, state *api.NSXTInfraState) (bool, error) {
	client := infra.NewDefaultTier1sClient(ctx.Connector())
	if state.Tier1GatewayRef == nil {
		return false, nil
	}
	err := client.Delete(state.Tier1GatewayRef.ID)
	if err != nil {
		return false, nicerVAPIError(err)
	}
	state.Tier1GatewayRef = nil
	return true, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type tier1GatewayLocaleServiceTask struct{ baseTask }

func NewTier1GatewayLocaleServiceTask() Task {
	return &tier1GatewayLocaleServiceTask{baseTask{label: "tier-1 gateway local service"}}
}

func (t *tier1GatewayLocaleServiceTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.LocaleServiceRef
}

func (t *tier1GatewayLocaleServiceTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	client := tier_1s.NewDefaultLocaleServicesClient(ctx.Connector())

	obj := model.LocaleServices{
		DisplayName:     strptr(spec.FullClusterName()),
		Description:     strptr(description),
		EdgeClusterPath: &state.EdgeClusterRef.Path,
		Tags:            spec.CreateTags(),
	}

	if state.LocaleServiceRef != nil {
		oldTier1, err := client.Get(state.LocaleServiceRef.ID, defaultPolicyLocaleServiceID)
		if isNotFoundError(err) {
			state.Tier1GatewayRef = nil
			return t.Ensure(ctx, spec, state)
		}
		if err != nil {
			return readingErr(err)
		}
		if *oldTier1.DisplayName != *obj.DisplayName ||
			oldTier1.EdgeClusterPath == nil ||
			*oldTier1.EdgeClusterPath != *obj.EdgeClusterPath ||
			!equalTags(oldTier1.Tags, obj.Tags) {
			err := client.Patch(state.LocaleServiceRef.ID, defaultPolicyLocaleServiceID, obj)
			if err != nil {
				return updatingErr(err)
			}
			return actionUpdated, nil
		}
		return actionUnchanged, nil
	}

	// The default ID of the locale service will be the Tier1 ID
	id := state.Tier1GatewayRef.ID
	err := client.Patch(id, defaultPolicyLocaleServiceID, obj)
	if err != nil {
		return creatingErr(err)
	}
	state.LocaleServiceRef = &api.Reference{ID: id, Path: ""}
	return actionCreated, nil
}

func (t *tier1GatewayLocaleServiceTask) SetRecoveredReference(state *api.NSXTInfraState, _ *api.Reference, _ *string) {
	state.LocaleServiceRef = &api.Reference{ID: state.Tier1GatewayRef.ID, Path: ""}
}

func (t *tier1GatewayLocaleServiceTask) ListAll(ctx EnsurerContext, state *api.NSXTInfraState, cursor *string) (interface{}, error) {
	client := tier_1s.NewDefaultLocaleServicesClient(ctx.Connector())
	return client.List(state.Tier1GatewayRef.ID, cursor, nil, nil, nil, nil, nil)
}

func (t *tier1GatewayLocaleServiceTask) EnsureDeleted(ctx EnsurerContext, state *api.NSXTInfraState) (bool, error) {
	client := tier_1s.NewDefaultLocaleServicesClient(ctx.Connector())
	if state.LocaleServiceRef == nil {
		return false, nil
	}
	err := client.Delete(state.LocaleServiceRef.ID, defaultPolicyLocaleServiceID)
	if err != nil {
		return false, nicerVAPIError(err)
	}
	state.LocaleServiceRef = nil
	return true, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type segmentTask struct{ baseTask }

func NewSegmentTask() Task {
	return &segmentTask{baseTask{label: "segment"}}
}

func (t *segmentTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.SegmentRef
}

func (t *segmentTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	client := infra.NewDefaultSegmentsClient(ctx.Connector())

	gatewayAddr, err := cidrHostAndPrefix(spec.WorkersNetwork, 1)
	if err != nil {
		return "", errors.Wrapf(err, "gateway address")
	}
	subnet := model.SegmentSubnet{
		GatewayAddress: strptr(gatewayAddr),
	}
	displayName := spec.FullClusterName() + "-" + RandomString(8)
	segment := model.Segment{
		DisplayName:       strptr(displayName),
		Description:       strptr(description),
		ConnectivityPath:  strptr(state.Tier1GatewayRef.Path),
		TransportZonePath: strptr(state.TransportZoneRef.Path),
		Tags:              spec.CreateTags(),
		Subnets:           []model.SegmentSubnet{subnet},
	}

	if state.SegmentRef != nil {
		oldSegment, err := client.Get(state.SegmentRef.ID)
		if isNotFoundError(err) {
			state.SegmentRef = nil
			return t.Ensure(ctx, spec, state)
		}
		if err != nil {
			return readingErr(err)
		}
		if !strings.HasPrefix(*oldSegment.DisplayName, spec.FullClusterName()) ||
			oldSegment.ConnectivityPath == nil ||
			*oldSegment.ConnectivityPath != *segment.ConnectivityPath ||
			oldSegment.TransportZonePath == nil ||
			*oldSegment.TransportZonePath != *segment.TransportZonePath ||
			len(oldSegment.Subnets) != 1 ||
			oldSegment.Subnets[0].GatewayAddress == nil ||
			*oldSegment.Subnets[0].GatewayAddress != *segment.Subnets[0].GatewayAddress ||
			!equalTags(oldSegment.Tags, segment.Tags) {
			err := client.Patch(state.SegmentRef.ID, segment)
			if err != nil {
				return updatingErr(err)
			}
			return actionUpdated, nil
		}
		return actionUnchanged, nil
	}

	id := generateID("segment")
	createdObj, err := client.Update(id, segment)
	if err != nil {
		return creatingErr(err)
	}
	state.SegmentRef = &api.Reference{ID: *createdObj.Id, Path: *createdObj.Path}
	state.SegmentName = createdObj.DisplayName
	return actionCreated, nil
}

func (t *segmentTask) SetRecoveredReference(state *api.NSXTInfraState, ref *api.Reference, name *string) {
	state.SegmentRef = ref
	state.SegmentName = name
}

func (t *segmentTask) ListAll(ctx EnsurerContext, _ *api.NSXTInfraState, cursor *string) (interface{}, error) {
	client := infra.NewDefaultSegmentsClient(ctx.Connector())
	return client.List(cursor, nil, nil, nil, nil, nil)
}

func (t *segmentTask) EnsureDeleted(ctx EnsurerContext, state *api.NSXTInfraState) (bool, error) {
	client := infra.NewDefaultSegmentsClient(ctx.Connector())
	if state.SegmentRef == nil {
		return false, nil
	}
	err := client.Delete(state.SegmentRef.ID)
	if err != nil {
		return false, nicerVAPIError(err)
	}
	state.SegmentRef = nil
	state.SegmentName = nil
	return true, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type snatIPAddressAllocationTask struct{ baseTask }

func NewSNATIPAddressAllocationTask() Task {
	return &snatIPAddressAllocationTask{baseTask{label: "SNAT IP address allocation"}}
}

func (t *snatIPAddressAllocationTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.SNATIPAddressAllocRef
}

func (t *snatIPAddressAllocationTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	client := ip_pools.NewDefaultIpAllocationsClient(ctx.Connector())

	allocation := model.IpAddressAllocation{
		DisplayName: strptr(spec.FullClusterName() + "_SNAT"),
		Description: strptr("SNAT IP address for all nodes. " + description),
		Tags:        spec.CreateTags(),
	}

	if state.SNATIPAddressAllocRef != nil {
		_, err := client.Get(state.SNATIPPoolRef.ID, state.SNATIPAddressAllocRef.ID)
		if err == nil {
			// IP address allocation is never updated
			return actionUnchanged, nil
		}
		if !isNotFoundError(err) {
			return readingErr(err)
		}
	}

	id := generateID("snatippool")
	createdObj, err := client.Update(state.SNATIPPoolRef.ID, id, allocation)
	if err != nil {
		return creatingErr(err)
	}
	state.SNATIPAddressAllocRef = &api.Reference{ID: *createdObj.Id, Path: *createdObj.Path}
	return actionCreated, nil
}

func (t *snatIPAddressAllocationTask) SetRecoveredReference(state *api.NSXTInfraState, ref *api.Reference, _ *string) {
	state.SNATIPAddressAllocRef = ref
}

func (t *snatIPAddressAllocationTask) ListAll(ctx EnsurerContext, state *api.NSXTInfraState, cursor *string) (interface{}, error) {
	client := ip_pools.NewDefaultIpAllocationsClient(ctx.Connector())
	return client.List(state.SNATIPPoolRef.ID, cursor, nil, nil, nil, nil, nil)
}

func (t *snatIPAddressAllocationTask) EnsureDeleted(ctx EnsurerContext, state *api.NSXTInfraState) (bool, error) {
	client := ip_pools.NewDefaultIpAllocationsClient(ctx.Connector())
	if state.SNATIPAddressAllocRef == nil {
		return false, nil
	}
	err := client.Delete(state.SNATIPPoolRef.ID, state.SNATIPAddressAllocRef.ID)
	if err != nil {
		return false, err
	}
	state.SNATIPAddressAllocRef = nil
	state.SNATIPAddress = nil
	return true, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type snatIPAddressRealizationTask struct{ baseTask }

func NewSNATIPAddressRealizationTask() Task {
	return &snatIPAddressRealizationTask{baseTask{label: "SNAT IP address realization"}}
}

func (t *snatIPAddressRealizationTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return toReference(state.SNATIPAddress)
}

func (t *snatIPAddressRealizationTask) Ensure(ctx EnsurerContext, _ vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	ipAddress, err := getRealizedIPAddress(ctx.Connector(), state.SNATIPAddressAllocRef.Path, 15*time.Second)
	if err != nil {
		return "", err
	}
	state.SNATIPAddress = ipAddress
	return actionFound, nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////

type snatRuleTask struct{ baseTask }

func NewSNATRuleTask() Task {
	return &snatRuleTask{baseTask{label: "SNAT rule"}}
}

func (t *snatRuleTask) Reference(state *api.NSXTInfraState) *api.Reference {
	return state.SNATRuleRef
}

func (t *snatRuleTask) Ensure(ctx EnsurerContext, spec vinfra.NSXTInfraSpec, state *api.NSXTInfraState) (string, error) {
	client := t1nat.NewDefaultNatRulesClient(ctx.Connector())

	rule := model.PolicyNatRule{
		DisplayName:    strptr(spec.FullClusterName()),
		Description:    strptr(description),
		Action:         model.PolicyNatRule_ACTION_SNAT,
		Enabled:        boolptr(true),
		Logging:        boolptr(true),
		SequenceNumber: int64ptr(100),
		Tags:           spec.CreateTags(),

		SourceNetwork:     strptr(spec.WorkersNetwork),
		TranslatedNetwork: strptr(fmt.Sprintf("%s/32", *state.SNATIPAddress)),
	}

	if state.SNATRuleRef != nil {
		oldRule, err := client.Get(state.Tier1GatewayRef.ID, model.PolicyNat_NAT_TYPE_USER, state.SNATRuleRef.ID)
		if isNotFoundError(err) {
			state.SNATRuleRef = nil
			return t.Ensure(ctx, spec, state)
		}
		if err != nil {
			return readingErr(err)
		}
		if *oldRule.DisplayName != *rule.DisplayName ||
			oldRule.Action != rule.Action ||
			oldRule.Enabled == nil ||
			*oldRule.Enabled != *rule.Enabled ||
			oldRule.Logging == nil ||
			*oldRule.Logging != *rule.Logging ||
			oldRule.SequenceNumber == nil ||
			*oldRule.SequenceNumber != *rule.SequenceNumber ||
			oldRule.SourceNetwork == nil ||
			*oldRule.SourceNetwork != *rule.SourceNetwork ||
			oldRule.TranslatedNetwork == nil ||
			*oldRule.TranslatedNetwork != *rule.TranslatedNetwork ||
			oldRule.DestinationNetwork != nil ||
			!equalTags(oldRule.Tags, rule.Tags) {
			err := client.Patch(state.Tier1GatewayRef.ID, model.PolicyNat_NAT_TYPE_USER, state.SNATRuleRef.ID, rule)
			if err != nil {
				return updatingErr(err)
			}
			return actionUpdated, nil
		}
		return actionUnchanged, nil
	}

	id := generateID("snatrule")
	createdObj, err := client.Update(state.Tier1GatewayRef.ID, model.PolicyNat_NAT_TYPE_USER, id, rule)
	if err != nil {
		return creatingErr(err)
	}
	state.SNATRuleRef = &api.Reference{ID: *createdObj.Id, Path: *createdObj.Path}
	return actionCreated, nil
}

func (t *snatRuleTask) SetRecoveredReference(state *api.NSXTInfraState, ref *api.Reference, _ *string) {
	state.SNATRuleRef = ref
}

func (t *snatRuleTask) ListAll(ctx EnsurerContext, state *api.NSXTInfraState, cursor *string) (interface{}, error) {
	client := t1nat.NewDefaultNatRulesClient(ctx.Connector())
	return client.List(state.Tier1GatewayRef.ID, model.PolicyNat_NAT_TYPE_USER, cursor, nil, nil, nil, nil, nil)
}

func (t *snatRuleTask) EnsureDeleted(ctx EnsurerContext, state *api.NSXTInfraState) (bool, error) {
	client := t1nat.NewDefaultNatRulesClient(ctx.Connector())
	if state.SNATRuleRef == nil {
		return false, nil
	}
	err := client.Delete(state.Tier1GatewayRef.ID, model.PolicyNat_NAT_TYPE_USER, state.SNATRuleRef.ID)
	if err != nil {
		return false, err
	}
	state.SNATRuleRef = nil
	return true, nil
}
