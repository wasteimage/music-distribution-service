package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const loginPageName = "login"

var _ Page = &loginPage{}

type loginPage struct {
	wrongPass string
	warning   string
	page
}

func (p *loginPage) Get(rq RequestContext) {
	trashCookie := http.Cookie{
		Name:    sessionCookie,
		Path:    "/",
		Expires: time.Now(),
	}
	http.SetCookie(rq.rw, &trashCookie)
	var currentTheme string
	var navLogo string
	var colorTheme string
	if rq.theme == "SGreen" {
		currentTheme = "style_black.css"
		navLogo = "logo_white.png"
		colorTheme = "success"
	} else {
		currentTheme = "style.css"
		navLogo = "logo.png"
		colorTheme = "primary"
	}
	locales, err := p.loc.TranslatePage(rq.r.Header.Get("Accept-Language"),
		"login_p", "login_user", "login_pass", "login_complete", "login_reg",
		"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
		"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	)
	var params = map[string]interface{}{
		"loggedIn": rq.userID > 0,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"wrong":    p.wrongPass,
		"warning":  p.warning,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
	}
	p.wrongPass = "display: none;"
	p.warning = "display: none;"
	if rq.userID > 0 {
		http.Redirect(rq.rw, rq.r, "../cabinet", http.StatusFound)
	}
	err = p.tmpl.Lookup("login").Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *loginPage) Post(rq RequestContext) {
	username := rq.r.FormValue("text")
	password := rq.r.FormValue("password")
	if len(username) < 1 || len(password) < 1 {
		p.warning = "display: block;"
		http.Redirect(rq.rw, rq.r, "/login/", http.StatusFound)
	}
	userId, err := p.db.GetUserId(username, password)
	if err != nil || userId <= 0 && len(username) > 0 && len(password) > 0 {
		p.wrongPass = "display: block;"
		http.Redirect(rq.rw, rq.r, "/login/", http.StatusFound)
	}
	session := http.Cookie{
		Name:    sessionCookie,
		Value:   strconv.Itoa(userId),
		Path:    "/",
		Domain:  "",
		Expires: time.Now().Add(time.Hour * 48),
	}
	rq.r.AddCookie(&session)
	http.SetCookie(rq.rw, &session)
	http.Redirect(rq.rw, rq.r, "../cabinet/", http.StatusFound)
	return
}
