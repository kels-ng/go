// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

// Locality hint to point a location in the cache hierarchy
const (
	// T0 (temporal data)-prefetch data into all levels of the cache hierarchy.
	PrefetchLocality0 = iota
	// T1 (temporal data with respect to first level cache)-prefetch data into
	// level 2 cache and higher.
	PrefetchLocality1
	// T2 (temporal data with respect to second level cache)-prefetch data into
	// level 2 cache and higher.
	PrefetchLocality2
	// NTA (non-temporal data with respect to all cache levels)-prefetch data
	// into non-temporal cache structure and into a location close to the
	// processor, minimizing cache pollution.
	PrefetchLocalityNTA
)

// PrefetchMemory - fetch data from memory addr to CPU cache line
// using locality hint - level
func PrefetchMemory(addr uintptr, level int) {}
