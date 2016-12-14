package main

import (
	"strings"

	"gopkg.in/olivere/elastic.v5"
)

func (nodeLog broadLogData) loginAndLogoutSendToEs() {
	l, n := nodeLog.data, nodeLog.nodeName
	nodeLogin, nodeLogout := n+"-login", n+"-logout"
	if strings.Contains(l, "LOGIN") && strings.Contains(l, "FAILED") == false {
		indexReq := elastic.NewBulkIndexRequest().Index(nodeLogin).Type("login").Doc(analysisLogin(l, n))
		wrToEs(indexReq, bulkRequest)
	} else if strings.Contains(l, "LOGOUT") {
		indexReq := elastic.NewBulkIndexRequest().Index(nodeLogout).Type("logout").Doc(analysisOffline(l, n))
		wrToEs(indexReq, bulkRequest)
	}
}

func wrToEs(indexReq *elastic.BulkIndexRequest, bulkRequest *elastic.BulkService) {
	bulkRequest = bulkRequest.Add(indexReq)
	if bulkRequest.NumberOfActions() > put2es {
		_, err := bulkRequest.Do(ctx)
		checkerr(err)

	}
}
