// Copyright 2016 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"fmt"
)

//2018/03/20, 1st release is Pangu 0.8.0
//2018/07/28, nuwa test version 1.0.0
//2018/09/20, nuwa version 1.0.3, fixed memory leaking issue under pressure test.
const (
	VersionName  = "nuwa" // Major version name in the Roadmap: Pangu 0.8; Nuwa 1.0; Fuxi 1.1; Shennong 1.2;
	VersionMajor = 1      // Major version component of the current release
	VersionMinor = 0      // Minor version component of the current release
	VersionPatch = 4      // Patch version component of the current release
	VersionMeta  = "rc"   // Version metadata to append to the version string
)

// Version holds the textual version string with Full name.
var VersionWithName = func() string {
	v := fmt.Sprintf("%s %d.%d.%d", VersionName, VersionMajor, VersionMinor, VersionPatch)
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// Version holds the textual version string.
var Version = func() string {
	v := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

// VersionNum only returns the version number.
var VersionNum = func() string {
	v := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	return v
}()

func VersionWithCommit(gitCommit string) string {
	vsn := Version
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}
