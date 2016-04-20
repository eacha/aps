package mysql

/*
Source Script:
https://svn.nmap.org/nmap/scripts/mysql-info.nse
https://svn.nmap.org/nmap/nselib/mysql.lua
https://dev.mysql.com/doc/internals/en/connection-phase-packets.html#packet-Protocol::Handshake
*/

import (
	"encoding/binary"
	"net"
	"strconv"
	"strings"
	"time"
)

func NewMySQLConnection(ip string, port int, connectionTimeout, ioTimeout time.Duration, data *MySQL) (*MySQLConnection, error) {
	var c MySQLConnection

	c.ip = ip
	c.port = port
	c.connectionTimeout = connectionTimeout * time.Second
	c.ioTimeout = ioTimeout * time.Second
	c.data = data

	conn, err := net.DialTimeout(PROTOCOL, c.ip+SEPARATOR+strconv.Itoa(c.port), c.connectionTimeout)
	if err != nil {
		return nil, err
	}
	c.conn = conn

	if err = c.conn.SetDeadline(time.Now().Add(c.ioTimeout)); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *MySQLConnection) Read(b []byte) (int, error) {
	n, err := c.conn.Read(b)

	if err != nil {
		return 0, err
	}

	if err = c.conn.SetDeadline(time.Now().Add(c.ioTimeout)); err != nil {
		return 0, err
	}

	return n, nil
}

func (c *MySQLConnection) Write(b []byte) (int, error) {
	n, err := c.conn.Write(b)

	if err != nil {
		return 0, err
	}

	if err = c.conn.SetDeadline(time.Now().Add(c.ioTimeout)); err != nil {
		return 0, err
	}

	return n, nil
}

func (c *MySQLConnection) GetBanner() error {
	//var pos int
	buffer := make([]byte, BUFFER_SIZE)

	readBytes, err := c.Read(buffer)
	if err != nil {
		return err
	}

	if err = c.data.decodeBanner(buffer[:readBytes]); err != nil {
		return err
	}
	//data.Banner = buffer[:readBytes]
	//
	//if errorField := buffer[MSG_TYPE]; errorField == ERROR_TYPE {
	//	data.MySQLError.Code = binary.LittleEndian.Uint16(buffer[ERROR_CODE : ERROR_CODE+SHORT])
	//	data.MySQLError.Message = string(buffer[ERROR_MSG:readBytes])
	//	return nil
	//}
	//
	//data.Proto = buffer[PROTO]
	//data.Version, pos = NullString(buffer, VERSION)
	//data.ThreadId = binary.LittleEndian.Uint32(buffer[pos : pos+INT])
	//pos += INT + 9
	//data.Capabilities = parseCapabilities(binary.LittleEndian.Uint16(buffer[pos : pos+SHORT]))
	//pos += SHORT
	//data.Charset = buffer[pos]
	//pos += BYTE
	//data.Status = binary.LittleEndian.Uint16(buffer[pos : pos+SHORT])

	return nil
}

func decodeHeader(header []byte) (len, id uint32) {
	tmp := binary.LittleEndian.Uint32(header)
	len = tmp & 0xffffff
	id = tmp >> 24

	return len, id
}

func (m *MySQL) decodeBanner(banner []byte) error {
	var pos int

	if length, _ := decodeHeader(banner[0:HEADER_SIZE]); int(length+HEADER_SIZE) != len(banner) {
		// TODO change to correct error
		return nil
	}

	m.Banner = banner

	if errorField := banner[MSG_TYPE]; errorField == ERROR_TYPE {
		m.MySQLError.Code = binary.LittleEndian.Uint16(banner[ERROR_CODE : ERROR_CODE+SHORT])
		m.MySQLError.Message = string(banner[ERROR_MSG:len(banner)])
		return nil
	}

	m.Proto = banner[PROTO]
	m.Version, pos = NullString(banner, VERSION)
	m.ThreadId = binary.LittleEndian.Uint32(banner[pos : pos+INT])
	pos += INT + 9
	m.Capabilities = parseCapabilities(binary.LittleEndian.Uint16(banner[pos : pos+SHORT]))
	pos += SHORT
	m.Charset = banner[pos]
	pos += BYTE
	m.Status = binary.LittleEndian.Uint16(banner[pos : pos+SHORT])

	return nil
}

func parseCapabilities(cap uint16) []string {
	capList, cnt := make([]string, 16), 0

	for k, v := range capabilities {
		if cap&v > 0 {
			capList[cnt] = k
			cnt += 1
		}
	}

	return capList[:cnt]
}

/* mandarlo a util */
func NullString(buffer []byte, init int) (str string, pos int) {
	nullPos := strings.Index(string(buffer[init:]), "\x00") + init
	s := buffer[init:nullPos]

	return string(s), nullPos + 1
}
