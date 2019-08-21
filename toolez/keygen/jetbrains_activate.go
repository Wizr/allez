package keygen

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

// base 36 integer
var _n = "t76eqou0a5rdnij8mixvut88ca8lk9y513x5xxj6dcihqar4buvofcyu7ymbt84rllt9awioc7b1ra3ridjxt4ofkzm1i2iriaz"
var _d = "ekwnx0sfrsk6fy9vsgvgp9yqz1r4vzcnnmb6efvn1klz22e614qfwh41r8sch3fvwmp9yc2t2zadz4ip7eodiagcffxovccch5t"
var pkey *rsa.PrivateKey

// prepare rsa constants
func init() {
	n, _ := new(big.Int).SetString(_n, 36)
	d, _ := new(big.Int).SetString(_d, 36)
	p, _ := new(big.Int).SetString("526cqefo0m5h5qydu25d301g8ntrhmbvq7om4ldki3kgmcnhrd", 36)
	q, _ := new(big.Int).SetString("5rq25sorip2lewym841g6h2qx68oq38n8pdxvpb88tpiq4b1mb", 36)
	pkey = &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: n,
			E: 65537,
		},
		D:      d,
		Primes: []*big.Int{p, q},
	}
}

// ActivateIdea process Get request for iris
func ActivateIdea(c *gin.Context) {
	salt := c.Query("salt")
	userName := c.Query("userName")
	log.Printf(`UserName: %v
HostName: %v
MachineId: %v
ProductCode: %v
ProductFamilyId: %v
BuildNumber: %v
ClientVersion: %v
VersionNumber: %v
secure: %v
salt: %v
`,
		userName,
		c.Query("hostName"),
		c.Query("machineId"),
		c.Query("productCode"),
		c.Query("productFamilyId"),
		c.Query("buildNumber"),
		c.Query("clientVersion"),
		c.Query("versionNumber"),
		c.Query("secure"),
		c.Query("salt"))

	ticket := fmt.Sprintf(`<ObtainTicketResponse>
	<message></message>
	<prolongationPeriod>86400000</prolongationPeriod>
	<responseCode>OK</responseCode>
	<salt>%s</salt>
	<ticketId>1</ticketId>
	<ticketProperties>licensee=%s	licenseType=0	</ticketProperties>
</ObtainTicketResponse>`, salt, userName)
	hashed := md5.Sum([]byte(ticket))
	sig, _ := rsa.SignPKCS1v15(rand.Reader, pkey, crypto.MD5, hashed[:])

	resp := fmt.Sprintf("<!-- %x -->\n%s", sig, ticket)
	c.String(http.StatusOK, resp)
}
