// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime_test

import (
	"runtime"
	"testing"
	"unsafe"
)

// Sizes of data fields
const (
	metaSize    = 64
	dataSize    = 64 * 1024
	packetCount = 1000
)

// Packet for emulation
type TestPacket struct {
	meta []byte
	data []byte
}

var messages = make(chan *TestPacket)
var processed = make(chan *TestPacket)
var next = make(chan *TestPacket)
var done = make(chan bool)

// Emulate packet processor with prefetch
func packetProcPref() {
	for i := 0; i < packetCount; i++ {
		// Get packet
		data := <-messages

		// Prefetch meta info from packet
		runtime.PrefetchMemory(uintptr(unsafe.Pointer(&data.meta)), runtime.PrefetchLocality0)

		// Create new packet and noise cache during zeroing
		arr := &TestPacket{
			meta: make([]byte, metaSize),
			data: make([]byte, dataSize),
		}

		// Create new meta with changed old one
		for i := 0; i < metaSize; i++ {
			arr.meta[i] = data.meta[i] + 1
		}
		// Return result
		processed <- data
		next <- arr
	}
}

// Emulate packet processor
func packetProc() {
	for i := 0; i < packetCount; i++ {
		// Get packet
		data := <-messages

		// Create new packet and noise cache during zeroing
		arr := &TestPacket{
			meta: make([]byte, metaSize),
			data: make([]byte, dataSize),
		}

		// Create new meta with changed old one
		for i := 0; i < metaSize; i++ {
			arr.meta[i] = data.meta[i] + 1
		}
		// Return result
		processed <- arr
		next <- arr
	}
}

// Packet generator
func packetGen() {
	// Emulate recieving new packet and send it to processor
	for i := 0; i < packetCount; i++ {
		arr := &TestPacket{
			meta: make([]byte, metaSize),
			data: make([]byte, dataSize),
		}

		messages <- arr
	}
}

// Translate packets
func packetResponse() {
	tmp := make([]TestPacket, packetCount)
	for i := 0; i < packetCount; i++ {
		tmp[i] = *<-next
	}
	done <- true
}

// Delete packets
func packetUtilization() {
	tmp := make([]TestPacket, packetCount)
	for i := 0; i < packetCount; i++ {
		tmp[i] = *<-processed
	}
	done <- true
}

// Bench performance using prefetch
func BenchmarkWithPrefetch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go packetGen()
		go packetProcPref()
		go packetResponse()
		go packetUtilization()
		<-done
		<-done
	}
}

// Bench performance without prefetch
func BenchmarkWithoutPrefetch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go packetGen()
		go packetProc()
		go packetResponse()
		go packetUtilization()
		<-done
		<-done
	}
}
