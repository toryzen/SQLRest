// internal/util/check_auth.go
package util

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	//"strconv"
	"strings"
)

func CheckAuthKey(authKey string, apiID string, requestID string) error {
	Logger.Infof("checkAuthKey called with authKey: %s, apiID: %s, requestID: %s", authKey, apiID, requestID)
	var authKeyRecord struct {
		ID        int
		AuthKey   string
		ProjectID string
		ApiIDs    string
	}
	err := MainDB.QueryRow("SELECT a.id, a.authkey, a.project_id, a.api_ids FROM auth_key a JOIN project p ON a.project_id = p.project_id WHERE a.authkey = ? AND a.is_del = 0", authKey).Scan(&authKeyRecord.ID, &authKeyRecord.AuthKey, &authKeyRecord.ProjectID, &authKeyRecord.ApiIDs)

	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("authkey 鉴权失败")
		}
		return err
	}
	found := false
	if authKeyRecord.ApiIDs == "-1" {
		row := MainDB.QueryRow("SELECT 1 FROM data_api WHERE api_id = ? AND project_id = ? AND is_del = 0", apiID, authKeyRecord.ProjectID)
		if err := row.Scan(new(int)); err == nil {
			found = true
		}
	} else {
		apiIDs := strings.Split(authKeyRecord.ApiIDs, ",")
		for _, id := range apiIDs {
			if id == apiID {
				// 检查API是否属于给定的project_id
				row := MainDB.QueryRow("SELECT 1 FROM data_api WHERE api_id = ? AND project_id = ? AND is_del = 0", apiID, authKeyRecord.ProjectID)
				if err := row.Scan(new(int)); err == nil {
					found = true
					break
				}
			}
		}
	}

	if !found {
		return errors.New("authkey 鉴权失败，该authkey没有此api的操作权限！")
	}

	return nil
}

type Response struct {
	Userid      string `json:"userid"`
	Employeeid  string `json:"employeeid"`
	Alldeptname string `json:"allDeptname"`
	Deptname    string `json:"deptname"`
	Account     string `json:"account"`
	Username    string `json:"username"`
	Status      bool   `json:"status"`
}

// 自定义二次鉴权及变量获取(返回给params)
func JointAuth(jointname string, token string, requestID string) (map[string]interface{}, error) {
	Logger.Infof("JointAuth called with jointname: %s, token: %s, requestID: %s", jointname, token, requestID)
	var params map[string]interface{}
	params = map[string]interface{}{
		"userid":      "",
		"employeeid":  "",
		"deptname":    "",
		"username":    "",
	}
	return params, nil
}
