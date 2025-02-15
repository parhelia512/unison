// Copyright (c) 2021-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package unison

import (
	"runtime"

	"github.com/richardwilkes/unison/enums/blur"
	"github.com/richardwilkes/unison/internal/skia"
)

// MaskFilter performs a transformation on the mask before drawing it.
type MaskFilter struct {
	filter skia.MaskFilter
}

func newMaskFilter(filter skia.MaskFilter) *MaskFilter {
	if filter == nil {
		return nil
	}
	f := &MaskFilter{filter: filter}
	runtime.AddCleanup(f, func(sf skia.MaskFilter) {
		ReleaseOnUIThread(func() {
			skia.MaskFilterUnref(sf)
		})
	}, f.filter)
	return f
}

func (f *MaskFilter) filterOrNil() skia.MaskFilter {
	if f == nil {
		return nil
	}
	return f.filter
}

// NewBlurMaskFilter returns a new blur mask filter. sigma is the standard deviation of the gaussian blur to apply. Must
// be greater than 0. If respectMatrix is true, the blur's sigma is modified by the current matrix.
func NewBlurMaskFilter(style blur.Enum, sigma float32, respectMatrix bool) *MaskFilter {
	if sigma <= 0 {
		sigma = 1
	}
	return newMaskFilter(skia.MaskFilterNewBlurWithFlags(skia.Blur(style), sigma, respectMatrix))
}

// NewTableMaskFilter returns a new table mask filter. The table should be 256 elements long. If shorter, it will be
// expanded to 256 elements and the new entries will be filled with 0.
func NewTableMaskFilter(table []byte) *MaskFilter {
	if len(table) < 256 {
		t := make([]byte, 256)
		copy(t, table)
		table = t
	}
	return newMaskFilter(skia.MaskFilterNewTable(table))
}

// NewGammaMaskFilter returns a new gamma mask filter.
func NewGammaMaskFilter(gamma float32) *MaskFilter {
	return newMaskFilter(skia.MaskFilterNewGamma(gamma))
}

// NewClipMaskFilter returns a new clip mask filter.
func NewClipMaskFilter(minimum, maximum byte) *MaskFilter {
	return newMaskFilter(skia.MaskFilterNewClip(minimum, maximum))
}

// NewShaderMaskFilter returns a new shader mask filter.
func NewShaderMaskFilter(shader *Shader) *MaskFilter {
	return newMaskFilter(skia.MaskFilterNewShader(shader.shader))
}
