package main

import (
	"context"
	"fmt"
	"html/template"
	"mvcweb/connection"
	"mvcweb/helper"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectData struct {
	Name,Description,Image,Duration string
	StartDate,EndDate pgtype.Date
	Technologies[]string	
}

var projects []ProjectData 




func main() {
	router := mux.NewRouter()

	connection.DatabaseConnect()

	directory := http.Dir("./public")
	fileServer := http.FileServer(directory)

    router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	// router
	// get
	router.HandleFunc("/", getHome).Methods("GET")
	router.HandleFunc("/form-add-project", getAddProject).Methods("GET")
	router.HandleFunc("/form-edit-project/{index}", getEditProject).Methods("GET")
	router.HandleFunc("/contact-me", getContactMe).Methods("GET")
	router.HandleFunc("/project/{projectId}", getProjectDetail).Methods("GET")
	// post
	router.HandleFunc("/add-project", postAddProject).Methods("POST")
	router.HandleFunc("/update-project/{index}", updateProject).Methods("POST")
	router.HandleFunc("/delete-project/{index}", deleteProject).Methods("POST")
	


	fmt.Println("running on port 5000")
	http.ListenAndServe("localhost:5000", router)

}

// show homepage, where showing project from postgre database
func getHome(w http.ResponseWriter, r *http.Request) {

	data, err := connection.Conn.Query(context.Background(), "SELECT name,start_date,end_date,description,technologies,image FROM public.tb_projects;")

	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	var dataResult []ProjectData

	var project ProjectData
	for data.Next() {

		var err = data.Scan(&project.Name, &project.StartDate, &project.EndDate, &project.Description, &project.Technologies, &project.Image)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		
		project.Duration = helper.GetDuration(project.StartDate.Time.Format("2006-01-02"), project.EndDate.Time.Format("2006-01-02"))
		dataResult = append(dataResult, project)
	}
	

	
	var view, templErr = template.ParseFiles("views/index.html")	
	if err != nil {
		panic(templErr.Error())
	}
	view.Execute(w, dataResult)
}

func getContactMe(w http.ResponseWriter, r *http.Request) {
	var view, err = template.ParseFiles("views/contact.html")	
	if err != nil {
		panic(err.Error())
	}
	view.Execute(w, nil)
}

func getProjectDetail(w http.ResponseWriter, r *http.Request) {
	projectIndex, indexError := strconv.Atoi(mux.Vars(r)["projectId"]);
	if indexError != nil {
		panic(indexError.Error())
	}
	data := projects[projectIndex]
	var view,err = template.ParseFiles("views/projectDetail.html")
	if err != nil {
		panic(err.Error())
	}
	view.Execute(w, data)

}

func postAddProject(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	

	http.Redirect(w,r,"/form-add-project", http.StatusFound)
}


func getAddProject(w http.ResponseWriter, r *http.Request) {
	var view, err = template.ParseFiles("views/project.html")	
	if err != nil {
		panic(err.Error())
	}

	view.Execute(w, nil)
}

func getEditProject(w http.ResponseWriter, r *http.Request) {
	indexVars := mux.Vars(r)["index"]
	projectIndex, parseErr := strconv.Atoi(indexVars)
	if parseErr != nil {
		panic(parseErr.Error())
	}
	currentData := projects[projectIndex]
	var view, err = template.ParseFiles("views/edit-project.html")
	if err != nil {
		panic(err.Error())
	}
	data := map[string]interface{} {
		"data": currentData,
		"index": indexVars,
	}

	view.Execute(w, data)
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	parseErr := r.ParseForm()
	// newData := r.PostForm;
	// projectIndex := mux.Vars(r)["index"]
	
	if parseErr != nil {
		panic(parseErr.Error())
	}
	// i, indexErr := strconv.Atoi(projectIndex)

	// if indexErr != nil {
	// 	panic(indexErr.Error())
	// }

	// projects[i].Name = newData.Get("name")
	// projects[i].StartDate = newData.Get("start-date")
	// projects[i].EndDate = newData.Get("end-date")
	// projects[i].Description = newData.Get("description")
	
	http.Redirect(w,r,"/",http.StatusFound)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	projectIndex := mux.Vars(r)["index"]
	
	i, indexErr := strconv.Atoi(projectIndex)

	if indexErr != nil {
		panic(indexErr.Error())
	}

	projects = append(projects[:i], projects[i+1:]...)

	http.Redirect(w,r,"/",http.StatusFound)
}








