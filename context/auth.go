package contexti

type AuthInfo interface {
	IdStr() string
}

/*

type AuthInterface interface {
	Validate() error
	GenerateToken(secret []byte) (string, error)
	ParseToken(token string, secret []byte) error
}


type Authorization struct {
	AuthInfo `json:"auth"`
	jwt.RegisteredClaims
	AuthInfoRaw string `json:"-"`
}

func (x *Authorization) Validate() error {
	return nil
}

func (x *Authorization) GenerateToken(secret []byte) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, x)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

func (x *Authorization) ParseToken(token string, secret []byte) error {
	_, err := jwti.ParseToken(x, token, secret)
	if err != nil {
		return err
	}
	x.ID = x.AuthInfo.IdStr()
	authBytes, _ := json.Marshal(x.AuthInfo)
	x.AuthInfoRaw = stringsi.BytesToString(authBytes)
	return nil
}
*/
