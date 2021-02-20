package netutils

//
//Copyright 2019 Telenor Digital AS
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
import (
	"fmt"
	"net"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindIPAddress(t *testing.T) {
	addr, err := FindPublicIPv4()
	t.Logf("Found address: %v", addr.String())
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindLoopback(t *testing.T) {
	FindLoopbackIPv4Interface()
}

func TestIsLoopback(t *testing.T) {
	if IsLoopbackAddress("0.0.0.0:1234") {
		t.Fatal("0.0.0.0 shouldn't be loopback")
	}

	if IsLoopbackAddress(":4711") {
		t.Fatal(":4711 should not be loopback")
	}

	if !IsLoopbackAddress("127.0.0.1:4711") {
		t.Fatal("127.0.0.1 is loopback")
	}

	if IsLoopbackAddress("example.com:4711") {
		t.Fatal("example.com is not loopback")
	}
}

func TestInterfaceSorter(t *testing.T) {
	assert := require.New(t)
	// Just grab the local interfaces and run through sort.Sort
	ifaces, err := net.Interfaces()
	assert.NoError(err)
	ifs := &interfaceSorter{ifaces}
	sort.Sort(ifs)

	oldIndex := 0
	for n, i := range ifs.interfaces {
		fmt.Println(i.Index)
		if i.Index <= oldIndex {
			assert.Failf("Sorting is broken", "Interface at index %d is %d but should be bigger", n, i.Index)
		}
	}
}
