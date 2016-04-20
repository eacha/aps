package mssql

import (
	"encoding/binary"
	"github.com/eacha/aps/util"
	"regexp"
)

func NewSqlServerVersion(version string) *SqlServerVersion {
	var ssv SqlServerVersion
	re, _ := regexp.Compile(`(\d+)\.(\d+)\.(\d+)\.(\d+)`)

	match := re.FindStringSubmatch(version)

	ssv.VersionNumber = version
	ssv.Major = uint8(util.Atoi(match[1]))
	ssv.Minor = uint8(util.Atoi(match[2]))
	ssv.Build = uint16(util.Atoi(match[3]))
	ssv.SubBuild = uint16(util.Atoi(match[4]))

	return &ssv
}

func (ssv *SqlServerVersion) ToBytes(packet []byte, init int) int {
	packet[init] = ssv.Major
	packet[init+1] = ssv.Minor
	binary.BigEndian.PutUint16(packet[init+2:init+4], ssv.Build)
	binary.BigEndian.PutUint16(packet[init+4:init+6], ssv.SubBuild)

	return init + 6
}

func (plh *PreLoginHeader) ToBytes(packet []byte, init int) int {
	packet[init] = plh.OptionType
	binary.BigEndian.PutUint16(packet[init+1:init+3], plh.Offset)
	binary.BigEndian.PutUint16(packet[init+3:init+5], plh.OptionLength)

	return init + 5
}

func NewPreLoginPacket() *PreLoginPacket {
	var plp PreLoginPacket

	plp.RequestEncryption = false
	plp.InstanceName = ""
	plp.ThreadId = 0
	plp.RequestMars = false

	return &plp
}

func (plp *PreLoginPacket) SetVersion(version string) {
	plp.VersionInfo = NewSqlServerVersion(version)
}

func (plp *PreLoginPacket) ToBytes() []byte {
	var offset uint16 = 21 // (1) Terminator + (5) Version + (5) Encryption + (5) InstOpt + (5) ThreadId
	packet, pos := make([]byte, 1024), 0

	if plp.RequestMars {
		offset += 3 // i think is 5
	}

	plp.SetVersion(DEFAULT_VERSION)

	/* HEADER */

	header := PreLoginHeader{Version, offset, OPTION_LENGTH_CLIENT[Version]}
	pos = header.ToBytes(packet, pos)
	offset += OPTION_LENGTH_CLIENT[Version]

	header = PreLoginHeader{Encryption, offset, OPTION_LENGTH_CLIENT[Encryption]}
	pos = header.ToBytes(packet, pos)
	offset += OPTION_LENGTH_CLIENT[Encryption]

	header = PreLoginHeader{InstOpt, offset, uint16(len(plp.InstanceName) + 1)}
	pos = header.ToBytes(packet, pos)
	offset += uint16(len(plp.InstanceName) + 1)

	header = PreLoginHeader{ThreadId, offset, OPTION_LENGTH_CLIENT[ThreadId]}
	pos = header.ToBytes(packet, pos)
	offset += OPTION_LENGTH_CLIENT[ThreadId]

	if plp.RequestMars {
		header = PreLoginHeader{MARS, offset, OPTION_LENGTH_CLIENT[MARS]}
		pos = header.ToBytes(packet, pos)
		offset += OPTION_LENGTH_CLIENT[MARS]
	}

	packet[pos] = Terminator
	pos += 1

	/* Data */
	pos = plp.VersionInfo.ToBytes(packet, pos)

	packet[pos] = util.BoolToByte(plp.RequestEncryption)
	pos += 1

	pos = util.CopySliceInto(util.ToByteArrayNullTerminated(plp.InstanceName), packet, pos)

	binary.LittleEndian.PutUint32(packet[pos:pos+4], plp.ThreadId)
	pos += 4

	if plp.RequestMars {
		packet[pos] = util.BoolToByte(plp.RequestEncryption)
		pos += 1
	}

	// other header
	realHeader := make([]byte, 8)
	data := packet[:pos]
	packetLength := uint16(len(data) + 8)

	realHeader[0] = PreLogin
	realHeader[1] = 1                                         // message status
	binary.BigEndian.PutUint16(realHeader[2:4], packetLength) // packet length
	binary.BigEndian.PutUint16(realHeader[4:6], 0)            // spid
	realHeader[6] = 1                                         // packet id
	realHeader[7] = 0                                         // window

	return append(realHeader, packet[:pos]...)
}
