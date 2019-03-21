package app

import (
	"shensuanzi/datastruct"
	"shensuanzi/handle"

	"github.com/gin-gonic/gin"
)

func checkFtToken(c *gin.Context, handle *handle.AppHandler) (int, string, bool) {
	tokens, isExist := c.Request.Header["Apptoken"]
	tf := false
	var token string
	var ft_id int
	var isBlackList bool
	if isExist {
		token = tokens[0]
		if token != "" {
			ft_id, tf, isBlackList = handle.IsExistFt(token)
			if tf && isBlackList {
				// url := eventHandler.GetBlackListRedirect()
				url := "http://www.baidu.com"
				c.JSON(200, gin.H{
					"code": datastruct.Redirect,
					"data": url,
				})
				return -1, "", false
			}
		}
	}
	if !tf {
		c.JSON(200, gin.H{
			"code": datastruct.TokenError,
		})
	}
	return ft_id, token, tf
}

func isExistFtPhone(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/isexistphone/:phone"
	isExistPhone(r, handle, url, true)
}

func isExistPhone(r *gin.Engine, handle *handle.AppHandler, url string, isFT bool) {
	r.GET(url, func(c *gin.Context) {
		phone := c.Param("phone")
		if phone == "" {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := handle.IsExistPhone(phone, isFT)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func isExistFtNickName(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/isexistnickname/:nickname"
	r.GET(url, func(c *gin.Context) {
		nickname := c.Param("nickname")
		if nickname == "" {
			c.JSON(200, gin.H{
				"code": datastruct.ParamError,
			})
			return
		}
		data, code := handle.IsExistNickName(nickname)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFtMarkInfo(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/mark"
	r.GET(url, func(c *gin.Context) {
		data, code := handle.GetFtMarkInfo()
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func ftRegister(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/register"
	r.POST(url, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": handle.FtRegister(c),
		})
	})
}

func ftRegisterWithID(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/registerwithid"
	r.POST(url, func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": handle.FtRegisterWithID(c),
		})
	})
}

func ftLogin(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/login"
	r.POST(url, func(c *gin.Context) {
		data, code := handle.FtLogin(c)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func ftInfo(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/info"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtInfo(id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func updateFtInfo(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/info"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.UpdateFtInfo(c, id),
		})
	})
}

func updateFtMark(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/mark"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.UpdateFtMark(c, id),
		})
	})
}

func updateFtIntroduction(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/introduce"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.UpdateFtIntroduction(c, id),
		})
	})
}

func updateFtAutoReply(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/autoreply"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.UpdateFtAutoReply(c, id),
		})
	})
}

func getFtIntroduction(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/introduce"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtIntroduction(id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFtAutoReply(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/autoreply"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtAutoReply(id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func ftSubmitIdentity(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/submitidentity"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.FtSubmitIdentity(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getAppraised(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/appraised/:pageindex/:pagesize"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetAppraised(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFtUnReadMsgCount(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/unreadcount"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtUnReadMsgCount(id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFtSystemMsg(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/sysmsg/:pageindex/:pagesize"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtSystemMsg(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFtDndList(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/dndlist/:pageindex/:pagesize"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtDndList(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getProduct(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/product"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetProduct(id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func removeFtDndList(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/dndlist"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.RemoveFtDndList(c, id),
		})
	})
}
func editProduct(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/editproduct"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.EditProduct(c, id)
		if code == datastruct.Sensitive {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func removeProduct(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/rmproduct"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.RemoveProduct(c, id),
		})
	})
}

func sortProducts(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/sortproducts"
	r.POST(url, func(c *gin.Context) {
		_, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.SortProducts(c),
		})
	})
}

func createFakeAppraised(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/fakeappraised"
	r.POST(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.CreateFakeAppraised(c, id),
		})
	})
}

func getAllFtOrder(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/allorder/:pageindex/:pagesize"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetAllFtOrder(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getAmountList(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/amountlist/:datatype/:pageindex/:pagesize"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetAmountList(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFtDrawCashInfo(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/drawinfo/:pageindex/:pagesize"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtDrawCashInfo(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFtAccountChangeInfo(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/accountchange/:pageindex/:pagesize"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFtAccountChangeInfo(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getIncomeList(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/incomelist/:datatype/:pageindex/:pagesize"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetIncomeList(c, id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFinance(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/finance"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetFinance(id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getQRcode(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/qrcode"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetQRcode(id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getFtDrawCashParams(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/drawparams"
	r.GET(url, func(c *gin.Context) {
		_, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetDrawCashParams(datastruct.FT)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func getProducts(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/products"
	r.GET(url, func(c *gin.Context) {
		id, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		data, code := handle.GetProducts(id)
		if code == datastruct.NULLError {
			c.JSON(200, gin.H{
				"code": code,
				"data": data,
			})
		} else {
			c.JSON(200, gin.H{
				"code": code,
			})
		}
	})
}

func isAgreeRefund(r *gin.Engine, handle *handle.AppHandler) {
	url := "/app/ft/isagreerefund"
	r.POST(url, func(c *gin.Context) {
		_, _, tf := checkFtToken(c, handle)
		if !tf {
			return
		}
		c.JSON(200, gin.H{
			"code": handle.IsAgreeRefund(c),
		})
	})
}

// func ftIsOnline(r *gin.Engine, handle *handle.AppHandler) {
// 	url := "/app/ft/online"
// 	r.POST(url, func(c *gin.Context) {
// 		id, _, tf := checkFtToken(c, handle)
// 		if !tf {
// 			return
// 		}
// 		c.JSON(200, gin.H{
// 			"code": handle.FtIsOnline(id),
// 		})
// 	})
// }

func FtRegisterRoutes(r *gin.Engine, handle *handle.AppHandler) {
	isExistFtPhone(r, handle)
	isExistFtNickName(r, handle)
	getFtMarkInfo(r, handle)
	ftRegister(r, handle)
	ftRegisterWithID(r, handle)
	ftLogin(r, handle)
	ftInfo(r, handle)
	updateFtInfo(r, handle)
	updateFtMark(r, handle)
	updateFtIntroduction(r, handle)
	getFtIntroduction(r, handle)
	getFtAutoReply(r, handle)
	updateFtAutoReply(r, handle)
	ftSubmitIdentity(r, handle)
	getAppraised(r, handle)
	getFtUnReadMsgCount(r, handle)
	getFtSystemMsg(r, handle)
	getFtDndList(r, handle)
	removeFtDndList(r, handle)
	editProduct(r, handle)
	removeProduct(r, handle)
	getProduct(r, handle)
	sortProducts(r, handle)
	getAllFtOrder(r, handle)
	createFakeAppraised(r, handle)
	isAgreeRefund(r, handle)
	getFinance(r, handle)
	getProducts(r, handle)
	getAmountList(r, handle)
	getIncomeList(r, handle)
	getQRcode(r, handle)
	getFtDrawCashParams(r, handle)
	getFtDrawCashInfo(r, handle)
	getFtAccountChangeInfo(r, handle)
	//ftIsOnline(r, handle)
}
