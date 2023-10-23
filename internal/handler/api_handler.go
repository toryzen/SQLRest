// internal/handeler/api_handler.go
package handler

import (
	"sqlrest/internal/util"
	"sqlrest/pkg/response"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	//"strconv"
	"strings"
)

func PrepareResponse(reqid string, code int, errMsg string, data interface{}) response.Response {
	if data == nil {
		data = []map[string]interface{}{}
	}
	return response.Response{
		Code:  code,
		Data:  data,
		ReqId: reqid,
		Msg:   errMsg,
	}
}

func checkSQLInjection(params map[string]interface{}) error {
	sqlKeywords := []string{"select", "insert", "update", "drop", "truncate"}
	for _, v := range params {
		value := fmt.Sprintf("%v", v)
		lowerValue := strings.ToLower(value)
		for _, keyword := range sqlKeywords {
			pattern := fmt.Sprintf(`;\s*%s`, keyword)
			matched, _ := regexp.MatchString(pattern, lowerValue)
			if matched {
				return errors.New("parameter values should not contain select, insert, update, drop, truncate keywords")
			}
		}
	}
	return nil
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	requestID := uuid.New().String()
	util.Logger.Debugf("ApiHandler called, requestID: %s", requestID)

	uparams := make(map[string]interface{})
	w.Header().Set("Content-Type", "application/json")

	if r.ContentLength > 0 {
		util.Logger.Debugf("Processing request body, requestID: %s", requestID)

		err := json.NewDecoder(r.Body).Decode(&uparams)
		if err != nil {
			util.Logger.Warningf("Invalid request body content, requestID: %s", requestID)
			resp := PrepareResponse(requestID, 500, "Invalid JSON request in the body", nil)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if err := checkSQLInjection(uparams); err != nil {
			resp := PrepareResponse(requestID, 400, err.Error(), nil)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	var apiIDList []string
	query := r.URL.Query()
	apiIDStr := query.Get("api_id")

	if apiIDStr == "" {
		util.Logger.Warningf("api_id is empty, requestID: %s", requestID)
		resp := PrepareResponse(requestID, 400, "api_id cannot be empty", nil)
		json.NewEncoder(w).Encode(resp)
		return
	}

	matched, err := regexp.MatchString(`^([a-zA-Z0-9]+)(,[a-zA-Z0-9]+)*$`, apiIDStr)  
	if err != nil || !matched {
		util.Logger.Warningf("Invalid api_id parameter, requestID: %s", requestID)
		resp := PrepareResponse(requestID, 500, "Invalid api_id parameter", nil)
		json.NewEncoder(w).Encode(resp)
		return
	}
	splitAPIIDs := strings.Split(apiIDStr, ",")
	for _, id := range splitAPIIDs {
		if id!=""{
			apiIDList = append(apiIDList, id)
		}
		
	}

	authKey := r.Header.Get("authkey")
	if authKey == "" {
		util.Logger.Warningf("authkey is empty, requestID: %s", requestID)
		resp := PrepareResponse(requestID, 401, "authkey cannot be empty", nil)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var subIDList []string
	subQueryApiIds := query.Get("sub_api_id")
	if subQueryApiIds != "" {
		ids, err := CheckSubQuery(subQueryApiIds, authKey, requestID)
		if err != nil {
			util.Logger.Warningf("Subquery authentication failed, requestID: %s", requestID)
			resp := PrepareResponse(requestID, 401, err.Error(), nil)
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(resp)
			return
		}
		subIDList = ids
		util.Logger.Debugf("Get SubIDList by query sub_query_api_ids: %v, requestID: %s", subIDList, requestID)
	}

	token := query.Get("token")

	var apiResponse []interface{}

	for i, apiID := range apiIDList {
		util.Logger.Infof("Processing apiIDList[%d] = %s, requestID: %s", i, apiID, requestID)
		err := util.CheckAuthKey(authKey, apiID, requestID)
		if err != nil {
			resp := PrepareResponse(requestID, 401, err.Error(), nil)
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(resp)
			return
		}
		resp, err := util.FetchData(apiID, uparams, token, subIDList, requestID)
		if err != nil {
			resp = PrepareResponse(requestID, 500, err.Error(), nil)
			json.NewEncoder(w).Encode(resp)
			return
		}
		apiResponse = append(apiResponse, resp)
	}
	if len(apiResponse) == 1 {
		resp := PrepareResponse(requestID, 200, "", apiResponse[0])
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := PrepareResponse(requestID, 200, "", apiResponse)
	json.NewEncoder(w).Encode(resp)
}

func CheckSubQuery(subQueryApiIds string, authKey string ,requestID string ) ([]string, error) {
	util.Logger.Debug("CheckSubQuery called")

	subQueryApiIdsList := []string{}
	matched, err := regexp.MatchString(`^([a-zA-Z0-9]+)(,[a-zA-Z0-9]+)*$`, subQueryApiIds)
	if err != nil || !matched {
		return subQueryApiIdsList, errors.New("invalid subquery_api_ids")
	}
	splitAPIIDs := strings.Split(subQueryApiIds, ",")
	for _, id := range splitAPIIDs {
		if id!=""  {
			subQueryApiIdsList = append(subQueryApiIdsList, id)
		}
	}
	for _, apiID := range subQueryApiIdsList {
		util.Logger.Infof("Checking authKey for subquery API ID %s", apiID)
		err := util.CheckAuthKey(authKey, apiID ,requestID)
		if err != nil {
			return subQueryApiIdsList, errors.New(err.Error())
		}
	}
	return subQueryApiIdsList, nil
}
