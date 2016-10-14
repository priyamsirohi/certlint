package subjectaltname

import (
	"fmt"
	"strings"

	"github.com/globalsign/certlint/certdata"
	"github.com/globalsign/certlint/checks"
)

const checkName = "Subject Alternative Names Check"

func init() {
	checks.RegisterCertificateCheck(checkName, nil, Check)
}

// Check performs a strict verification on the extention according to the standard(s)
func Check(d *certdata.Data) []error {
	var errors []error

	switch d.Type {
	case "PS":
		//if len(d.Cert.EmailAddresses) == 0 {
		//	return fmt.Errorf("Certificate doesn't contain any subjectAltName")
		//}
		for _, s := range d.Cert.EmailAddresses {
			if strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") {
				errors = append(errors, fmt.Errorf("Certificate subjectAltName '%s' starts or ends with a whitespace", s))
			}
		}

	case "DV", "OV", "EV":
		if len(d.Cert.DNSNames) == 0 && len(d.Cert.IPAddresses) == 0 {
			errors = append(errors, fmt.Errorf("Certificate doesn't contain any subjectAltName"))
		}

		var cnInSan bool
		for _, s := range d.Cert.DNSNames {
			if strings.EqualFold(d.Cert.Subject.CommonName, s) {
				cnInSan = true
			}
			if strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") {
				errors = append(errors, fmt.Errorf("Certificate subjectAltName '%s' starts or ends with a whitespace", s))
			}
		}

		// Maybe it's an IP address
		if !cnInSan {
			for _, s := range d.Cert.IPAddresses {
				if strings.EqualFold(d.Cert.Subject.CommonName, s.String()) {
					cnInSan = true
				}
			}
		}

		if !cnInSan {
			errors = append(errors, fmt.Errorf("Certificate CN is not listed in subjectAltName"))
		}
	}

	return errors
}