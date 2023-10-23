// internal/util/fetch_data.go
package util

import (
	"sqlrest/internal/config"
	"sqlrest/internal/model"
	"sqlrest/pkg/response"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var (
	MainDB            *sql.DB
	TargetDBs         = make(map[int]*sql.DB)
	ErrNoAvailableApi = errors.New("未找到该ID的API，请检查api_id")
	ErrNoAvailableDb  = errors.New("未找到该API的DBSource，请联系该API管理员")
)

func ConnectDB() {
	var err error
	MainDB, err = sql.Open("mysql", config.DBConnectionString)
	if err != nil {
		Logger.Fatal(err)
	}
	MainDB.SetMaxOpenConns(10)
	MainDB.SetMaxIdleConns(5)
}

func FetchApiInfo(apiID string, requestID string) (model.DataAPI, error) {
	Logger.Infof("FetchApiInfo called with apiID: %s, requestID: %s", apiID, requestID)
	var dataAPI model.DataAPI

	err := MainDB.QueryRow("SELECT id, apiname, db_id, joint, sourcesql FROM data_api WHERE api_id = ? AND is_del = 0", apiID).Scan(&dataAPI.ID, &dataAPI.ApiName, &dataAPI.DbID, &dataAPI.Joint, &dataAPI.SourceSQL)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.DataAPI{}, ErrNoAvailableApi
		}
		return model.DataAPI{}, err
	}
	return dataAPI, nil
}

func FetchData(apiID string, uparams map[string]interface{}, token string, subIDList []string, requestID string) (interface{}, error) {
	Logger.Infof("FetchData called with apiID: %s, requestID: %s", apiID, requestID)

	//获取API信息
	dataAPI, err := FetchApiInfo(apiID, requestID)
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})

	//判断是否需要joint
	if dataAPI.Joint != "" {
		if token == "" {
			return nil, errors.New("联合鉴权已开启，token不能为空")
		}
		params, err = JointAuth(dataAPI.Joint, token, requestID)
		if err != nil {
			return nil, err
		}
	}

	//不允许用户变量与系统变量（joint变量）冲突
	merged, err := mergeMapsCheckConflict(params, uparams)
	if err != nil {
		return nil, err
	} else {
		params = merged
	}

	var dbSource model.DBSource
	err = MainDB.QueryRow("SELECT id, dbname, ip, port, user, pwd FROM db_source WHERE db_id = ? AND is_del = 0", dataAPI.DbID).Scan(&dbSource.ID, &dbSource.DbName, &dbSource.IP, &dbSource.Port, &dbSource.User, &dbSource.Pwd)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoAvailableDb
		}
		return nil, err
	}

	Logger.Debugf("Fetched dbSource, requestID: %s", requestID)
	targetDB, ok := TargetDBs[dbSource.ID]
	if !ok {
		targetDB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true", dbSource.User, dbSource.Pwd, dbSource.IP, dbSource.Port, dbSource.DbName))
		if err != nil {
			return nil, err
		}
		targetDB.SetMaxOpenConns(10)
		targetDB.SetMaxIdleConns(5)
		TargetDBs[dbSource.ID] = targetDB
	}

	// Replace placeholders in SourceSQL with values from params
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		dataAPI.SourceSQL = strings.Replace(dataAPI.SourceSQL, placeholder, fmt.Sprintf("%v", value), -1)
	}

	Logger.Debugf("Executing query on targetDB: %s, requestID: %s", dataAPI.SourceSQL, requestID)
	// 检查是否为INSERT操作和UPDATE操作
	sqlTrimmed := strings.ToUpper(strings.TrimSpace(dataAPI.SourceSQL))
	isInsert := strings.HasPrefix(sqlTrimmed, "INSERT")
	isUpdate := strings.HasPrefix(sqlTrimmed, "UPDATE")

	//data := make(map[string]interface{})
	if isInsert {
		Logger.Debugf("isInsert %s, requestID: %s", dataAPI.SourceSQL, requestID)
		// 执行INSERT操作
		result, err := targetDB.Exec(dataAPI.SourceSQL)
		if err != nil {
			return nil, err
		}

		// 获取插入的主键ID
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		// 获取插入的行数
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, err
		}

		// 将主键ID作为数据返回
		data := map[string]interface{}{
			"last_insert_id": lastInsertID,
			"rows_affected":  rowsAffected,
		}
		Logger.Debugf("Query result data: %v, requestID: %s", data, requestID)
		return data, nil
	} else if isUpdate {
		Logger.Debugf("isUpdate %s, requestID: %s", dataAPI.SourceSQL, requestID)
		// 执行UPDATE操作
		result, err := targetDB.Exec(dataAPI.SourceSQL)
		if err != nil {
			return nil, err
		}

		// 获取受影响的行数
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, err
		}

		// 返回成功与否的结果
		//success := rowsAffected > 0
		data := map[string]interface{}{
			"rows_affected": rowsAffected,
		}
		Logger.Debugf("Query result data: %v, requestID: %s", data, requestID)
		return data, nil
	} else {
		Logger.Debugf("isSelect %s, requestID: %s", dataAPI.SourceSQL, requestID)
		var data []map[string]interface{}
		// 执行非INSERT和非UPDATE操作
		rows, err := targetDB.Query(dataAPI.SourceSQL)
		if err != nil {
			if err == sql.ErrNoRows {
				data := map[string]interface{}{}
				return data, nil
			}
			return nil, err
		}
		defer rows.Close()

		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))

		for rows.Next() {
			for i := range values {
				pointers[i] = &values[i]
			}

			err = rows.Scan(pointers...)
			if err != nil {
				return response.Response{}, err
			}

			rowData := make(map[string]interface{})
			for i, colName := range columns {
				var v interface{}
				val := values[i]
				b, ok := val.([]byte)
				if ok {
					v = string(b)
				} else {
					v = val
				}
				rowData[colName] = v
			}
			data = append(data, rowData)
		}
		for k, subID := range subIDList {
			Logger.Debugf("Begin subQuery by ID %s, requestID: %s", subID, requestID)
			for id, dataValue := range data {
				subdata, err := FetchData(subID, dataValue, token, []string{}, requestID)
				if err != nil {
					data[id]["sub_query_"+strconv.Itoa(k)] = err.Error()
				} else {
					data[id]["sub_query_"+strconv.Itoa(k)] = subdata
				}
			}
		}
		return data, nil
	}
}

func mergeMapsCheckConflict(map1, map2 map[string]interface{}) (map[string]interface{}, error) {
	for k, v := range map1 {
		if _, exists := map2[k]; exists {
			return nil, fmt.Errorf("用户变量 '%s' 与系统变量冲突，请联系管理员", k)
		}
		map2[k] = v
	}

	return map2, nil
}