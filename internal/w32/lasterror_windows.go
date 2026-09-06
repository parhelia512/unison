// Copyright (c) 2021-2026 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package w32

import (
	"errors"
	"syscall"

	"github.com/richardwilkes/toolbox/v2/errs"
)

// lastError describes a failed Windows API call using the last-error value that Call or SyscallN captured for it.
// It must only be used once the call's primary return value has shown that the call failed, since the value is
// meaningless otherwise. A code of 0 is reported as such rather than dropped: an entry point that fails without
// setting an error is itself a useful clue. The optional names map supplies text for codes that FormatMessage does not
// know, such as those defined by OpenGL extensions.
func lastError(fn string, e error, names map[syscall.Errno]string) error {
	var errno syscall.Errno
	if !errors.As(e, &errno) {
		return errs.Newf("%s failed: %v", fn, e)
	}
	if errno == 0 {
		return errs.Newf("%s failed without setting an error code", fn)
	}
	text, ok := names[errno]
	if !ok {
		text = errno.Error()
	}
	return errs.Newf("%s failed: %s (error code 0x%X)", fn, text, uintptr(errno))
}
