//go:build ignore

package templates

import "github.com/bctnry/gitus/pkg/gitus"

type DepotNoNamespaceModel struct {
	Config *gitus.GitusConfig
	RepositoryList []struct{RelPath string; Description string}
	DepotName string
	LoginInfo *LoginInfoModel
}

