package mysql

const (
	BYTE  = 1
	SHORT = 2
	INT   = 4
	//STRING_END = 0X00
)

const (
	SEPARATOR = ":"
	PROTOCOL  = "tcp"
)

// Protocol Constants
const (
	//DEFAULT_PORT = 3306
	BUFFER_SIZE = 8192
	HEADER_SIZE = 4
	MSG_TYPE    = HEADER_SIZE
	ERROR_CODE  = MSG_TYPE + BYTE
	ERROR_MSG   = ERROR_CODE + SHORT
	PROTO       = MSG_TYPE
	VERSION     = PROTO + BYTE
	ERROR_TYPE  = 0xFF
)

var capabilities = map[string]uint16{
	"LongPassword":                 0x1,
	"FoundRows":                    0x2,
	"LongColumnFlag":               0x4,
	"ConnectWithDatabase":          0x8,
	"DontAllowDatabaseTableColumn": 0x10,
	"SupportsCompression":          0x20,
	"ODBCClient":                   0x40,
	"SupportsLoadDataLocal":        0x80,
	"IgnoreSpaceBeforeParenthesis": 0x100,
	"Speaks41ProtocolNew":          0x200,
	"InteractiveClient":            0x400,
	"SwitchToSSLAfterHandshake":    0x800,
	"IgnoreSigpipes":               0x1000,
	"SupportsTransactions":         0x2000,
	"Speaks41ProtocolOld":          0x4000,
	"Support41Auth":                0x8000,
}

var charset = map[byte]string{
	0x08: "latin1_swedish_ci",
	0x21: "utf8_general_ci",
	0x63: "binary",
}
