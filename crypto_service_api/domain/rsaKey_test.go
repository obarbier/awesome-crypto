package domain

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

//  openssl genrsa 2048
const validRsa2048 = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAvwjOscHHNbEDSshyznRB+7FUjYcnShYMgjZHC8Mk6OLa804q\nTWFYuEicrjSeQgJv0u9iXpd4KIubQvIS4v7swNQt/5VIY7ad4wC4ElEaMN393gQ5\nlVpguu/oH4LVxKCBASILCt+nJwKRH847+Mc2/SRRfZc6nD1rf4QUc2z6dxrSfpWk\n9DQ8tohyokJtSLJBxFzrhOJph03VUMsS19SNWyorhyRZ+S4J3id9EVqTPbRRT7Om\nKILpkKReHugKqqgcaAMZ8uYSQ3fXB8m7dZau2zyLxIABL9mOmm7V3FMOVDp0cEEW\ndm3QspUAGPn85U5M5bngPwWmqKBTCpDXBdQUWQIDAQABAoIBAGADojdPaLuAo2Hz\ny1gesIc7Qn77nfCrkk/jzeOIikWOt3MwJyzLL5c3z5/zDFOq+98tt+hJAviF6lxE\nZN+4NvBX9GKP+lk7kLFmTB3Qq3RCLvVmdDavvQxXxfgIfz6bWRtjq98kN9Pwg1ZU\nBmpsXiPvHBIebzPV3vCxRsIY1O4FKjgaYbawrSVM06dssfO2kF+zfaMiH5mnluXP\nYjO5kbqkS3NJPJAFrmZHNQcJFXp4p2Z9rLPC6FtGHJmm8RUQyz5ebeHOpaftWM9+\nmChqIARsZcgnC6dSS/haMecrFzBmUGnJypuv5wVVFIFbzwy414yjllDCPtAx0+gq\n3tiMNr0CgYEA+Y7WgxNsS1XGB/J896BMjZOxBsrqPdGzKRjn527Vpjchyl3Ygd5S\nfV5B0vmZSOi1QqVXFx4OOvwOrBeSX7AzMe2oTKAYQqGKAk4zXRfW5Q+GJTb2ZFHc\nY9NAMsNhmx7imP9mLUIq1Y1b6TYQfFIFrWtPlJ2aNuxKTvIseCu1OjcCgYEAw/c5\n8C8GofTojdRR7z0n8IeIH9XBhz3EY2Jj04I1rQ80KH3Xhc5JWvYodj6JQijN1K+H\ngUf3SdzGUVL/gJ8HTkqvJCZWKI38d7Ivco7BwE9Jyzvu+dcMiakgdYcYTpVKq27J\nQ/wHFb83KcjJxiVHzAnsTgmDQMYxkJa2D7EJne8CgYEAj3hT0KVg8+qK39TDjWUF\nvbrz4hDUG9gr5Ouhnwa0I0u8zGepafgTemmu9Ah03FqUoo0FhY/M5JI2KS+gAgz9\nUa3svKipad0Ox4aHtvRWofeLymdPvZrmVimD1etHePOHmCf0aP6KO516ApgHYEGT\nbACujqUQnJS5n6tQb4HJPX8CgYAAzGINC3QAdun3ofTPf7VI8pRoZMuMDIFfUkhL\n1Uz4roYs4A5fui5sU3JowOp4PYhRJIHt0eg9AcxBCpCF6p/x/rXl9M4HDkUIC87L\nra82ZFxNmqnnlKu1Z938/Jbpwwvx1Nq3DzDMMuI7pljEGOTI/QVccAd73RLYnvSQ\n7cy3OwKBgDXk4rQ1a0jrPyYA8smwNKFsxhjogzXjKcAQdyJjjx4Ef9h4AK1KVuvE\n0lC3DtX8NlQosro0SXr3CY2SL0FoeTEZWJ4z+HQH59ggvIDm4Fe1GNOeZ/deMO+9\n//tjD/tx63YIB3IiJ1jNJeSdB+NXP4Vis0fiOUlCuLQ3NL8GGBEt\n-----END RSA PRIVATE KEY-----\n"

//  openssl genrsa 4096
const validRsa4096 = "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"

// 	openssl genrsa 2048 | openssl pkcs8 -topk8 -nocrypt
const pkcs8Rsa = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC5hQRZdKYie6D/\nn6+M+EszzRVjfjPiSMznd1Zhj9GQFPYxohhWqhgLixTrxM9DXP3/TXDRX23TRSu4\n6MRFy6qhcBHVf5VHUuSM2x4iaWEXRWkWrJYmyPXD/9TYCkAl8Z6Nq87kb3WpdzJ8\nXJOslTmtuNLSHXTr2yVV3ngeygJh8m8v8MV5f19a07gTuvFBSnurWPw+d0UrRIcI\nYI52zbN0Ajn06vWPZhAxT8a5NalLsgPebJXhNgIFzilFZXUnXpLCqWPNKZ4R4CoK\nWJPWHqSBwfArxpww+NTxf0yhfQcrU36Ff+dmBzbKgIrBKUlYbYtlg5HBZ7HoRNBG\nsGI4OPA7AgMBAAECggEBAKSWC9vJDzL8d0MRSk9IYH9ebKFN733LlG2tg+ceDo9C\n6X/zDKCmWpqzEyZv/mkG8Rg0fehiPy716OotJyO8om6C3G+KtscGFVmZc8yXrNlW\nbPr+tl9GXjM6nnvj7DE9gKqzR+OFtt9XrmSCRUkRQpCKrg5Wr8onK6JYsjyufxqg\n1B3QOjy+P2kbZ7hQPv6hAC/RYjWOGSvjcnNoTTdKeAqITDxk/fU8ayBsxVQPs4r3\nhD+iEMgb2BsgdKMk684zUFi87Dj/ctX1zSswnB0qhwMAcDXAwdt8PX3gzGpxCTi9\nDEFYIb6OVrrLWOKWZHzc4fFmTyhyI4vm/plF0LSaPVECgYEA7oDJaPdAfbiMhc7e\nMWqnZr3QzTaXU01d32wO0G4vKzDKwAfMj20p7Cg5lTVUfsU09zD3/EgjUHRk/buu\nJ6+eGc4w6N5aQVVxz/cx9iuzTLLvklc+vIU4T0PVie5+MoVSFbidh0dHddT/+qSt\n1+YTJOdLxFmMwkeBYjKiEcsRu28CgYEAxyEsNFcYDhAqOIyGoSVEHSztro41mQpc\nImprya1R3HPcFX2PAASeXJEEWD2KkcOmhUMhT45ctRPMl1WT1iimtCGYjmvCmmF/\nSihvnLKxs8540A1lVFLCBhHZF3Xm0W9ssFa5KtKw5rAVfesVZOkCIA1pYvtFbWd+\n1C7cESu84fUCgYAwc9V3D5P4dn+Fx4r4OxSbGMDMj+SaNcN2WjuAOII4ogbukCcM\nlD7KDTn1iAoMXv/tn/MhO36BH8RMj85HnbPexjbFeDaZw0QF3dA2lJYuZMOq1TKX\nlfDkmYFOLjdRCCiu5PyLuP1ZgNYoE0CF9eW5v4ty7kZcSa6NRoAKYVjO3QKBgQCO\nXyhkz7RyMZqOTeLgCm310i6p9CFcJ20SajZgvpvd27SKZPg+Eg9LrZ+Gm5GcgF9p\nvkJtyCJ+kQZhWR1XLD9sYOzbPy6nBHhnBBww1A57uW7lif5d2MHCZzZpMLH0Ig96\n0LaZaIR0m4byPYdRW8taMVydGXxdKXcjq9FKMZRdlQKBgFv8vvTUiFBirGp7FiYn\nF6WfDcAm/Ml4WqQ4Ote/lWKAWMiqgiy7/+znTzZpcsdnqQfGaafdK+Pr+FC/e6L4\n7BuV1wV4wwVQ26JkE5CMk6+bMFwBmcCZOd5WZf5q1p8ODDXlJNOEfXGkgFEUBZSq\nHG0sehEYN48zG0ZVeEcl4fNP\n-----END PRIVATE KEY-----"

func TestRsaKey(t *testing.T) {
	tests := []struct {
		name      string
		pemString string
		wants     []byte
		wantsErr  bool
	}{
		{
			name:      "Valid Test for RSA 2048",
			pemString: validRsa2048,
			wants:     []byte(validRsa2048),
			wantsErr:  false,
		},
		{
			name:      "Valid Test for RSA 4096",
			pemString: validRsa4096,
			wants:     []byte(validRsa4096),
			wantsErr:  false,
		},
		{
			name:      "Should Error pkcs8 format",
			pemString: pkcs8Rsa,
			wants:     []byte(pkcs8Rsa),
			wantsErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testKey := &RSA{}
			err := testKey.Marshal(test.pemString)
			if (err != nil) != test.wantsErr {
				t.Errorf("%s: failed to create rsa key from string: %s", test.name, err)
				return
			}

			if test.wantsErr {
				return // Hacked way to stop at marshalling
			}

			assert.Equal(t, testKey.Type(), "rsa", "validate key type is rsa")

			rsaKeyBytes, err := testKey.UnMarshal()
			if (err != nil) != test.wantsErr {
				t.Errorf("%s: failed to create rsa key from string: %s", test.name, err)
				return
			}

			if !bytes.Equal(rsaKeyBytes, test.wants) {
				t.Errorf("Failled")
			}
		})
	}
}
