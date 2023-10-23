package model  

import "time"  
  
type DBSource struct {  
	ID           int    `db:"id"`  
	DbName       string `db:"dbname"`  
	IP           string `db:"ip"`  
	Port         int    `db:"port"`  
	User         string `db:"user"`  
	Pwd          string `db:"pwd"`  
	Memo         string `db:"memo"`  
	IsDel        int    `db:"is_del"`  
	CreatedStime *time.Time `db:"created_stime"`  
	ModifiedStime *time.Time `db:"modified_stime"`  
	CreatedUser  string `db:"created_user"`  
	ModifiedUser string `db:"modified_user"`  
}  
