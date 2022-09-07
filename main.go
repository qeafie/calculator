package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Task struct {
	FirstNumber, SecondNumber float64
	Sign                      string
}
type TaskForm struct {
	*Task
	Answer float64
}

var templates = make(map[string]*template.Template, 3)

func loadTemplates() {

	//массив имён шаблонов для парсинга
	templateName := [2]string{"form", "answer"}

	//проходим по именам и пробуем пропарсить файл с этим именем

	for index, name := range templateName {
		t, err := template.ParseFiles("layout.html", name+".html")

		//если успешно то добавляем в  нашу мапу шаблонов
		if err == nil {
			templates[name] = t
			fmt.Println("Loaded template", index, name)
		}
	}
}

// обработчик страницы формы
func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, TaskForm{
			Task: &Task{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()

		firstpars, err := strconv.ParseFloat(request.Form["first"][0], 32)
		secondpars, err := strconv.ParseFloat(request.Form["second"][0], 32)
		if err != nil {
			panic(err)
		}
		responseData := Task{
			FirstNumber:  firstpars,
			SecondNumber: secondpars,
			Sign:         request.Form["sign"][0],
		}
		answer := 0.0
		if responseData.Sign == "+" {
			answer = float64(responseData.FirstNumber + responseData.SecondNumber)
		} else if responseData.Sign == "-" {
			answer = float64(responseData.FirstNumber - responseData.SecondNumber)
		} else if responseData.Sign == "*" {
			answer = float64(responseData.FirstNumber * responseData.SecondNumber)
		} else if responseData.Sign == "/" {
			answer = float64(responseData.FirstNumber / responseData.SecondNumber)
		}
		fmt.Println(answer)
		templates["answer"].Execute(writer, answer)

	}

}

func main() {
	//загрузка шаблонов
	loadTemplates()

	//регистрация обработчиrков
	http.HandleFunc("/", formHandler)

	//создание http сервера

	err := http.ListenAndServe(":5000", nil)

	if err != nil {
		fmt.Println(err)
	}

}
