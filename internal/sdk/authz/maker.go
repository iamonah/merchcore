package authz

type TokenMaker interface {
	GenerateToken(JWTData) (string, *Payload, error)
	VerifyToken(string) (*Payload, error)
}

var _ TokenMaker = (*JWTAuthMaker)(nil)
