// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type FakeDashboardUrlGenerator struct {
	DashboardUrlStub        func(instanceID string, plan serviceadapter.Plan, manifest bosh.BoshManifest) (serviceadapter.DashboardUrl, error)
	dashboardUrlMutex       sync.RWMutex
	dashboardUrlArgsForCall []struct {
		instanceID string
		plan       serviceadapter.Plan
		manifest   bosh.BoshManifest
	}
	dashboardUrlReturns struct {
		result1 serviceadapter.DashboardUrl
		result2 error
	}
	dashboardUrlReturnsOnCall map[int]struct {
		result1 serviceadapter.DashboardUrl
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDashboardUrlGenerator) DashboardUrl(instanceID string, plan serviceadapter.Plan, manifest bosh.BoshManifest) (serviceadapter.DashboardUrl, error) {
	fake.dashboardUrlMutex.Lock()
	ret, specificReturn := fake.dashboardUrlReturnsOnCall[len(fake.dashboardUrlArgsForCall)]
	fake.dashboardUrlArgsForCall = append(fake.dashboardUrlArgsForCall, struct {
		instanceID string
		plan       serviceadapter.Plan
		manifest   bosh.BoshManifest
	}{instanceID, plan, manifest})
	fake.recordInvocation("DashboardUrl", []interface{}{instanceID, plan, manifest})
	fake.dashboardUrlMutex.Unlock()
	if fake.DashboardUrlStub != nil {
		return fake.DashboardUrlStub(instanceID, plan, manifest)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.dashboardUrlReturns.result1, fake.dashboardUrlReturns.result2
}

func (fake *FakeDashboardUrlGenerator) DashboardUrlCallCount() int {
	fake.dashboardUrlMutex.RLock()
	defer fake.dashboardUrlMutex.RUnlock()
	return len(fake.dashboardUrlArgsForCall)
}

func (fake *FakeDashboardUrlGenerator) DashboardUrlArgsForCall(i int) (string, serviceadapter.Plan, bosh.BoshManifest) {
	fake.dashboardUrlMutex.RLock()
	defer fake.dashboardUrlMutex.RUnlock()
	return fake.dashboardUrlArgsForCall[i].instanceID, fake.dashboardUrlArgsForCall[i].plan, fake.dashboardUrlArgsForCall[i].manifest
}

func (fake *FakeDashboardUrlGenerator) DashboardUrlReturns(result1 serviceadapter.DashboardUrl, result2 error) {
	fake.DashboardUrlStub = nil
	fake.dashboardUrlReturns = struct {
		result1 serviceadapter.DashboardUrl
		result2 error
	}{result1, result2}
}

func (fake *FakeDashboardUrlGenerator) DashboardUrlReturnsOnCall(i int, result1 serviceadapter.DashboardUrl, result2 error) {
	fake.DashboardUrlStub = nil
	if fake.dashboardUrlReturnsOnCall == nil {
		fake.dashboardUrlReturnsOnCall = make(map[int]struct {
			result1 serviceadapter.DashboardUrl
			result2 error
		})
	}
	fake.dashboardUrlReturnsOnCall[i] = struct {
		result1 serviceadapter.DashboardUrl
		result2 error
	}{result1, result2}
}

func (fake *FakeDashboardUrlGenerator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.dashboardUrlMutex.RLock()
	defer fake.dashboardUrlMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDashboardUrlGenerator) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ serviceadapter.DashboardUrlGenerator = new(FakeDashboardUrlGenerator)
