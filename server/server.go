package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type pet struct {
	id             int
	family         string
	sex            string
	name           string
	age            int
	breed          string
	checkIn        string
	checkOut       string
	idRoom         int
	idOwner        int
	idHealthReport int
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "write.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		family := r.FormValue("family")
		sex := r.FormValue("sex")
		name := r.FormValue("name")
		breed := r.FormValue("breed")
		checkIn := r.FormValue("checkIn")
		checkOut := r.FormValue("checkOut")

		id, err := strconv.Atoi(r.FormValue("id"))
		age, err := strconv.Atoi(r.FormValue("age"))
		idRoom, err := strconv.Atoi(r.FormValue("idRoom"))
		idOwner, err := strconv.Atoi(r.FormValue("idOwner"))
		idHealthReport, err := strconv.Atoi(r.FormValue("idHealthReport"))
		if err != nil {
			panic(err)
		}

		db, err := sql.Open("mysql", "root:@/zoohome")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO `pets`(`id`, `family`, `sex`, `name`, `age`, `breed`, `check-in`, `check-out`, `id_room`, `id_owner`, `id_health-report`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			id, family, sex, name, age, breed, checkIn, checkOut, idRoom, idOwner, idHealthReport)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w, "Data added successfully.")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "read.html")
	case "POST":
		//Read data
		db, err := sql.Open("mysql", "root:@/zoohome")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		rows, err := db.Query("SELECT * FROM `pets` LIMIT 10")
		if err != nil {
			panic(err)
		}
		defer rows.Close()
		pets := []pet{}

		for rows.Next() {
			p := pet{}
			err := rows.Scan(&p.id, &p.family, &p.sex, &p.name, &p.age, &p.breed, &p.checkIn, &p.checkOut, &p.idRoom, &p.idOwner, &p.idHealthReport)
			if err != nil {
				fmt.Println(err)
				continue
			}
			pets = append(pets, p)
		}
		for _, p := range pets {
			fmt.Fprintf(w, "id: %d, family: %s, sex: %s, name: %s, age: %d, breed: %s, checkIn: %s, checkOut: %s, idRoom: %d, idOwner: %d, idHealthReport: %d\n",
				p.id, p.family, p.sex, p.name, p.age, p.breed, p.checkIn, p.checkOut, p.idRoom, p.idOwner, p.idHealthReport)
		}

		fmt.Fprintf(w, "\nData readed successfully.")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "index.html")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/write", writeHandler)
	router.HandleFunc("/read", readHandler)
	router.HandleFunc("/", mainHandler)
	http.Handle("/", router)

	fmt.Printf("Server started!\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
