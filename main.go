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
)


func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
}

func dbFunc(c *gin.Context) {
    if _, err := db.Exec("select * from player"); err != nil {
        c.String(http.StatusInternalServerError,
            fmt.Sprintf("Error inserting institution: %q", err))
        return
    }

    c.String(http.StatusOK, fmt.Sprintf("finished!"))

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



    router.GET("/db", dbFunc)


    router.Run(":" + port)

}

