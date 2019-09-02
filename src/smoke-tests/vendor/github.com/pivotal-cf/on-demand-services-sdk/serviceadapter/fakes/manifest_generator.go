// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/on-demand-services-sdk/serviceadapter"
)

type FakeManifestGenerator struct {
	GenerateManifestStub        func(params serviceadapter.GenerateManifestParams) (serviceadapter.GenerateManifestOutput, error)
	generateManifestMutex       sync.RWMutex
	generateManifestArgsForCall []struct {
		params serviceadapter.GenerateManifestParams
	}
	generateManifestReturns struct {
		result1 serviceadapter.GenerateManifestOutput
		result2 error
	}
	generateManifestReturnsOnCall map[int]struct {
		result1 serviceadapter.GenerateManifestOutput
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManifestGenerator) GenerateManifest(params serviceadapter.GenerateManifestParams) (serviceadapter.GenerateManifestOutput, error) {
	fake.generateManifestMutex.Lock()
	ret, specificReturn := fake.generateManifestReturnsOnCall[len(fake.generateManifestArgsForCall)]
	fake.generateManifestArgsForCall = append(fake.generateManifestArgsForCall, struct {
		params serviceadapter.GenerateManifestParams
	}{params})
	fake.recordInvocation("GenerateManifest", []interface{}{params})
	fake.generateManifestMutex.Unlock()
	if fake.GenerateManifestStub != nil {
		return fake.GenerateManifestStub(params)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.generateManifestReturns.result1, fake.generateManifestReturns.result2
}

func (fake *FakeManifestGenerator) GenerateManifestCallCount() int {
	fake.generateManifestMutex.RLock()
	defer fake.generateManifestMutex.RUnlock()
	return len(fake.generateManifestArgsForCall)
}

func (fake *FakeManifestGenerator) GenerateManifestArgsForCall(i int) serviceadapter.GenerateManifestParams {
	fake.generateManifestMutex.RLock()
	defer fake.generateManifestMutex.RUnlock()
	return fake.generateManifestArgsForCall[i].params
}

func (fake *FakeManifestGenerator) GenerateManifestReturns(result1 serviceadapter.GenerateManifestOutput, result2 error) {
	fake.GenerateManifestStub = nil
	fake.generateManifestReturns = struct {
		result1 serviceadapter.GenerateManifestOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeManifestGenerator) GenerateManifestReturnsOnCall(i int, result1 serviceadapter.GenerateManifestOutput, result2 error) {
	fake.GenerateManifestStub = nil
	if fake.generateManifestReturnsOnCall == nil {
		fake.generateManifestReturnsOnCall = make(map[int]struct {
			result1 serviceadapter.GenerateManifestOutput
			result2 error
		})
	}
	fake.generateManifestReturnsOnCall[i] = struct {
		result1 serviceadapter.GenerateManifestOutput
		result2 error
	}{result1, result2}
}

func (fake *FakeManifestGenerator) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.generateManifestMutex.RLock()
	defer fake.generateManifestMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeManifestGenerator) recordInvocation(key string, args []interface{}) {
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

var _ serviceadapter.ManifestGenerator = new(FakeManifestGenerator)