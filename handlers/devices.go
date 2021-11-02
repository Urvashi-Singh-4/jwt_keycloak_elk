package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"token_based_auth/logger"
	"token_based_auth/middleware"
	"token_based_auth/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllDevices(ctx *gin.Context) {

	authorizationKeys := ctx.Request.Header.Get("Authorization")
	authorized := middleware.CheckAuthorised(authorizationKeys, "d.r")
	if !authorized {
		logger.Write("error", "Not Authorized: User didn't have required permissions to read devices")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You dont't have permission to read devices",
		})
		return
	}

	logger.Write("success", "User Authorized to read all devices")
	var result []model.Device
	err := model.GetAllDevices(&result)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occured while fetching all devices data",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
	logger.Write("success", "All device data fetched")
}

func AddNewDevice(ctx *gin.Context) {

	authorizationKeys := ctx.Request.Header.Get("Authorization")
	authorized := middleware.CheckAuthorised(authorizationKeys, "d.c")
	if !authorized {
		logger.Write("error", "Not Authorized: User didn't have required permissions to create device")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You dont't have permission to create devices",
		})
		return
	}

	logger.Write("success", "User Authorized to create new devices")
	createDevice := make(map[string]interface{})
	err := ctx.BindJSON(&createDevice)
	if err != nil {
		logger.Write("error", "error: Incorrect insert data format")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect insert data format",
		})
		return
	}

	createDevice["id"] = uuid.Must(uuid.NewRandom()).String()

	jsonbyteData, _ := json.Marshal(createDevice)
	jsonStrData := string(jsonbyteData)

	if _, ok := createDevice["sensor_data"]; ok {
		if !middleware.CheckAuthorised(authorizationKeys, "d.i") {
			logger.Write("error", "Not Authorized: User didn't have required permissions to insert sensor data")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "You don't have permission to insert sensor data",
			})
			return
		}
	}

	err = model.InsertNewDevice(jsonStrData)
	if err != nil {
		logger.Write("error", err.Error())
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while inserting new device details",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"device_id": createDevice["id"],
		"message":   "Device inserted successfully",
	})
	logger.Write("success", "New device inserted")
}

func DeleteDevices(ctx *gin.Context) {
	authorizationKeys := ctx.Request.Header.Get("Authorization")
	authorized := middleware.CheckAuthorised(authorizationKeys, "d.d")

	if !authorized {
		logger.Write("error", "Not Authorized: User didn't have required permissions to delete devices")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You dont't have permission to delete devices",
		})
		return
	}
	logger.Write("success", "User authorized to delete devices")
	dID := ctx.DefaultQuery("deviceID", "")
	if dID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid device ID",
		})
		return
	}

	err, message := model.DeleteDeviceByID(dID)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"device_id": dID,
		"message":   message,
	})
	logger.Write("success", "Device data deleted")
}

func UdpateNewDevice(ctx *gin.Context) {
	authorizationKeys := ctx.Request.Header.Get("Authorization")
	authorized := middleware.CheckAuthorised(authorizationKeys, "d.u")

	if !authorized {
		logger.Write("error", "Not Authorized: User didn't have required permissions to update device data")
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "You dont't have permission to update devices",
		})
		return
	}
	logger.Write("success", "User authorized to update device data")

	dID := ctx.DefaultQuery("deviceID", "")
	if dID == "" {
		logger.Write("error", "Invalid device id")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid device id",
		})
		return
	}

	updateData := make(map[string]interface{})
	err := ctx.BindJSON(&updateData)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect update data format",
		})
		return
	}

	if _, ok := updateData["sensor_data"]; ok {
		if !middleware.CheckAuthorised(authorizationKeys, "d.i") {
			logger.Write("error", "Not Authorized: User didn't have required permissions to insert sensor data")
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "You don't have permission to insert sensor data",
			})
			return
		}
	}

	err, message := model.UpdateDeviceByID(dID, updateData)
	if err != nil {
		logger.Write("error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": message,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"device_id": dID,
		"message":   message,
	})
	logger.Write("success", "Device data updated")
}
