package server

import (
	"github.com/0xPolygonHermez/zkevm-bridge-service/log"
	"github.com/0xPolygonHermez/zkevm-bridge-service/nacos"
)

func RegisterNacos(cfg nacos.Config) {
	var err error
	if cfg.NacosUrls != "" {
		err = nacos.InitNacosClient(cfg.NacosUrls, cfg.NamespaceId, cfg.ApplicationName, cfg.ExternalListenAddr)
	}
	log.Debugf("Init nacos NacosUrls[%s] NamespaceId[%s] ApplicationName[%s] ExternalListenAddr[%s] Error[%v]", cfg.NacosUrls, cfg.NamespaceId, cfg.ApplicationName, cfg.ExternalListenAddr, err)
}
