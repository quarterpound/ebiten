// Copyright 2019 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jsutil

import (
	"syscall/js"
)

// isTypedArrayWritable represents whether TypedArray is writable or not.
// TypedArray's properties are not writable in the Web standard, but are writable with go2cpp.
// This enables to avoid unnecessary allocations of js.Value.
var isTypedArrayWritable = js.Global().Get("go2cpp").Truthy()

// temporaryBuffer is a temporary buffer used at gl.readPixels or gl.texSubImage2D.
// The read data is converted to Go's byte slice as soon as possible.
// To avoid often allocating ArrayBuffer, reuse the buffer whenever possible.
var temporaryBuffer = js.Global().Get("ArrayBuffer").New(16)

func ensureTemporaryBufferSize(byteLength int) {
	if bufl := temporaryBuffer.Get("byteLength").Int(); bufl < byteLength {
		for bufl < byteLength {
			bufl *= 2
		}
		temporaryBuffer = js.Global().Get("ArrayBuffer").New(bufl)
	}
}

func TemporaryUint8Array(byteLength int) js.Value {
	ensureTemporaryBufferSize(byteLength)
	return uint8Array(temporaryBuffer, 0, byteLength)
}

var uint8ArrayObj js.Value

func uint8Array(buffer js.Value, byteOffset, byteLength int) js.Value {
	if isTypedArrayWritable {
		if Equal(uint8ArrayObj, js.Undefined()) {
			uint8ArrayObj = js.Global().Get("Uint8Array").New()
		}
		uint8ArrayObj.Set("buffer", buffer)
		uint8ArrayObj.Set("byteOffset", byteOffset)
		uint8ArrayObj.Set("byteLength", byteLength)
		return uint8ArrayObj
	}
	return js.Global().Get("Uint8Array").New(buffer, byteOffset, byteLength)
}

func TemporaryFloat32Array(byteLength int) js.Value {
	ensureTemporaryBufferSize(byteLength)
	return float32Array(temporaryBuffer, 0, byteLength)
}

var float32ArrayObj js.Value

func float32Array(buffer js.Value, byteOffset, byteLength int) js.Value {
	if isTypedArrayWritable {
		if Equal(float32ArrayObj, js.Undefined()) {
			float32ArrayObj = js.Global().Get("Float32Array").New()
		}
		float32ArrayObj.Set("buffer", buffer)
		float32ArrayObj.Set("byteOffset", byteOffset)
		float32ArrayObj.Set("byteLength", byteLength)
		return float32ArrayObj
	}
	return js.Global().Get("Float32Array").New(buffer, byteOffset/4, byteLength/4)
}
