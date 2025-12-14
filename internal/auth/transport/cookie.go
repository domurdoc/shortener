package transport

import "net/http"

// CookieTransport is a Transport implementation that reads and writes authentication tokens
// using HTTP cookies. It allows configuration of various cookie attributes to control
// security and scope, such as name, lifetime, path, and SameSite policy.
type CookieTransport struct {
	name     string        // name is the cookie name used to store the token.
	maxAge   int           // maxAge specifies the lifetime of the cookie in seconds.
	secure   bool          // secure indicates that the cookie should only be sent over HTTPS.
	httpOnly bool          // httpOnly prevents client-side JavaScript access to the cookie.
	path     string        // path limits the scope of the cookie to a specific URL path.
	sameSite http.SameSite // sameSite controls whether the cookie is sent with cross-site requests.
}

func NewCookie(name string, maxAge int, secure bool) *CookieTransport {
	return &CookieTransport{
		name:     name,
		maxAge:   maxAge,
		secure:   secure,
		httpOnly: true,
		path:     "/",
		sameSite: http.SameSiteStrictMode,
	}
}

func (c *CookieTransport) Read(r *http.Request) (string, error) {
	cookie, noCookieErr := r.Cookie(c.name)
	if noCookieErr != nil {
		return "", noCookieErr
	}
	if invalidCookieErr := cookie.Valid(); invalidCookieErr != nil {
		return "", invalidCookieErr
	}
	return cookie.Value, nil
}

func (c *CookieTransport) Write(w http.ResponseWriter, tokenString string) error {
	cookie := http.Cookie{
		Name:     c.name,
		Value:    tokenString,
		SameSite: c.sameSite,
		MaxAge:   c.maxAge,
		Secure:   c.secure,
		HttpOnly: c.httpOnly,
		Path:     c.path,
	}
	if err := cookie.Valid(); err != nil {
		return err
	}
	http.SetCookie(w, &cookie)
	return nil
}
