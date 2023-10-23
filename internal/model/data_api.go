package model  
  
import "time"  
  
type DataAPI struct {  
	ID           int    `db:"id"`  
	ApiName      string `db:"apiname"`  
	ProjectID    string    `db:"project_id"`  
	DbID         string    `db:"db_id"`  
	Joint        string `db:"joint"`  
	SourceSQL    string `db:"sourcesql"`  
	Memo         string `db:"memo"`  
	IsDel        int    `db:"is_del"`  
	CreatedStime *time.Time `db:"created_stime"`  
	ModifiedStime *time.Time `db:"modified_stime"`  
	CreatedUser  string `db:"created_user"`  
	ModifiedUser string `db:"modified_user"`  
}  
