package main

import (
	//"fmt"
	"log"
	"net/http"
	"database/sql"
	"encoding/json"
	 _ "github.com/go-sql-driver/mysql"
)

type Book struct {
	Bookname string  `json:"bookname"`
	Aname string     `json:"Aname"`
}

type Res struct {
	Status int
	Message error
}

func checkErr(err error) {
    if err != nil {
        panic(err)
	}
}

func add_book(db *sql.DB, b,a string) error {
	insert,err := db.Query("INSERT INTO books (book_name,author_name) VALUES(?,?)",b,a)
	checkErr(err)
	defer insert.Close()
	return err
}


func main(){
	//create db connection
	db, err := sql.Open("mysql","root:Mysql-pinki1@tcp(127.0.0.1:3306)/library")
	checkErr(err)
	defer db.Close()


	f:=http.FileServer(http.Dir("/home/pratiksha/lib_api/static1"))
	g:=http.FileServer(http.Dir("/home/pratiksha/lib_api/static2"))


	http.Handle("/",g)
	http.Handle("/add/", http.StripPrefix("/add",f))
	http.Handle("/display/", http.StripPrefix("/display",g))


	http.HandleFunc("/add",func(w http.ResponseWriter, r *http.Request) {
    
	    var bb Book
	      err := json.NewDecoder(r.Body).Decode(&bb)

	      if err != nil {
	        http.Error(w, err.Error(), http.StatusBadRequest)
	        return
	      }
	      
	      var rr Res
	   
	      rr.Message = add_book(db,bb.Bookname,bb.Aname)
	      if rr.Message != nil{
	      rr.Status = 1
	      } else {
	      rr.Status = 0
	  	  }
	      json.NewEncoder(w).Encode(rr)
    })

  

	http.HandleFunc("/display",func(w http.ResponseWriter, r *http.Request) {
	  
		var b Book
		s := make([]Book,0)
	    list,err := db.Query("select book_name,author_name from books order by book_name")
	    checkErr(err)
	    for list.Next() {
	        // for each row, scan the result into our tag composite object
	        err = list.Scan(&b.Bookname,&b.Aname)
	        checkErr(err)
			s = append(s,b)      
	    }
	    json.NewEncoder(w).Encode(s)
	})

	log.Fatal(http.ListenAndServe(":3001",nil))
}

 
// url-> http://localhost:3001/add/		||		http://localhost:3001/display
// http.Handle("/",http.FileServer(http.Dir("/home/pratiksha/go_project/static")))
//fmt.Fprintln(w, "Welcome Server : "+r.URL.String()+"\n")
/*
func display_list(db *sql.DB){
	var b book
	list := db.Query("select * from books order by book_name")

    for list.Next() {
        // for each row, scan the result into our tag composite object
        err = list.Scan(&b.bookName,&b.authorName)
        checkErr(err)
        
    }
}
*/
