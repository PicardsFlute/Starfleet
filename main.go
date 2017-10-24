package main
//all go files must be part of package main if you want to use their functionality without importing

import (
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"fmt"
	"github.com/gorilla/handlers"
	"os"
 	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	//"os/user"
)

var (
	tpl *template.Template
	db *gorm.DB
	err error
)

type Person struct {
	Email string
	Password string
}





// Then, initialize the session manager
func init() {
	dbPassword := os.Getenv("PG_DATABASE_PW")
	db, err = gorm.Open("postgres", "host=127.0.0.1 dbname=Starfleet sslmode=disable password="+dbPassword)
	if err != nil {
		fmt.Println("Cannot connect to database...")
		fmt.Println("DB Error: ", err)
	}



	db.AutoMigrate(&User{}, &Student{})
}



type User struct {
	UserID uint `gorm:"primary_key"`
	UserEmail string `gorm:"type:varchar(20);unique"`
	UserPassword string `gorm:"type:varchar(300)"`
	FirstName string `gorm:"type:varchar(50)"`
	LastName string `gorm:"type:varchar(50)"`
	UserType int
}

type Student struct {
	StudentID uint `gorm:"primary_key"`
	User  User `gorm:"ForeignKey:UserRefer"`
	UserRefer uint
}


func main() {







	routes := mux.NewRouter()
	tpl = template.Must(template.ParseGlob("templates/*"))
	routes.PathPrefix("/style").Handler(http.StripPrefix("/style/",http.FileServer(http.Dir("style"))))
	routes.PathPrefix("/public").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	routes.HandleFunc("/",index)
	//routes.HandleFunc("/about/{number}", about)
	routes.HandleFunc("/login", loginPage).Methods("GET")
	routes.HandleFunc("/login/{num}", loginUser).Methods("POST")
	routes.HandleFunc("/user/{id:[0-9]+}", displayUser).Methods("GET")



	// USED FOR HEROKU
	//http.ListenAndServe(":" + os.Getenv("PORT"), routes)

	//USED FOR LOCAL, only use one
	http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout,routes))

	//defer db.Close(), want to keep db connection open
}


// routes for site

func index(w http.ResponseWriter, r *http.Request){

	tpl.ExecuteTemplate(w, "index", nil)

}

func loginPage(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w,"login",nil)
}

func loginUser(w http.ResponseWriter, r *http.Request){

	vars := mux.Vars(r)
	fmt.Println(vars)

	formEmail := r.FormValue("email")
	formPassword :=	r.FormValue("password")

	// Handle db checks here, if they are valid, render the user template and pass in some data,
	// otherwise,
	// render the login template with an error message

	p := Person{}

	fmt.Println("Email: ", formEmail)
	fmt.Println("Password: ", formPassword)

	u := User{}
	db.First(&u)
	if u.UserEmail != "" {
		displayUser(w,r)
	} else {
		p.Email = "Not Found"
		p.Password = "Not found"
	}



	tpl.ExecuteTemplate(w,"user",p)
}

func displayUser(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w, "user", nil)
}


/*
func about(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r) //returns a mapping responses
	personId := vars["number"] //get map with key id number
	if num, _ := strconv.Atoi(personId); num > 3 {
		p := Person{"Bob", 4}
		tpl.ExecuteTemplate(w, "index.html", p)
	}else {
		p := Person{"Steve", 2}
		tpl.ExecuteTemplate(w, "index.html", p)
	}
}
func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.public", file))
	}
	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}
*/