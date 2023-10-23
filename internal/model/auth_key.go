// internal/model/auth_key.go  
package model  

import "time"  

type AuthKey struct {  
	ID           int    `db:"id"`  
	AuthKey      string `db:"authkey"`  
	//Joint        string `db:"joint"`  
	ProjectID    int    `db:"project_id"`  
	ApiIDs       string `db:"api_ids"`  
	Memo         string `db:"memo"`  
	IsDel        int    `db:"is_del"`  
	CreatedStime *time.Time `db:"created_stime"`  
	ModifiedStime *time.Time `db:"modified_stime"`  
	CreatedUser  string `db:"created_user"`  
	ModifiedUser string `db:"modified_user"`  
}  
