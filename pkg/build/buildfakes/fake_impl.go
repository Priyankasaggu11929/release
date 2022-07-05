/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by counterfeiter. DO NOT EDIT.
package buildfakes

import (
	"sync"

	"github.com/Priyankasaggu11929/release-sdk/git"
)

type FakeImpl struct {
	CheckoutStub        func(*git.Repo, string) error
	checkoutMutex       sync.RWMutex
	checkoutArgsForCall []struct {
		arg1 *git.Repo
		arg2 string
	}
	checkoutReturns struct {
		result1 error
	}
	checkoutReturnsOnCall map[int]struct {
		result1 error
	}
	CommandStub        func(string, ...string) error
	commandMutex       sync.RWMutex
	commandArgsForCall []struct {
		arg1 string
		arg2 []string
	}
	commandReturns struct {
		result1 error
	}
	commandReturnsOnCall map[int]struct {
		result1 error
	}
	OpenRepoStub        func(string) (*git.Repo, error)
	openRepoMutex       sync.RWMutex
	openRepoArgsForCall []struct {
		arg1 string
	}
	openRepoReturns struct {
		result1 *git.Repo
		result2 error
	}
	openRepoReturnsOnCall map[int]struct {
		result1 *git.Repo
		result2 error
	}
	RenameStub        func(string, string) error
	renameMutex       sync.RWMutex
	renameArgsForCall []struct {
		arg1 string
		arg2 string
	}
	renameReturns struct {
		result1 error
	}
	renameReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeImpl) Checkout(arg1 *git.Repo, arg2 string) error {
	fake.checkoutMutex.Lock()
	ret, specificReturn := fake.checkoutReturnsOnCall[len(fake.checkoutArgsForCall)]
	fake.checkoutArgsForCall = append(fake.checkoutArgsForCall, struct {
		arg1 *git.Repo
		arg2 string
	}{arg1, arg2})
	stub := fake.CheckoutStub
	fakeReturns := fake.checkoutReturns
	fake.recordInvocation("Checkout", []interface{}{arg1, arg2})
	fake.checkoutMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeImpl) CheckoutCallCount() int {
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	return len(fake.checkoutArgsForCall)
}

func (fake *FakeImpl) CheckoutCalls(stub func(*git.Repo, string) error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = stub
}

func (fake *FakeImpl) CheckoutArgsForCall(i int) (*git.Repo, string) {
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	argsForCall := fake.checkoutArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeImpl) CheckoutReturns(result1 error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = nil
	fake.checkoutReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeImpl) CheckoutReturnsOnCall(i int, result1 error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = nil
	if fake.checkoutReturnsOnCall == nil {
		fake.checkoutReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkoutReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeImpl) Command(arg1 string, arg2 ...string) error {
	fake.commandMutex.Lock()
	ret, specificReturn := fake.commandReturnsOnCall[len(fake.commandArgsForCall)]
	fake.commandArgsForCall = append(fake.commandArgsForCall, struct {
		arg1 string
		arg2 []string
	}{arg1, arg2})
	stub := fake.CommandStub
	fakeReturns := fake.commandReturns
	fake.recordInvocation("Command", []interface{}{arg1, arg2})
	fake.commandMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2...)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeImpl) CommandCallCount() int {
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	return len(fake.commandArgsForCall)
}

func (fake *FakeImpl) CommandCalls(stub func(string, ...string) error) {
	fake.commandMutex.Lock()
	defer fake.commandMutex.Unlock()
	fake.CommandStub = stub
}

func (fake *FakeImpl) CommandArgsForCall(i int) (string, []string) {
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	argsForCall := fake.commandArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeImpl) CommandReturns(result1 error) {
	fake.commandMutex.Lock()
	defer fake.commandMutex.Unlock()
	fake.CommandStub = nil
	fake.commandReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeImpl) CommandReturnsOnCall(i int, result1 error) {
	fake.commandMutex.Lock()
	defer fake.commandMutex.Unlock()
	fake.CommandStub = nil
	if fake.commandReturnsOnCall == nil {
		fake.commandReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.commandReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeImpl) OpenRepo(arg1 string) (*git.Repo, error) {
	fake.openRepoMutex.Lock()
	ret, specificReturn := fake.openRepoReturnsOnCall[len(fake.openRepoArgsForCall)]
	fake.openRepoArgsForCall = append(fake.openRepoArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.OpenRepoStub
	fakeReturns := fake.openRepoReturns
	fake.recordInvocation("OpenRepo", []interface{}{arg1})
	fake.openRepoMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeImpl) OpenRepoCallCount() int {
	fake.openRepoMutex.RLock()
	defer fake.openRepoMutex.RUnlock()
	return len(fake.openRepoArgsForCall)
}

func (fake *FakeImpl) OpenRepoCalls(stub func(string) (*git.Repo, error)) {
	fake.openRepoMutex.Lock()
	defer fake.openRepoMutex.Unlock()
	fake.OpenRepoStub = stub
}

func (fake *FakeImpl) OpenRepoArgsForCall(i int) string {
	fake.openRepoMutex.RLock()
	defer fake.openRepoMutex.RUnlock()
	argsForCall := fake.openRepoArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeImpl) OpenRepoReturns(result1 *git.Repo, result2 error) {
	fake.openRepoMutex.Lock()
	defer fake.openRepoMutex.Unlock()
	fake.OpenRepoStub = nil
	fake.openRepoReturns = struct {
		result1 *git.Repo
		result2 error
	}{result1, result2}
}

func (fake *FakeImpl) OpenRepoReturnsOnCall(i int, result1 *git.Repo, result2 error) {
	fake.openRepoMutex.Lock()
	defer fake.openRepoMutex.Unlock()
	fake.OpenRepoStub = nil
	if fake.openRepoReturnsOnCall == nil {
		fake.openRepoReturnsOnCall = make(map[int]struct {
			result1 *git.Repo
			result2 error
		})
	}
	fake.openRepoReturnsOnCall[i] = struct {
		result1 *git.Repo
		result2 error
	}{result1, result2}
}

func (fake *FakeImpl) Rename(arg1 string, arg2 string) error {
	fake.renameMutex.Lock()
	ret, specificReturn := fake.renameReturnsOnCall[len(fake.renameArgsForCall)]
	fake.renameArgsForCall = append(fake.renameArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.RenameStub
	fakeReturns := fake.renameReturns
	fake.recordInvocation("Rename", []interface{}{arg1, arg2})
	fake.renameMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeImpl) RenameCallCount() int {
	fake.renameMutex.RLock()
	defer fake.renameMutex.RUnlock()
	return len(fake.renameArgsForCall)
}

func (fake *FakeImpl) RenameCalls(stub func(string, string) error) {
	fake.renameMutex.Lock()
	defer fake.renameMutex.Unlock()
	fake.RenameStub = stub
}

func (fake *FakeImpl) RenameArgsForCall(i int) (string, string) {
	fake.renameMutex.RLock()
	defer fake.renameMutex.RUnlock()
	argsForCall := fake.renameArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeImpl) RenameReturns(result1 error) {
	fake.renameMutex.Lock()
	defer fake.renameMutex.Unlock()
	fake.RenameStub = nil
	fake.renameReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeImpl) RenameReturnsOnCall(i int, result1 error) {
	fake.renameMutex.Lock()
	defer fake.renameMutex.Unlock()
	fake.RenameStub = nil
	if fake.renameReturnsOnCall == nil {
		fake.renameReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.renameReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeImpl) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	fake.commandMutex.RLock()
	defer fake.commandMutex.RUnlock()
	fake.openRepoMutex.RLock()
	defer fake.openRepoMutex.RUnlock()
	fake.renameMutex.RLock()
	defer fake.renameMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeImpl) recordInvocation(key string, args []interface{}) {
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
