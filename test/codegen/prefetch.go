// asmcheck

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"runtime"
	"sort"
	"unsafe"
)

// This test check instruction generation for Prefetch API
func CheckPrefetchGeneration() {
	arr := []int{4, 5, 6, 7, 0, 1, 2, 3}
	p := uintptr(unsafe.Pointer(&arr))
	// amd64:`PREFETCHNTA\t\(.*\)`
	// arm64:`PRFM\t\(.*\), \$4`
	runtime.PrefetchMemory(p, runtime.PrefetchLocalityNTA)
	// amd64:`PREFETCHT2\t\(.*\)`
	// arm64:`PRFM\t\(.*\), \$4`
	runtime.PrefetchMemory(p, runtime.PrefetchLocality2)
	// amd64:`PREFETCHT1\t\(.*\)`
	// arm64:`PRFM\t\(.*\), \$2`
	runtime.PrefetchMemory(p, runtime.PrefetchLocality1)
	// amd64:`PREFETCHT0\t\(.*\)`
	// arm64:`PRFM\t\(.*\), \$0`
	runtime.PrefetchMemory(p, runtime.PrefetchLocality0)
	sort.Ints(arr)
}
