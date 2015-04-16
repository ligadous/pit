package api

import (
	"fmt"
	"github.com/alonsovidales/pit/accounts_manager"
	"github.com/alonsovidales/pit/cfg"
	"github.com/alonsovidales/pit/log"
	"github.com/alonsovidales/pit/shards_manager"
	"net/http"
)

const (
	cHealtyPath = "/check_healty"
)

type Api struct {
	shardsManager   *shardsmanager.Manager
	accountsManager *accountsmanager.Manager
	staticPath      string

	muxHttpServer *http.ServeMux
}

func Init(shardsManager *shardsmanager.Manager, accountsManager *accountsmanager.Manager, staticPath string) (api *Api) {
	api = &Api{
		shardsManager:   shardsManager,
		accountsManager: accountsManager,
		muxHttpServer:   http.NewServeMux(),
		staticPath:      staticPath,
	}

	api.registerApis()

	log.Info("Starting API server on port:", int(cfg.GetInt("rec-api", "port")))
	go http.ListenAndServe(fmt.Sprintf(":%d", int(cfg.GetInt("rec-api", "port"))), api.muxHttpServer)

	return
}

func (api *Api) registerApis() {
	api.muxHttpServer.HandleFunc(cHealtyPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	api.muxHttpServer.HandleFunc(shardsmanager.CRecPath, api.shardsManager.ScoresApiHandler)

	api.muxHttpServer.HandleFunc(shardsmanager.CGroupInfoPath, api.shardsManager.GroupInfoApiHandler)

	api.muxHttpServer.HandleFunc(shardsmanager.CRegenerateGroupKey, api.shardsManager.RegenerateGroupKey)
	api.muxHttpServer.HandleFunc(shardsmanager.CGetGroupsByUser, api.shardsManager.GetGroupsByUser)
	api.muxHttpServer.HandleFunc(shardsmanager.CAddUpdateGroup, api.shardsManager.AddUpdateGroup)
	api.muxHttpServer.HandleFunc(shardsmanager.CSetShardsGroup, api.shardsManager.SetShards)

	api.muxHttpServer.HandleFunc(accountsmanager.CRegisterPath, api.accountsManager.Register)
	api.muxHttpServer.HandleFunc(accountsmanager.CVerifyPath, api.accountsManager.Verify)
	api.muxHttpServer.HandleFunc(accountsmanager.CLogsPath, api.accountsManager.Logs)
	api.muxHttpServer.HandleFunc(accountsmanager.CRecoverPassPath, api.accountsManager.RecoverPass)
	api.muxHttpServer.HandleFunc(accountsmanager.CChangePass, api.accountsManager.ChangePass)
	api.muxHttpServer.HandleFunc(accountsmanager.CDisablePath, api.accountsManager.Disable)

	api.muxHttpServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, api.staticPath+r.URL.Path[1:])
	})
}