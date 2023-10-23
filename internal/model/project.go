// internal/model/project.go  
package model  

import "time"  

  
type Project struct {  
	ID           int    `db:"id"`  
	ProjectName  string `db:"project_name"`  
	ProjectType  string `db:"project_type"` 
	Memo         string `db:"memo"`  
	IsDel        int    `db:"is_del"`  
	CreatedStime *time.Time `db:"created_stime"`  
	ModifiedStime *time.Time `db:"modified_stime"`  
	CreatedUser  string `db:"created_user"`  
	ModifiedUser string `db:"modified_user"`  
}  
