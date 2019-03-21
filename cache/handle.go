package cache

import (
	"shensuanzi/datastruct"
	"shensuanzi/log"
	"shensuanzi/tools"

	"github.com/gomodule/redigo/redis"
)

func (handle *CACHEHandler) clearData() {
	conn := handle.GetConn()
	defer conn.Close()
	conn.Do("flushdb")
}

func (handle *CACHEHandler) GetConn() redis.Conn {
	conn := handle.redisClient.Get()
	return conn
}

/*

//创建充值订单
func (handle *CACHEHandler) CreateOrderForm(pay *datastruct.OrderForm) datastruct.CodeType {
	conn := handle.GetConn()
	defer conn.Close()
	value, isError := tools.InterfaceToString(pay)
	if isError {
		log.Error("CreateOrderForm InterfaceToString err")
		return datastruct.UpdateDataFailed
	}
	conn.Send("MULTI")
	conn.Send("SETNX", pay.Id, value)
	conn.Send("expire", pay.Id, datastruct.WXOrderMaxSec)
	_, err := conn.Do("EXEC")
	if err != nil {
		log.Error("CreateOrderForm payInfo:%v err:%s", pay, err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

//获取充值订单
func (handle *CACHEHandler) GetPayData(payId string) (*datastruct.OrderForm, datastruct.CodeType) {
	conn := handle.GetConn()
	defer conn.Close()
	value, err := redis.String(conn.Do("get", payId))
	if err != nil || value == "" {
		log.Error("GetPayData payId:%v err", payId)
		return nil, datastruct.GetDataFailed
	}
	pay, isError := tools.BytesToOrderForm([]byte(value))
	if isError {
		log.Error("GetPayData BytesToOrderForm err")
		return nil, datastruct.GetDataFailed
	}
	return pay, datastruct.NULLError
}


func (handle *CACHEHandler) UpdateBlackList(token string, isBlackList int) {
	conn := handle.GetConn()
	defer conn.Close()

	args := make([]interface{}, 0, 3)
	args = append(args, token)

	args = append(args, datastruct.IsBlackListField)
	args = append(args, isBlackList)

	_, err := conn.Do("hset", args...)

	if err != nil {
		log.Debug("CACHEHandler UpdateBlackList err:%s", err.Error())
	}
}

func (handle *CACHEHandler) DeletedKeys(keys []interface{}) {
	conn := handle.GetConn()
	defer conn.Close()
	_, err := conn.Do("del", keys...)
	if err != nil {
		log.Debug("CACHEHandler DeletedKeys err:%s", err.Error())
	}
}

func (handle *CACHEHandler) GetTmpAvatars() []string {

	return nil
}

func (handle *CACHEHandler) GetTmpUsers() []*datastruct.TmpUser {

	return nil
}

*/

func (handle *CACHEHandler) SetFtToken(conn redis.Conn, data *datastruct.FtRedisData) {
	key := data.Token

	args := make([]interface{}, 0, 5)
	args = append(args, key)
	args = append(args, datastruct.Reidis_IdField)
	args = append(args, data.FtId)
	args = append(args, datastruct.Reidis_AccountStateField)
	args = append(args, int(data.AccountState))

	_, err := conn.Do("hmset", args...)
	if err != nil {
		log.Debug("CACHEHandler SetFtToken err:%s", err.Error())
	}
}

func (handle *CACHEHandler) SetUserToken(conn redis.Conn, data *datastruct.UserRedisData) {
	key := data.Token

	args := make([]interface{}, 0, 5)
	args = append(args, key)
	args = append(args, datastruct.Reidis_IdField)
	args = append(args, data.UserId)
	args = append(args, datastruct.Reidis_AccountStateField)
	args = append(args, int(data.AccountState))

	_, err := conn.Do("hmset", args...)
	if err != nil {
		log.Debug("CACHEHandler SetFtToken err:%s", err.Error())
	}
}

const ftOnlineKey = "ft_online"

func (handle *CACHEHandler) SetFtOnline(conn redis.Conn, ft_id_arr []interface{}) datastruct.CodeType {
	args := make([]interface{}, 0)
	args = append(args, ftOnlineKey)
	args = append(args, ft_id_arr...)
	_, err := conn.Do("SADD", args...)
	if err != nil {
		log.Error("CACHEHandler SetFtOnline err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

// func (handle *CACHEHandler) GetOnlineFtIds(conn redis.Conn, ft_id int) datastruct.CodeType {
// 	// value, err := redis.Values(conn.Do("hmget", key, datastruct.Reidis_IdField, datastruct.Reidis_AccountStateField))
// }

func (handle *CACHEHandler) FtOffline(conn redis.Conn, ft_id_arr []interface{}) {
	args := make([]interface{}, 0)
	args = append(args, ftOnlineKey)
	args = append(args, ft_id_arr...)
	_, err := conn.Do("SREM", args...)
	if err != nil {
		log.Debug("CACHEHandler FtOffline err:%s", err.Error())
	}
}

func (handle *CACHEHandler) AddExpire(conn redis.Conn, token string) {
	_, err := conn.Do("expire", token, datastruct.RedisExpireTime)
	if err != nil {
		log.Debug("expire token:%s err:%s", token, err.Error())
		return
	}
}

func (handle *CACHEHandler) IsExistFtWithConn(conn redis.Conn, key string) (int, bool, bool) {
	value, err := redis.Values(conn.Do("hmget", key, datastruct.Reidis_IdField, datastruct.Reidis_AccountStateField))
	length := len(value)
	if err != nil || length == 0 {
		return -1, false, false
	}
	var id int
	var accountState int
	for i := 0; i < length; i++ {
		x_value := value[i]
		if x_value == nil {
			return -1, false, false
		}
		tmp := x_value.([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			id = tools.StringToInt(str)
		case 1:
			accountState = tools.StringToInt(str)
		}
	}
	isBlackList := false
	if accountState == int(datastruct.BlackList) {
		isBlackList = true
	}
	return id, true, isBlackList
}

func (handle *CACHEHandler) IsExistUserWithConn(conn redis.Conn, key string) (int64, bool, bool) {
	value, err := redis.Values(conn.Do("hmget", key, datastruct.Reidis_IdField, datastruct.Reidis_AccountStateField))
	length := len(value)
	if err != nil || length == 0 {
		return -1, false, false
	}
	var id int64
	var accountState int
	for i := 0; i < length; i++ {
		x_value := value[i]
		if x_value == nil {
			return -1, false, false
		}
		tmp := x_value.([]byte)
		str := string(tmp[:])
		switch i {
		case 0:
			id = tools.StringToInt64(str)
		case 1:
			accountState = tools.StringToInt(str)
		}
	}
	isBlackList := false
	if accountState == int(datastruct.BlackList) {
		isBlackList = true
	}
	return id, true, isBlackList
}
