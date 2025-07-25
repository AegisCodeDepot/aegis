package admin

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/bctnry/aegis/routes"
	"github.com/bctnry/aegis/templates"
)

func bindAdminIndexConfigController(ctx *routes.RouterContext) {
	http.HandleFunc("GET /admin/index-config", routes.WithLog(func(w http.ResponseWriter, r *http.Request) {
		loginInfo, err := routes.GenerateLoginInfoModel(ctx, r)
		if err != nil { routes.FoundAt(w, "/") }
		if !loginInfo.LoggedIn { routes.FoundAt(w, "/") }
		if !loginInfo.IsAdmin { routes.FoundAt(w, "/") }
		iType := ctx.Config.FrontPageType
		var namespace, repository, fileContent string
		if strings.HasPrefix(iType, "static/") {
			fileContent = ctx.Config.FrontPageContent
		} else if strings.HasPrefix(iType, "namespace/") {
			namespace = iType[len("namespace/"):]
			iType = "namespace"
		} else if strings.HasPrefix(iType, "repository/") {
			k := strings.Split(iType, ":")
			if len(k) <= 1 {
				namespace = ""
				repository = k[0]
			} else {
				namespace = k[0]
				repository = k[1]
			}
			iType = "repository"
		}
		routes.LogTemplateError(ctx.LoadTemplate("admin/index-config").Execute(w, &templates.AdminIndexConfigTemplateModel{
			Config: ctx.Config,
			LoginInfo: loginInfo,
			ErrorMsg: "",
			IndexType: iType,
			IndexNamespace: namespace,
			IndexRepository: repository,
			IndexFileContent: fileContent,
		}))
	}))
	http.HandleFunc("POST /admin/index-config", routes.WithLog(func(w http.ResponseWriter, r *http.Request) {
		loginInfo, err := routes.GenerateLoginInfoModel(ctx, r)
		if err != nil { routes.FoundAt(w, "/") }
		if !loginInfo.LoggedIn { routes.FoundAt(w, "/") }
		if !loginInfo.IsAdmin { routes.FoundAt(w, "/") }
		err = r.ParseForm()
		if err != nil {
			routes.LogTemplateError(ctx.LoadTemplate("admin/index-config").Execute(w, templates.AdminIndexConfigTemplateModel{
				Config: ctx.Config,
				LoginInfo: loginInfo,
				ErrorMsg: fmt.Sprintf("Failed to parse request: %s\n", err.Error()),
			}))
			return
		}
		indexType := r.Form.Get("index-type")
		namespace := r.Form.Get("namespace")
		repository := r.Form.Get("repository")
		fileName := r.Form.Get("file-name")
		if !strings.HasPrefix(fileName, "/") {
			fileName = "/" + fileName
		}
		fileContent := r.Form.Get("file-content")
		switch indexType {
		case "all/namespace":
			ctx.Config.FrontPageType = "all/namespace"
		case "all/repository":
			ctx.Config.FrontPageType = "all/repository"
		case "namespace":
			ctx.Config.FrontPageType = fmt.Sprintf("namespace/%s", namespace)
		case "repository":
			ctx.Config.FrontPageType = fmt.Sprintf("repository/%s:%s", namespace, repository)
		case "static":
			ctx.Config.FrontPageType = fileName
			p := path.Join(ctx.Config.StaticAssetDirectory, fileName)
			f, err := os.OpenFile(p, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0664)
			if err != nil {
				routes.LogTemplateError(ctx.LoadTemplate("admin/index-config").Execute(w, templates.AdminIndexConfigTemplateModel{
					Config: ctx.Config,
					LoginInfo: loginInfo,
					ErrorMsg: fmt.Sprintf("Failed to open static file: %s\n", err.Error()),
					IndexType: indexType,
					IndexNamespace: namespace,
					IndexRepository: repository,
					IndexFileContent: fileContent,
				}))
				return
			}
			defer f.Close()
			_, err = f.Write([]byte(fileContent))
			if err != nil {
				routes.LogTemplateError(ctx.LoadTemplate("admin/index-config").Execute(w, templates.AdminIndexConfigTemplateModel{
					Config: ctx.Config,
					LoginInfo: loginInfo,
					ErrorMsg: fmt.Sprintf("Failed to save file: %s\n", err.Error()),
					IndexType: indexType,
					IndexNamespace: namespace,
					IndexRepository: repository,
					IndexFileContent: fileContent,
				}))
				return
			}
		}
		err = ctx.Config.Sync()
		if err != nil {
			routes.LogTemplateError(ctx.LoadTemplate("admin/index-config").Execute(w, templates.AdminIndexConfigTemplateModel{
				Config: ctx.Config,
				LoginInfo: loginInfo,
				ErrorMsg: fmt.Sprintf("Failed to save config: %s\n", err.Error()),
				IndexType: indexType,
				IndexNamespace: namespace,
				IndexRepository: repository,
				IndexFileContent: fileContent,
			}))
			return
		}
		
		routes.LogTemplateError(ctx.LoadTemplate("admin/index-config").Execute(w, templates.AdminIndexConfigTemplateModel{
			Config: ctx.Config,
			LoginInfo: loginInfo,
			ErrorMsg: "Updated.",
			IndexType: indexType,
			IndexNamespace: namespace,
			IndexRepository: repository,
			IndexFileContent: fileContent,
		}))
		return
	}))
}

