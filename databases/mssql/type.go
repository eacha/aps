package mssql

type SqlServerVersion struct {
	VersionNumber    string
	Major            uint8
	Minor            uint8
	Build            uint16
	SubBuild         uint16
	ProductName      string
	BrandedVersion   string
	ServicePackLevel string
}

type PreLoginHeader struct {
	OptionType   byte
	Offset       uint16
	OptionLength uint16
}

type PreLoginPacket struct {
	VersionInfo       *SqlServerVersion
	RequestEncryption bool
	InstanceName      string
	ThreadId          uint32
	RequestMars       bool
}

/*
_InferProductVersion = function(self)

		local VERSION_LOOKUP_TABLE = {
			["^6%.0"] = "6.0", ["^6%.5"] = "6.5", ["^7%.0"] = "7.0",
			["^8%.0"] = "2000",	["^9%.0"] = "2005",	["^10%.0"] = "2008",
			["^10%.50"] = "2008 R2", ["^11%.0"] = "2011",
		}

		local product = ""

		for m, v in pairs(VERSION_LOOKUP_TABLE) do
			if ( self.versionNumber:match(m) ) then
				product = v
				self.brandedVersion = product
				break
			end
		end

		self.productName = ("Microsoft SQL Server %s"):format(product)

	end,
*/
