package session

type Sid struct {
	HexString string
	SidString string
}

func NewSid(hexString string) *Sid {
	return &Sid{
		HexString: hexString,
		SidString: "sid:" + hexString,
	}
}
