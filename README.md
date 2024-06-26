## Motivation 
In my first attempt at my [Go project](https://github.com/willsu88/findlabs-go), I was able to hack away and create a 3 layer project REST - Server - DB. However, many
issues persisted in that project that I wasn't sure how to tackle in the Go way. Namely:
1. No separation of concerns. HTTP handler code, my server side logic, and DB code were all clunked together in the same functions and file. 
2. No management of global states.

## Learnings
In this new project, taken inspiration and learnings from the the book [Let's Go](https://lets-go.alexedwards.net/) by Alex Edwards, I've implemented  quite a few useful patterns to structure my Go code to be as idiomatic and clean as possible. Below are some key insights.

### Separation of Concerns
Separating the code into 3 main categories: `main.go`, `handlers.go`, and `/models/*.go`. This is similar to the famous MVC pattern I've used in Java/Spring.  

1. Controller: **main.go** handles all the initialiaztion -> connecting to the DB, routing paths to handler functions, and keeps track of any "global" states.
2. View: **handlers.go**  consists of each of the functions that main.go routes to. I think of these as the entry points to a HTTP request
3.  Model: **/models/*.go** consists of the underlying interactions with the DB and also the "objects" or "structs" that reflect the data from the DB. The good thing about this is that I can further separate these into directories by the actual DB. So this makes it easier for me to swap DB's in the future, and my Controller/View code aren't dependent directly on this.

### Dependency Injection
Since I've separated my code into multiple files/packages, there are bound to be times when I need to share some global state among these files. One way I've accomplished this is to create an **application** struct that lives in the main file. This struct serves as the object that will keep track of global states like
- connection to the DB to snippets
- connection to the DB to users
- the specific log objects we use
- session manager

See here:
```
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	users         *mysql.UserModel
	session       *sessions.Session
}
```
Then in my handler and helper functions, I simply declare the function header like so

```
func (app *application) home
```
so that these functions have direct access to this application struct and thus access to the states it manages. 

**Note:** This was the first time I've used this type of function declaration. In my mind, I imagine this to be Go's way of doing OOP-like functionalities. However, instead of declaring classes and the methods in the same file, I can now have the functionalities spread across different files as needed.

### Middleware Pattern
I like to think of the Middleware Pattern to be synonymous to the famous **Chain of Responsibility Pattern**. The Go web application is just a chain of `ServeHTTP()` methods that are called one after another. There are quite a few routes that my mux helps set up to handle different requests. But what's common among all of these requests is that I need a way to
1. handle panic recovery for the goroutine serving the HTTP request
2. log HTTP requests as they come in
3. handle some secure headers

It would not make sense to copy this code across each of the different routes and their respective functions. Instead, I can chain them together like so:
```
panicRecovery -> logRequests -> secureHeaders -> serveMux -> handler functions
```
So that the process of panic recovery, logs, secure headers are always handled. For this, I used a third party library `github.com/justinas/alice` to help me chain them together.
Eg
```
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	...
	return standardMiddleware.Then(mux)
```

Note: Another important detail I learned is the keyword `defer`. I have used this to handle defering the closing of a DB conenction. But in this particular instance, I used it to defer the `recover()` method. I thought this was a pretty neat feature of Go-lang. In that, in an event of a panic, as Go unwinds the stack, it'll eventually hit this recover function and handles the necessary tasks as I imposed.

### Session Management
Session management is something I knew in theory, but have never really implemented, so it was a good experience to learn this. As a note to myself, sesion management is different from state management (think React.js states). Session is handled server side, while state is working on client side.

In this project, one of the needs to have session management is that when a user creates a snippet, the user should recieve a confirmation message. To create the session manager, I used a third party library `github.com/golangcollege/sessions` and injected into my main application struct. This way I could use it as necessary in my handler functions.

For instsance, upon a successful creation of a snippet, in my `createSnippet()` function, I will run to put this state into the session.
```
app.session.Put(r, "flash", "Snippet successfully created!")
```
In the `showSnippet()` function, I then retrieve this state as necessary using. 
```
app.session.PopString(r, "flash")
```
**Note:** In my code, I actually placed the PopString function in a helper function in `helpers.go` that adds any default data to the `showSnippet()` function. This just makes it easier to scale and add any other default data I need for the form.

**Tradeoff:** I understand the need for session managers for other cases. However, for this particular use case, I wonder if it's possible to just handle this client side? For, instance upon creation of a succesful user, the client side will receive a 200 Status ok. Using that Status ok, the client code (say React) can just re-render as necessary on their side. Perhaps sessions become more important when the data being passed around has higher security concerns.

### User Authentication


