package db

import (
	"fmt"
	"shensuanzi/conf"
	"shensuanzi/datastruct"
	"shensuanzi/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DBHandler struct {
	mysqlEngine *xorm.Engine
}

func CreateDBHandler(isWeb bool) *DBHandler {
	dbHandler := new(DBHandler)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", conf.Server.DB_UserName, conf.Server.DB_Pwd, conf.Server.DB_IP, conf.Server.DB_Name)
	engine, err := xorm.NewEngine("mysql", dsn)
	errhandle(err)
	err = engine.Ping()
	errhandle(err)
	//日志打印SQL
	engine.ShowSQL(true)
	//设置连接池的空闲数大小
	maxIdleConns := 40
	if conf.Common.Mode == conf.Debug || isWeb {
		maxIdleConns = 5
	}
	engine.SetMaxIdleConns(maxIdleConns)
	SyncDB(engine)
	initData(engine)
	dbHandler.mysqlEngine = engine
	//go timerTask(engine)
	return dbHandler
}

func SyncDB(engine *xorm.Engine) {
	arr := make([]interface{}, 0)
	var err error

	arr = append(arr, new(datastruct.ColdUserInfo))
	arr = append(arr, new(datastruct.HotUserInfo))
	arr = append(arr, new(datastruct.WXUserInfo))
	arr = append(arr, new(datastruct.WXFTInfo))
	arr = append(arr, new(datastruct.AdInfo))
	arr = append(arr, new(datastruct.ImgNav))
	arr = append(arr, new(datastruct.AppraisedInfo))
	arr = append(arr, new(datastruct.AppraisedMark))
	arr = append(arr, new(datastruct.AppraisedImgs))
	arr = append(arr, new(datastruct.UserAttention))
	arr = append(arr, new(datastruct.InviteUserRegister))

	arr = append(arr, new(datastruct.UserDepositInfo))
	arr = append(arr, new(datastruct.FTDepositInfo))
	arr = append(arr, new(datastruct.UserGoldChange))
	arr = append(arr, new(datastruct.BalanceInfo))
	arr = append(arr, new(datastruct.UserDrawCashInfo))
	arr = append(arr, new(datastruct.FTDrawCashInfo))
	arr = append(arr, new(datastruct.DrawCashParams))
	arr = append(arr, new(datastruct.UserOrderInfo))
	arr = append(arr, new(datastruct.UserOrderCheck))
	arr = append(arr, new(datastruct.UserOrderRefundFinished))

	arr = append(arr, new(datastruct.UserOrderRefunding))
	arr = append(arr, new(datastruct.UserOrderRightsFinished))
	arr = append(arr, new(datastruct.UserOrderRighting))
	arr = append(arr, new(datastruct.OrderInfoMsg))
	arr = append(arr, new(datastruct.OrderRightsFinishedMsg))
	arr = append(arr, new(datastruct.UserOrderRefundMsg))
	arr = append(arr, new(datastruct.UserRegisterMsg))
	arr = append(arr, new(datastruct.FTRegisterMsg))
	arr = append(arr, new(datastruct.FTOrderRefundMsg))
	arr = append(arr, new(datastruct.ColdFTInfo))

	arr = append(arr, new(datastruct.HotFTInfo))
	arr = append(arr, new(datastruct.ShopInfo))
	arr = append(arr, new(datastruct.ShopImgs))
	//arr = append(arr, new(datastruct.FTOnlineTime))
	arr = append(arr, new(datastruct.ProductInfo))
	arr = append(arr, new(datastruct.Authentication))
	arr = append(arr, new(datastruct.PayConsumerProtection))
	arr = append(arr, new(datastruct.FTDndList))
	arr = append(arr, new(datastruct.ChattingList))
	arr = append(arr, new(datastruct.ChatEndList))
	arr = append(arr, new(datastruct.FTAccountChange))

	arr = append(arr, new(datastruct.FTQuickReply))
	arr = append(arr, new(datastruct.FTMarkInfo))
	//arr = append(arr, new(datastruct.DomainInfo))
	arr = append(arr, new(datastruct.CustomerServiceInfo))
	//arr = append(arr, new(datastruct.ServerInfo))
	arr = append(arr, new(datastruct.SensitiveWord))

	//.DropTables(arr...) //test

	err = engine.Sync2(arr...)
	errhandle(err)

	// role := createRoleData()
	// _, err = engine.Insert(&role)
	// errhandle(err)
	// webUser := createLoginData(datastruct.AdminLevelID)
	// _, err = engine.Insert(&webUser)
	// errhandle(err)
}

/*
func createRoleData() []datastruct.Role {
	admin := datastruct.Role{
		Id:   datastruct.AdminLevelID,
		Desc: "admin",
	}
	guest := datastruct.Role{
		Id:   datastruct.NormalLevelID,
		Desc: "normal",
	}

	return []datastruct.Role{admin, guest}
}

func createLoginData(adminLevelID int) []datastruct.WebUser {
	now_time := time.Now().Unix()
	admin := datastruct.WebUser{
		Name:      "管理员",
		LoginName: "admin",
		Pwd:       "123@s678",
		RoleId:    adminLevelID,
		Token:     commondata.CommonDataInfo.UniqueId(),
		CreatedAt: now_time,
		UpdatedAt: now_time,
	}
	return []datastruct.WebUser{admin}
}*/

func initData(engine *xorm.Engine) {
	execStr := fmt.Sprintf("ALTER DATABASE %s CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;", conf.Server.DB_Name)
	_, err := engine.Exec(execStr)
	errhandle(err)

	_, err = engine.Exec("ALTER TABLE hot_f_t_info CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)

	_, err = engine.Exec("ALTER TABLE cold_f_t_info CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)

	_, err = engine.Exec("ALTER TABLE f_t_quick_reply CONVERT TO CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)

	// execStr = fmt.Sprintf("ALTER TABLE user_info AUTO_INCREMENT = %d", datastruct.UserIdStart)
	// _, err = engine.Exec(execStr)
	// errhandle(err)

	_, err = engine.Exec("ALTER TABLE hot_f_t_info CHANGE auto_reply auto_reply VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)

	_, err = engine.Exec("ALTER TABLE cold_f_t_info CHANGE nick_name nick_name VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)

	_, err = engine.Exec("ALTER TABLE f_t_quick_reply CHANGE reply reply VARCHAR(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;")
	errhandle(err)
}

func errhandle(err error) {
	if err != nil {
		log.Fatal("db error is %v", err.Error())
	}
}
