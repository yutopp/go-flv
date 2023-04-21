//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package flv

import (
	"github.com/yutopp/go-flv/tag"
)

type Flags uint8

const (
	FlagsAudio Flags = 0x01
	FlagsVideo Flags = 0x02
)

var HeaderSignature = []byte{0x46, 0x4c, 0x56} // F, L, V

const HeaderLength uint32 = 9

type Header struct {
	Version uint8
	Flags
	DataOffset uint32
}

type Body struct {
	Tags []*tag.FlvTag
}
