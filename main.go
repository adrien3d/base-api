package main

import (
	"github.com/adrien3d/jump-technical-test/models"
	"github.com/adrien3d/jump-technical-test/server"
	"github.com/adrien3d/jump-technical-test/services"
	"github.com/adrien3d/jump-technical-test/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	api := &server.API{Router: gin.Default(), Config: viper.New()}

	// Configuration setup
	err := api.SetupViper()
	utils.CheckErr(err)

	// Email sender setup
	api.EmailSender = services.NewEmailSender(api.Config)
	api.TextSender = services.NewTextSender(api.Config)

	// Database setup
	dbType := api.Config.GetString("db_type")
	switch dbType {
	case "mongo":
		_, err := api.SetupMongoDatabase()
		if err == nil {
			utils.Log(nil, "info", "SetupMongoDatabase OK")
		} else {
			utils.Log(nil, "error", "SetupMongoDatabase KO:", err)
		}
		utils.CheckErr(err)
		//defer session.Close()

		err = api.SetupMongoIndexes()
		if err == nil {
			utils.Log(nil, "info", "SetupMongoIndexes OK")
		} else {
			utils.Log(nil, "error", "SetupMongoIndexes KO:", err)
		}
		utils.CheckErr(err)

		// Seeds setup
		err = api.SetupMongoSeeds()
		utils.CheckErr(err)
		if err == nil {
			utils.Log(nil, "info", "SetupMongoSeeds OK")
		} else {
			utils.Log(nil, "error", "SetupMongoSeeds KO:", err)
		}
		utils.CheckErr(err)

	case "postgresql":
		db, err := api.SetupPostgreDatabase()
		utils.CheckErr(err)

		db.AutoMigrate(&models.Organization{})
		db.AutoMigrate(&models.Group{})
		db.AutoMigrate(&models.User{})

		err = api.SetupPostgreSeeds()
		utils.CheckErr(err)
	}

	// Router setup
	api.SetupRouter()

	logrus.Infoln("SetupRouter OK")
	err = api.Router.Run(api.Config.GetString("host_address"))
	if err != nil {
		logrus.Infoln("api.Router.Run OK")
	}
	utils.CheckErr(err)
}
