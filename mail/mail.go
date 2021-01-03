package mail

import (
	"crypto/tls"
	"html/template"
	"path/filepath"
	"sync"
	"time"

	"github.com/NekoWheel/NekoCAS/conf"
	"gopkg.in/gomail.v2"
	"gopkg.in/macaron.v1"
)

var (
	tplRender     *macaron.TplRender
	tplRenderOnce sync.Once
)

// render 根据给定的信息渲染邮件模板
func render(tpl string, data map[string]interface{}) (string, error) {
	tplRenderOnce.Do(func() {
		opt := &macaron.RenderOptions{
			Directory:         filepath.Join("templates", "mail"),
			AppendDirectories: []string{filepath.Join("templates", "mail")},
			Extensions:        []string{".tmpl", ".html"},
			Funcs: []template.FuncMap{map[string]interface{}{
				"Year": func() int {
					return time.Now().Year()
				},
			}},
		}

		ts := macaron.NewTemplateSet()
		ts.Set(macaron.DEFAULT_TPL_SET_NAME, opt)
		tplRender = &macaron.TplRender{
			TemplateSet: ts,
			Opt:         opt,
		}
	})

	return tplRender.HTMLString(tpl, data)
}

func SendActivationMail(to, code string) error {
	data := map[string]interface{}{
		"Email": to,
		"Link":  conf.Get().Site.BaseURL + "/activate_code?code=" + code,
	}
	body, err := render("activate", data)
	if err != nil {
		return err
	}

	return send(to, "激活您的 Neko 账号", body)
}

func SendLostPasswordMail(to, code string) error {
	data := map[string]interface{}{
		"Email": to,
		"Link":  conf.Get().Site.BaseURL + "/reset_password?code=" + code,
	}
	body, err := render("reset_password", data)
	if err != nil {
		return err
	}

	return send(to, "您正在找回您的 Neko 账号密码", body)
}

func send(to, title, content string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Get().Mail.Account)
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)

	d := gomail.NewDialer(
		conf.Get().Mail.SMTP,
		conf.Get().Mail.Port,
		conf.Get().Mail.Account,
		conf.Get().Mail.Password,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}