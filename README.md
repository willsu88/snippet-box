## Motivation 
In my first attempt at my [Go project](https://github.com/willsu88/findlabs-go), I was able to hack away and create a 3 layer project REST - Server - DB. However, many
issues persisted in that project that I wasn't sure how to tackle in the Go way. Namely:
1. No separation of concerns. HTTP handler code, my server side logic, and DB code were all clunked together in the same functions and file. 
2. No management of global states.

## Learnings
In this new project, taken inspiration and learnings from the the book [Let's Go](https://lets-go.alexedwards.net/) by Alex Edwards, I've implemented  quite a few useful patterns to structure my Go code to be as idiomatic and clean as possible. Below are some key insights.

### MVC Pattern
Separating the code into 3 main categories: main.go, handlers.go, and /models/*.go. This is similar to the famous MVC pattern used in Java/Spring.  

1. Controller: **main.go** handles all the initialiaztion -> connecting to the DB, routing paths to handler functions, and keeps track of any "global" states.
2. View: **handlers.go**  consists of each of the functions that main.go routes to. I think of these as the entry points to a HTTP request
3.  Model: **/models/*.go** consists of the underlying interactions with the DB and also the "objects" or "structs" that reflect the data from the DB. The good thing about this is that I can further separate these into directories by the actual DB. So this makes it easier for me to swap DB's in the future, and my Controller/View code aren't dependent directly on this.

### Dependency Injection
Since I've separated my code into multiple files/packages, there are bound to be times when I need to share some global state among these files. One way I've accomplished this is to create an **application** struct that lives in the main file. This struct serves as the object that will keep track of global states like
- the opened connection to the DB
- the specific log objects we use
- caches to HTML templates

See here:
```
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}
```
Then in my handler and helper functions, I simply declare the function header like so

```
func (app *application) home
```
so that these functions have direct access to this application struct and thus access to the states it manages. 

**Note:** This was the first time I've used this type of function declaration. In my mind, I imagine this to be Go's way of doing OOP-like functionalities. However, instead of declaring classes and the methods in the same file, I can now have the functionalities spread across different files as needed.



