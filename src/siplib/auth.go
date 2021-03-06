package siplib

import "fmt"
import "errors"
import "strings"
import "crypto/md5"

// Struct used to keep track of authorisation info for a REGISTER request.

type SIPAuth struct {
	User string
	Password string
	Nonce string
	NonceCount int
	Realm string
	Algo string
	Opaque string
	Qop string
	Proxy string
	Type string
}

func (a *SIPAuth) SetCreds(user string, password string) {
	a.User = user
	a.Password = password
}

func (a SIPAuth) Generate(uri string, method string) (string,error) {
	output := fmt.Sprintf("Digest username=\"%s\",realm=\"%s\",nonce=\"%s\",uri=\"%s\"", a.User, a.Realm, a.Nonce, uri)
	auth_digest := ""
	method_digest := ""
	cnonce := GenerateHexUUID(16)
	if (strings.ToLower(a.Qop) == "auth") {
		nc := fmt.Sprintf("%08d", a.NonceCount)
		output += fmt.Sprintf(",cnonce=\"%s\",nc=%s", cnonce, nc)
	}
	if (a.Algo == "") || (strings.ToLower(a.Algo) == "md5") {
		authstr := fmt.Sprintf("%s:%s:%s", a.User, a.Realm, a.Password)
		auth_digest = fmt.Sprintf("%x", md5.Sum([]byte(authstr)))
		output += ",algorithm=md5"
	}
	if (strings.ToLower(a.Algo) == "md5-sess") {
		cnonce := GenerateHexUUID(16)
		nc := fmt.Sprintf("%08d", a.NonceCount)
		output += fmt.Sprintf(",cnonce=\"%s\",nc=%s", cnonce, nc)
		authstr := fmt.Sprintf("%s:%s:%s", a.User, a.Realm, a.Password)
		auth_digest = fmt.Sprintf("%x", md5.Sum([]byte(authstr)))
		sess := fmt.Sprintf("%s:%s:%s", auth_digest, a.Nonce, cnonce)
		auth_digest = fmt.Sprintf("%x", md5.Sum([]byte(sess)))
		output += ",algorithm=md5-sess"
	}
	if auth_digest == "" {
		return "",errors.New("Unknown algorithm: " + a.Algo)
	}
	if (a.Qop == "") || (strings.ToLower(a.Qop) == "auth") {
		method_str := fmt.Sprintf("%s:%s", method, uri)
		method_digest = fmt.Sprintf("%x", md5.Sum([]byte(method_str)))
		output += ",qop=auth"
	}
	if strings.ToLower(a.Qop) == "auth-int" {
		return "",errors.New("Authentication type 'auth-int' is not supported.")
	}
	if strings.ToLower(a.Qop) == "auth" {
		nc := fmt.Sprintf("%08d", a.NonceCount)
		res := fmt.Sprintf("%s:%s:%s:%s:%s:%s", auth_digest, a.Nonce, nc, cnonce, a.Qop, method_digest)
		res_digest := fmt.Sprintf("%x", md5.Sum([]byte(res)))
		output += fmt.Sprintf(",response=\"%s\"", res_digest)
	} else {
		res := fmt.Sprintf("%s:%s:%s", auth_digest, a.Nonce, method_digest)
		res_digest := fmt.Sprintf("%x", md5.Sum([]byte(res)))
		output += fmt.Sprintf(",response=\"%s\"", res_digest)
	}
	if a.Opaque != "" {
		output += fmt.Sprintf(",opaque=\"%s\"", a.Opaque)
	}
	return output,nil
}

// Extract authentication info from SIP response; returns a SIPAuth object with everything but the username and password set.

func GetAuthInfo(resp SIPResponse) (SIPAuth,error) {
	output := SIPAuth{}
	output.NonceCount = 1
	if val,ok := resp.Headers["WWW-Authenticate"]; ok {
		parts := strings.Split(val, " ")
		if len(parts) < 2 {
			return output,errors.New("Failed to parse the WWW-Authenticate header.")
		}
		output.Type = parts[0]
		body := strings.Join(parts[1:], " ")
		parameters := strings.Split(body, ",")
		for _,part := range parameters {
			element := strings.Split(part, "=")
			if len(element) != 2 {
				return output,errors.New("Failed to parse the WWW-Authenticate header.")
			}
			key := strings.TrimSpace(element[0])
			value := strings.TrimSpace(element[1])
			switch strings.ToLower(key) {
				case "realm":
					output.Realm = strings.Replace(value, "\"", "", -1)
				case "nonce":
					output.Nonce = strings.Replace(value, "\"", "", -1)
				case "algorithm":
					output.Algo = value
				case "opaque":
					output.Opaque = strings.Replace(value, "\"", "", -1)
				case "qop":
					output.Qop = strings.Replace(value, "\"", "", -1)
			}
		}
		return output,nil
	} else {
		return output,errors.New("There was no WWW-Authenticate header in the response.")
	}
}

