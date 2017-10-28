package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/gin-gonic/gin"

"net/http"
    "os"
    "log"


    "fmt"
)

var (
    db *sql.DB
    id int
    usernameRes string
    passwordRes string
)


func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func dbFunc(c *gin.Context) {
    if _, err := db.Exec("INSERT INTO institutions (name,institution_id) VALUES ('Citi', 'ins_5')"); err != nil {
        c.String(http.StatusInternalServerError,
            fmt.Sprintf("Error inserting institution: %q", err))
        return
    }

    c.String(http.StatusOK, fmt.Sprintf("finished!"))

}

func handleLogin(c *gin.Context) {
    username := c.PostForm("username")
    password := c.PostForm("password")

    err := db.QueryRow("SELECT * FROM users WHERE username=$1 and password=$2", username, password).Scan(&id, &usernameRes, &passwordRes)

    switch {
    case err == sql.ErrNoRows:
        //no user. invalid username password combo.
        response := gin.H{
            "status" : http.StatusNotFound,
            "error" : "Invalid username and password combination",
        }
        c.JSON(http.StatusInternalServerError, response)
        break
    case err != nil:
        response := gin.H{
            "status" : http.StatusInternalServerError,
            "error" : "Something is wrong!",
        }
        c.JSON(http.StatusInternalServerError, response)
        break
    default:
        response := gin.H{
            "status" : http.StatusOK,
            "id" : id,
        }
        c.JSON(http.StatusOK, response)
    }
}


func handleRegistration(c *gin.Context) {

    username := c.PostForm("username")
    password := c.PostForm("password")

    response := queryForUser(username, password)


    c.JSON(http.StatusOK, gin.H{"response": response,})

}

func queryForUser(username string, password string) string  {
    err := db.QueryRow("SELECT * FROM users WHERE username=$1 and password=$2", username, password).Scan(&id, &usernameRes, &passwordRes)

    switch {
    case err == sql.ErrNoRows:
        //no user. invalid username password combo.
        //response := gin.H{
        //    "status" : http.StatusNotFound,
        //    "error" : "Invalid username and password combination",
        //}
        return "good to go"
        break
    case err != nil:
        //response := gin.H{
        //    "status" : http.StatusInternalServerError,
        //    "error" : "Something is wrong!",
        //}
        return "error"
        break
    default:
        //response := gin.H{
        //    "status" : http.StatusOK,
        //    "id" : id,
        //}
        return "user already exists"

    }

    return "no switch"
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
    var err error
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }


    db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))

    if err != nil {
        log.Fatalf("Error opening databae: %q", err)
    }

    router := gin.New()
    router.Use(gin.Logger())
    router.Use(CORSMiddleware())
    router.LoadHTMLGlob("templates/*.tmpl.html")
    router.Static("/static", "static")

    router.GET("/setup", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl.html", nil)
    })
    router.POST("/register", handleRegistration)
    router.GET("/db", dbFunc)
    router.POST("/login", handleLogin)

    router.Run(":" + port)

}

