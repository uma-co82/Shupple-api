package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/uma-co82/Shupple-api/src/api/service"
)

type AdminController struct{}

/************************************************************
 *                         ADMIN                            *
 ************************************************************/
func (adminController AdminController) GetAllUser(c *gin.Context) {
	var userService service.UserService
	p, err := userService.GetAllUser(c)

	if err != nil {
		c.JSON(err.(*service.Error).Code, err)
	} else {
		c.JSON(200, p)
	}
}
