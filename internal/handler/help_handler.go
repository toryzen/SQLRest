// internal/handeler/help_handler.go
package handler

import (
	"sqlrest/internal/util"
	"sqlrest/internal/config" 
	"fmt"
	"net/http"
	//"strconv"
)

func HelpHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	projectID := query.Get("project_id")
	if projectID=="" {
		http.Error(w, "project_id 无效", http.StatusBadRequest)
		return
	}

	projectInfo, apiInfos, err := util.FetchApiInfoByProjectID(projectID)
	if err != nil {
		if query.Get("check") == "check" {
			fmt.Fprintf(w, "false")
			return
		}
		http.Error(w, "未找到对应的Project信息", http.StatusInternalServerError)
		return
	}

	if query.Get("check") == "check" {
		fmt.Fprintf(w, "ok")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `  
	<html>  
	<head>  
		<title>API 使用信息</title>  
		<style>  
			body{font-family:Arial,sans-serif;line-height:1.6;margin:20px;color:#333}  
			h1,h2,h3{color:#444}  
			h2{margin-bottom:0}  
			h3{margin-top:0;font-weight:normal;font-size:18px}  
			pre{background-color:#f4f4f4;border:1px solid #ddd;padding:10px;border-radius:5px}  
			a{color:#1e90ff;text-decoration:none}  
			a:hover{text-decoration:underline}  
			ul{padding-left:20px}  
			li{margin-bottom:5px}  
			hr{border:none;border-top:1px solid #ddd;margin:30px 0}  
			.container{display:flex}  
			.sidebar{width: 250px; padding-right:20px;box-sizing:border-box;position: fixed;top: 20px;   }  
			.main{flex: 1;margin-left: 270px;}  
		</style>  
	</head>  
	<body>  
	<div class="container">  
		<div class="sidebar">  
			<h2>API导航</h2>  
    <ul>`)
	for _, apiInfo := range apiInfos {
		fmt.Fprintf(w, "<li><a href=\"#api-%d\">%s</a></li>", apiInfo.ID, apiInfo.ApiName)
	}
	fmt.Fprintf(w, `</ul>  
    </div>  
    <div class="main">  
        <h2 id="api-%s">项目: %s - API信息</h2>`, projectID, projectInfo.ProjectName)
	for _, apiInfo := range apiInfos {
		if apiInfo.Joint != "" {
			fmt.Fprintf(w, `<h2 id="api-%d">API名称: %s（%s）</h2>`, apiInfo.ID, apiInfo.ApiName, apiInfo.Joint)
		} else {
			fmt.Fprintf(w, `<h2 id="api-%d">API名称: %s</h2>`, apiInfo.ID, apiInfo.ApiName)
		}
		if apiInfo.Memo != "" {
			fmt.Fprintf(w, `<pre>%s</pre>`, apiInfo.Memo)
		}
		fmt.Fprintf(w, `<h3>请求URL:</h3>`)
		fmt.Fprintf(w, `<pre>POST %s/api?api_id=%s</pre>`,config.BaseUrl,apiInfo.ApiId)
		fmt.Fprintf(w, `<h3>Headers:</h3>`)
		if apiInfo.Joint != "" {
			fmt.Fprintf(w, `<pre>Content-Type: application/json<br>authkey: your_auth_key<br>%s_token: your_%s_token</pre>`, apiInfo.Joint, apiInfo.Joint)
		} else {
			fmt.Fprintf(w, `<pre>Content-Type: application/json<br>authkey: your_auth_key</pre>`)
		}
		fmt.Fprintf(w, `<h3>请求Body参数:</h3>`)
		if len(apiInfo.Params) == 0 {
			fmt.Fprintf(w, `<p>无需参数</p>`)
		} else {
			fmt.Fprintf(w, `<pre>{<br>`)
			paramCount := len(apiInfo.Params)
			currentParam := 0
			for paramName, paramDesc := range apiInfo.Params {
				currentParam++
				if currentParam == paramCount {
					fmt.Fprintf(w, `    "%s": "%s"<br>`, paramName, paramDesc)
				} else {
					fmt.Fprintf(w, `    "%s": "%s",<br>`, paramName, paramDesc)
				}
			}
			fmt.Fprintf(w, `}</pre>`)
		}
		fmt.Fprintf(w, `<hr>`)
	}
	fmt.Fprintf(w, `</div>  
</div>  
</body>  
</html>`)
}
