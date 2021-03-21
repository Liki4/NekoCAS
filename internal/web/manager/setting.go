package manager

import (
	"github.com/NekoWheel/NekoCAS/internal/db"
	"github.com/NekoWheel/NekoCAS/internal/web/context"
	"github.com/NekoWheel/NekoCAS/internal/web/form"
	log "unknwon.dev/clog/v2"
)

func SiteViewHandler(c *context.Context) {
	c.Success("manage/site")
}

func SiteActionHandler(c *context.Context, f form.Site) {
	// 表单报错
	if c.HasError() {
		c.Success("manage/site")
		return
	}

	if f.OpenRegister {
		err := db.SetSetting("open_setting", "on")
		if err != nil {
			log.Error("Failed to set %q to %q", "open_setting", "on")
		}
	} else {
		err := db.SetSetting("open_setting", "off")
		if err != nil {
			log.Error("Failed to set %q to %q", "open_setting", "off")
		}
	}

	_ = db.SetSetting("site_logo", f.SiteLogo)
	_ = db.SetSetting("mail_whitelist", f.MailWhitelist)
	_ = db.SetSetting("privacy", f.Privacy)

	c.Redirect("/manage/site")
}