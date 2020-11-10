package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/role", Get)

	router.GET("/role/:id", GetOne)

	router.POST("/role", Post)

	router.PUT("/role/:id", Put)

	router.DELETE("/role/:id", Delete)

	router.Run(":8080")
}

// 取得全部資料
func Get(c *gin.Context) {
	c.JSON(http.StatusOK, Data)
}

// 取得單一筆資料
func GetOne(c *gin.Context) {
	strID := c.Param("id")
	if id, err := strconv.Atoi(strID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for _, model := range Data {
			if uint(id) == model.ID {
				c.JSON(http.StatusOK, model)
			}
		}
	}

	c.AbortWithStatus(http.StatusNotFound)
}

// 新增資料
func Post(c *gin.Context) {

	jsonIRole := &Role{}
	if err := c.ShouldBindJSON(&jsonIRole); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if jsonIRole.Name == "" || jsonIRole.Summary == "" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	lastData := Data[len(Data)-1]
	jsonIRole.ID = uint(lastData.ID + 1)
	Data = append(Data, *jsonIRole)
	c.JSON(http.StatusOK, *jsonIRole)
}

type RoleVM struct {
	ID      uint   `json:"id"`      // Key
	Name    string `json:"name"`    // 角色名稱
	Summary string `json:"summary"` // 介紹
}

// 更新資料, 更新角色名稱與介紹
func Put(c *gin.Context) {
	roleVM := &RoleVM{}
	strID := c.Param("id")
	if err := c.ShouldBindJSON(&roleVM); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if id, err := strconv.Atoi(strID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		roleVM.ID = uint(id)
		for index, model := range Data {
			if roleVM.ID == model.ID {
				Data[index].Name = roleVM.Name
				Data[index].Summary = roleVM.Summary
				c.JSON(http.StatusOK, Data[index])
				return
			}
		}
	}

	c.AbortWithStatus(http.StatusNotFound)
}

// 刪除資料
func Delete(c *gin.Context) {
	strID := c.Param("id")
	if id, err := strconv.Atoi(strID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		for index, model := range Data {

			if uint(id) == model.ID {
				if Data, err = delAt(Data, index); err != nil {
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				} else {
					c.AbortWithStatus(http.StatusOK)
					return
				}
			}
		}
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func delAt(slice []Role, i int) (result []Role, err error) {
	if i > len(slice) {
		err = errors.New("i can't be larger than len(slice)")
		return
	}
	tempSlice := append(make([]Role, 0, len(slice)), slice...)
	result = append(tempSlice[:i], tempSlice[i+1:]...)
	return
}
