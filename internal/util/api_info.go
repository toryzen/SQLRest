// internal/util/api_info.go
package util

import (
	"regexp"
)

type ProjectApiInfo struct {
	ProjectId int `json:"project_id"`
	ProjectName string `json:"project_name"`
}

type ApiInfo struct {
	ID          int                    `json:"id"`
	ApiId       string                 `json:"api_id"`
	ApiName     string                 `json:"api_name"`
	ProjectName string                 `json:"project_name"`
	Joint       string                 `json:"joint"`
	SourceSQL   string                 `json:"-"`
	Params      map[string]interface{} `json:"params"`
	ProjectType string                 `json:"project_type,omitempty"`
	Memo        string                 `json:"memo"`
}

func FetchApiInfoByProjectID(projectID string) (ProjectApiInfo, []ApiInfo, error) {
	var project ProjectApiInfo
	err := MainDB.QueryRow("SELECT id,project_name FROM project WHERE project_id = ? AND is_del = 0", projectID).Scan(&project.ProjectId, &project.ProjectName)
	if err != nil {
		return ProjectApiInfo{}, nil, err
	}

	rows, err := MainDB.Query("SELECT id, api_id , apiname,joint, sourcesql, memo FROM data_api WHERE project_id = ? AND is_del = 0", projectID)
	if err != nil {
		return ProjectApiInfo{}, nil, err
	}
	defer rows.Close()

	var apiInfos []ApiInfo
	for rows.Next() {
		var apiInfo ApiInfo
		err = rows.Scan(&apiInfo.ID, &apiInfo.ApiId, &apiInfo.ApiName, &apiInfo.Joint, &apiInfo.SourceSQL, &apiInfo.Memo)
		if err != nil {
			Logger.Warning(err)
			return ProjectApiInfo{}, nil, err
		}

		params, err := extractCustomParams(apiInfo.SourceSQL)
		if err != nil {
			return ProjectApiInfo{}, nil, err
		}
		apiInfo.Params = params
		apiInfos = append(apiInfos, apiInfo)
	}

	return project, apiInfos, nil
}

func extractCustomParams(sourceSQL string) (map[string]interface{}, error) {
	re, err := regexp.Compile("{(.*?)}")
	if err != nil {
		return nil, err
	}

	matches := re.FindAllStringSubmatch(sourceSQL, -1)

	params := make(map[string]interface{})
	for _, match := range matches {
		param := match[1]

		if param == "userid" || param == "employeeid" || param == "alldeptname" || param == "deptname" || param == "account" || param == "username" {
			continue
		}

		params[param] = ""
	}

	return params, nil
}
