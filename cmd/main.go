package main

import (
	"log"
	"time"

	"github.com/classAndrew/govalor/services"

	"github.com/classAndrew/govalor/controllers"

	"github.com/classAndrew/govalor/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(err)
	}

	r := gin.Default()

	models.ConnectDatabase()
	// let's gin know that anytime /static is referenced, it's using ./static relative
	// when localhost:8080/abc, will try to look for index.html static vs staticFile will do an entire directory
	// the first argument is when user enters url, the second is where to look for on the server side
	// r.Static("/abc", "./web/static")
	r.StaticFile("/", "./web/static/index.html")

	// :name is path parameter, not query string parameter
	r.GET("/usertotalxp/:guild/:name", controllers.FindUserTotalXP)
	r.GET("/usertotalxp/:guild", controllers.FindUserTotalXP)
	r.GET("/usertotalxp", controllers.FindUserTotalXP)
	r.POST("/usertotalxp", controllers.CreateUserTotalXP)
	r.PATCH("/usertotalxp/:name", controllers.UpdateUserTotalXP)

	log.SetPrefix("[Valor Engine] ")
	log.Println("Starting Server")
	enemies := []string{"GYP ON TOP", "Aequitas", "Avicia", "IceBlue Team", "BlueStoneGroup",
		"Bovemists", "Guardian of Wynn", "Eden", "Cyphrus Code", "Fuzzy Spiders", "Nerfuria",
		"Nethers Ascent", "ShadowFall", "Spectral Cabbage", "The Dark Phoenix", "The Mage Legacy",
		"TheNoLifes", "Wheres The Finish", "Ultra Violet", "The Multiverse", "Elit Magyar Legio",
		"WrathOfTheFallen",
	}
	go func() {
		time.Sleep(time.Second * 60 * 5)                                          // take five minutes before starting each up
		go services.UpdateMemberXP([]string{"Titans%20Valor"}, time.Second*60*30) // thirty minutes
		go services.CheckActivity(enemies, time.Second*60*60)                     // hourly
	}()

	r.Run(":8080")
}
