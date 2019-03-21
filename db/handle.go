package db

import (
	"fmt"
	"shensuanzi/commondata"
	"shensuanzi/datastruct"
	"shensuanzi/datastruct/important"
	"shensuanzi/log"
	"shensuanzi/tools"
	"strings"
	"time"

	"github.com/go-xorm/xorm"
)

func (handle *DBHandler) GetDirectDownloadApp() string {
	engine := handle.mysqlEngine
	column_name := "down_load_u_i_addr"
	sql := "select down_load_u_i_addr from domain_info where id = ?"
	results, err := engine.Query(sql, datastruct.DefaultId)
	if err != nil || len(results) <= 0 {
		return ""
	}
	return string(results[0][column_name][:])
}

func (handle *DBHandler) IsExistPhone(phone string, isFT bool) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	var sql string
	if isFT {
		sql = "select 1 from cold_f_t_info where phone = ? limit 1"
	} else {
		sql = "select 1 from cold_user_info where phone = ? limit 1"
	}
	results, err := engine.Query(sql, phone)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	count := len(results)
	isExist := false
	if count > 0 {
		isExist = true
	}
	return isExist, datastruct.NULLError
}

func (handle *DBHandler) IsExistNickName(nickname string) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select 1 from cold_f_t_info where nick_name = ? limit 1"
	results, err := engine.Query(sql, nickname)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	count := len(results)
	isExist := false
	if count > 0 {
		isExist = true
	}
	return isExist, datastruct.NULLError
}
func (handle *DBHandler) GetFtMarkInfo() (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select ftm.desc from f_t_mark_info ftm order by id desc"
	results, err := engine.Query(sql)
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	arr := make([]string, 0, len(results))
	for _, v := range results {
		arr = append(arr, string(v["desc"][:]))
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) FtRegister(body *datastruct.FTRegisterBody) datastruct.CodeType {
	return ftRegister(handle.mysqlEngine, body, nil)
}

func (handle *DBHandler) FtRegisterWithID(body *datastruct.FTRegisterWithIDBody) datastruct.CodeType {
	return ftRegister(handle.mysqlEngine, &body.FTRegisterBody, body)
}

func ftRegister(engine *xorm.Engine, body *datastruct.FTRegisterBody, IDbody *datastruct.FTRegisterWithIDBody) datastruct.CodeType {
	nowTime := time.Now().Unix()
	cold_ft := new(datastruct.ColdFTInfo)
	cold_ft.CreatedAt = nowTime
	cold_ft.Introduction = body.Desc
	cold_ft.NickName = body.NickName
	cold_ft.Phone = body.Phone
	cold_ft.Pwd = body.Pwd
	cold_ft.Avatar = body.Avatar
	cold_ft.Registration = body.Platform
	cold_ft.AuthState = datastruct.Authing

	if IDbody != nil {
		cold_ft.IdentityCard = IDbody.Identity
		cold_ft.IdFrontCover = IDbody.IdFrontCover
		cold_ft.IdBehindCover = IDbody.IdBehindCover
		cold_ft.ActualName = IDbody.ActualName
	}

	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	_, err := session.InsertOne(cold_ft)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister insert ColdFTInfo :%s", err.Error())
		rollbackError(str, session)
		return datastruct.NickNamePhoneIsUsed
	}

	hot_ft := new(datastruct.HotFTInfo)
	hot_ft.FTId = cold_ft.Id
	hot_ft.Mark = body.Mark

	auth := new(datastruct.Authentication)
	auth.FTId = cold_ft.Id
	if IDbody != nil {
		auth.IdCardState = datastruct.IdCardSubmited
	} else {
		auth.IdCardState = datastruct.IdCardNotSubmit
	}
	auth.IsCP = false
	auth.IsHR = false

	//create shop
	shop := new(datastruct.ShopInfo)
	shop.FTId = cold_ft.Id

	_, err = session.Insert(hot_ft, auth, shop)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister insert HotFTInfo and Authentication :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtRegister Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func rollbackError(err_str string, session *xorm.Session) {
	log.Error("will rollback,err_str:%v", err_str)
	session.Rollback()
}

func (handle *DBHandler) FtLogin(body *datastruct.FtLoginBody) (*datastruct.RespFtLogin, datastruct.CodeType) {
	engine := handle.mysqlEngine
	cold_ft := new(datastruct.ColdFTInfo)
	has, err := engine.Where("phone=?", body.Phone).Get(cold_ft)
	if err != nil {
		log.Error("FtLogin get err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	if !has {
		return nil, datastruct.NotRegisterPhone
	}
	if cold_ft.Pwd != body.Pwd {
		return nil, datastruct.PwdError
	}
	var code datastruct.CodeType
	var rs *datastruct.RespFtLogin
	rs = nil
	switch cold_ft.AuthState {
	case datastruct.Authing:
		code = datastruct.AuthingCode
	case datastruct.AuthFailed:
		code = datastruct.AuthFailedCode
	default:
		rs, code = ftLoginSucceed(cold_ft.Id, engine)
	}
	return rs, code
}

func ftLoginSucceed(ft_id int, enging *xorm.Engine) (*datastruct.RespFtLogin, datastruct.CodeType) {
	token := commondata.UniqueId()
	nowTime := time.Now().Unix()
	im_id := fmt.Sprintf("ft_%d", ft_id)
	im_privatekey, code := tools.AccountGenForIM(im_id, important.IM_SDK_APPID)
	if code != datastruct.NULLError {
		return nil, code
	}
	sql := "update hot_f_t_info set token = ?,login_time = ?,i_m_private_key = ? where f_t_id=?"
	_, err := enging.Exec(sql, token, nowTime, im_privatekey, ft_id)
	if err != nil {
		log.Error("ftLoginSucceed err:%s", err.Error())
		return nil, datastruct.UpdateDataFailed
	}
	respFtInfo, code := getFtInfo(ft_id, enging)
	rs := new(datastruct.RespFtLogin)
	rs.FtInfo = respFtInfo
	rs.Token = token
	rs.IMPrivateKey = im_privatekey
	return rs, code
}

func (handle *DBHandler) GetFtInfo(ft_id int) (*datastruct.RespFtInfo, datastruct.CodeType) {
	return getFtInfo(ft_id, handle.mysqlEngine)
}

func getFtInfo(ft_id int, engine *xorm.Engine) (*datastruct.RespFtInfo, datastruct.CodeType) {
	sql := "select cold.avatar,cold.nick_name,hot.enable_free,hot.mark,auth.is_c_p,auth.is_h_r,auth.id_card_state from cold_f_t_info cold join hot_f_t_info hot on cold.id = hot.f_t_id join authentication auth on auth.f_t_id = cold.id where id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil {
		log.Error("GetFtInfo err0: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	if len(results) <= 0 {
		log.Error("GetFtInfo err1: len(results) is zero")
		return nil, datastruct.GetDataFailed
	}
	mp := results[0]
	info := new(datastruct.RespFtInfo)
	info.Id = ft_id
	info.Avatar = string(mp["avatar"][:])
	info.NickName = string(mp["nick_name"][:])
	info.Mark = string(mp["mark"][:])
	info.EnableFree = tools.StringToBool(string(mp["enable_free"][:]))
	authIcon := new(datastruct.AuthIcon)
	authIcon.IsCP = tools.StringToBool(string(mp["is_c_p"][:]))
	authIcon.IsHR = tools.StringToBool(string(mp["is_h_r"][:]))
	authIcon.IdCard = tools.StringToIdCardState(string(mp["id_card_state"][:]))
	info.AuthIcon = authIcon
	return info, datastruct.NULLError
}

func (handle *DBHandler) GetFtDataWithToken(token string) (*datastruct.FtRedisData, bool) {
	engine := handle.mysqlEngine
	sql := "select f_t_id,account_state from hot_f_t_info where token = ?"
	results, err := engine.Query(sql, token)
	if err != nil || len(results) <= 0 {
		return nil, false
	}
	redis_data := new(datastruct.FtRedisData)
	redis_data.Token = token

	mp := results[0]
	redis_data.FtId = tools.StringToInt(string(mp["f_t_id"][:]))
	redis_data.AccountState = datastruct.AccountState(tools.StringToInt(string(mp["account_state"][:])))
	sql = "update hot_f_t_info set login_time = ? where token = ?"
	_, err = engine.Exec(sql, time.Now().Unix(), token)
	if err != nil {
		log.Error("GetFtDataWithToken err:%s", err.Error())
		return nil, false
	}

	return redis_data, true
}

func (handle *DBHandler) GetUserDataWithToken(token string) (*datastruct.UserRedisData, bool) {
	engine := handle.mysqlEngine
	sql := "select user_id,account_state from hot_user_info where token = ?"
	results, err := engine.Query(sql, token)
	if err != nil || len(results) <= 0 {
		return nil, false
	}
	redis_data := new(datastruct.UserRedisData)
	redis_data.Token = token

	mp := results[0]
	redis_data.UserId = tools.StringToInt64(string(mp["user_id"][:]))
	redis_data.AccountState = datastruct.AccountState(tools.StringToInt(string(mp["account_state"][:])))
	sql = "update hot_user_info set login_time = ? where token = ?"
	_, err = engine.Exec(sql, time.Now().Unix(), token)
	if err != nil {
		log.Error("GetUserDataWithToken err:%s", err.Error())
		return nil, false
	}

	return redis_data, true
}

func (handle *DBHandler) UpdateFtInfo(body *datastruct.UpdateFtInfoBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	var err error
	var results []map[string][]byte
	sql := "select avatar from cold_f_t_info where id = ?"
	results, err = engine.Query(sql, ft_id)
	if err != nil {
		log.Error("DBHandler->UpdateFtInfo get avatar err:", err.Error())
		return datastruct.UpdateDataFailed
	}
	avatar := ""
	if len(results) > 0 {
		avatar = string(results[0]["avatar"][:])
	}
	cold_ft := new(datastruct.ColdFTInfo)
	cold_ft.Avatar = body.Avatar
	cold_ft.NickName = body.NickName
	_, err = engine.Where("id=?", ft_id).Cols("nick_name", "avatar").Update(cold_ft)
	if err != nil {
		log.Error("DBHandler->UpdateFtInfo err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	if avatar != "" {
		commondata.DeleteOSSFileWithUrl(avatar)
	}
	return datastruct.NULLError
}

func (handle *DBHandler) UpdateFtMark(body *datastruct.UpdateFtMarkBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	hot_ft := new(datastruct.HotFTInfo)
	hot_ft.Mark = body.Mark
	_, err := engine.Where("f_t_id=?", ft_id).Cols("mark").Update(hot_ft)
	if err != nil {
		log.Error("DBHandler->UpdateFtMark err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) UpdateFtIntroduction(body *datastruct.UpdateFtIntroductionBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()
	cold_ft := new(datastruct.ColdFTInfo)
	cold_ft.Introduction = body.Desc
	_, err := session.Where("id=?", ft_id).Cols("introduction").Update(cold_ft)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction err0: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	shop_id, code := getShopIdWithSession(session, ft_id)
	if code != datastruct.NULLError {
		return code
	}

	sql := "select img_url from shop_imgs where shop_id = ?"
	results, err := session.Query(sql, shop_id)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction get imgs_url err")
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	will_DeleteImgs := make([]string, 0, len(results))
	for _, v := range results {
		will_DeleteImgs = append(will_DeleteImgs, string(v["img_url"][:]))
	}

	sql = "delete from shop_imgs where shop_id = ?"
	_, err = session.Exec(sql, shop_id)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction err1: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = "insert into shop_imgs(shop_id,img_url) values "
	values := ""
	tmp := ""
	for i := 0; i < len(body.Imgs); i++ {
		if i == 0 {
			tmp = fmt.Sprintf("(%d,'%s')", shop_id, body.Imgs[i])
		} else {
			tmp = fmt.Sprintf(",(%d,'%s')", shop_id, body.Imgs[i])
		}
		values += tmp
	}

	sql = sql + values
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction err2: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	if len(will_DeleteImgs) > 0 {
		for _, v := range will_DeleteImgs {
			commondata.DeleteOSSFileWithUrl(v)
		}
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetFtIntroduction(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select introduction from cold_f_t_info where id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->GetFtIntroduction err0")
		return nil, datastruct.GetDataFailed
	}
	desc := string(results[0]["introduction"][:])

	sql = "select id from shop_info where f_t_id = ?"
	results, err = engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->GetFtIntroduction err1")
		return nil, datastruct.GetDataFailed
	}
	shop_id := tools.StringToInt(string(results[0]["id"][:]))

	sql = "select img_url from shop_imgs where shop_id = ? order by id asc"
	results, err = engine.Query(sql, shop_id)
	if err != nil || len(results) < 0 {
		log.Error("DBHandler->GetFtIntroduction err2")
		return nil, datastruct.GetDataFailed
	}
	arr := make([]string, 0, len(results))
	for _, v := range results {
		arr = append(arr, string(v["img_url"][:]))
	}
	resp := new(datastruct.UpdateFtIntroductionBody)
	resp.Desc = desc
	resp.Imgs = arr
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetFtAutoReply(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select auto_reply from hot_f_t_info where f_t_id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->GetFtAutoReply err0")
		return nil, datastruct.GetDataFailed
	}
	autoreply := string(results[0]["auto_reply"][:])

	sql = "select reply from f_t_quick_reply where f_t_id = ? order by id asc"
	results, err = engine.Query(sql, ft_id)
	if err != nil || len(results) < 0 {
		log.Error("DBHandler->GetFtIntroduction err1")
		return nil, datastruct.GetDataFailed
	}
	arr := make([]string, 0, len(results))
	for _, v := range results {
		arr = append(arr, string(v["reply"][:]))
	}
	resp := new(datastruct.UpdateFtAutoReplyBody)
	resp.AutoReply = autoreply
	resp.QuickReply = arr
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UpdateFtAutoReply(body *datastruct.UpdateFtAutoReplyBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	hot_ft := new(datastruct.HotFTInfo)
	hot_ft.AutoReply = body.AutoReply
	_, err := session.Where("f_t_id=?", ft_id).Cols("auto_reply").Update(hot_ft)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtAutoReply err0: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	sql := "delete from f_t_quick_reply where f_t_id = ?"
	_, err = session.Exec(sql, ft_id)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtAutoReply err1: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = "insert into f_t_quick_reply(f_t_id,reply) values "
	values := ""
	tmp := ""
	for i := 0; i < len(body.QuickReply); i++ {
		if i == 0 {
			tmp = fmt.Sprintf("(%d,'%s')", ft_id, body.QuickReply[i])
		} else {
			tmp = fmt.Sprintf(",(%d,'%s')", ft_id, body.QuickReply[i])
		}
		values += tmp
	}

	sql = sql + values
	_, err = session.Exec(sql)
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtAutoReply err2: %s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->UpdateFtIntroduction Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	return datastruct.NULLError
}

func (handle *DBHandler) FtSubmitIdentity(body *datastruct.FtIdentity, ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine

	var err error
	var results []map[string][]byte
	sql := "select id_front_cover,id_behind_cover from cold_f_t_info where id = ?"
	results, err = engine.Query(sql, ft_id)
	if err != nil {
		log.Error("DBHandler->FtSubmitIdentity get id_cover err:", err.Error())
		return nil, datastruct.UpdateDataFailed
	}
	last_front_cover := ""
	last_behind_cover := ""
	if len(results) > 0 {
		last_front_cover = string(results[0]["id_front_cover"][:])
		last_behind_cover = string(results[0]["id_behind_cover"][:])
	}

	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	ft_cold := new(datastruct.ColdFTInfo)
	ft_cold.IdentityCard = body.Identity
	ft_cold.ActualName = body.ActualName
	ft_cold.IdBehindCover = body.IdBehindCover
	ft_cold.IdFrontCover = body.IdFrontCover

	_, err = session.Where("id=?", ft_id).Cols("identity_card", "actual_name", "id_behind_cover", "id_front_cover").Update(ft_cold)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtSubmitIdentity err0: %s", err.Error())
		rollbackError(str, session)
		return nil, datastruct.UpdateDataFailed
	}

	auth := new(datastruct.Authentication)
	auth.IdCardState = datastruct.IdCardSubmited
	_, err = session.Where("f_t_id=?", ft_id).Cols("id_card_state").Update(auth)
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtSubmitIdentity err1: %s", err.Error())
		rollbackError(str, session)
		return nil, datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->FtSubmitIdentity Commit :%s", err.Error())
		rollbackError(str, session)
		return nil, datastruct.UpdateDataFailed
	}
	if last_front_cover != "" || last_behind_cover != "" {
		commondata.DeleteOSSFileWithUrl(last_front_cover)
		commondata.DeleteOSSFileWithUrl(last_behind_cover)
	}
	return auth.IdCardState, datastruct.NULLError
}

func (handle *DBHandler) GetAppraised(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select ap.user_id,ap.score,ap.mark,ap.appraised_type,ap.desc,ap.is_anonym,ap.created_at,pr.product_name from appraised_info ap join product_info pr on ap.product_id = pr.id join shop_info sh on sh.id = pr.shop_id where sh.f_t_id = ? LIMIT ?,?"
	results, err := engine.Query(sql, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetAppraised err0: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	rs := make([]*datastruct.RespAppraise, 0, len(results))
	for _, v := range results {
		ap := new(datastruct.RespAppraise)
		is_anonym := tools.StringToBool(string(v["is_anonym"][:]))
		if is_anonym {
			user_id := tools.StringToInt64(string(v["user_id"][:]))
			sql = "select nick_name,avatar from cold_user_info where id = ?"
			sub_results, sub_err := engine.Query(sql, user_id)
			if sub_err != nil || len(sub_results) <= 0 {
				log.Error("DBHandler->GetAppraised err1: the user is not exist, id is %v", user_id)
				return nil, datastruct.GetDataFailed
			}
			ap.Name = string(sub_results[0]["nick_name"][:])
			ap.Avatar = string(sub_results[0]["avatar"][:])
		}
		ap.IsAnonym = is_anonym
		ap.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		ap.Desc = string(v["desc"][:])
		ap.Mark = string(v["mark"][:])
		ap.ProductName = string(v["product_name"][:])
		ap.Score = tools.StringToFloat64(string(v["score"][:]))
		rs = append(rs, ap)
	}
	return rs, datastruct.NULLError
}

func (handle *DBHandler) GetFtUnReadMsgCount(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	where_str := "f_t_read = 0 and f_t_id = ?"
	arr := make([]interface{}, 0)
	arr = append(arr, new(datastruct.OrderInfoMsg))
	arr = append(arr, new(datastruct.OrderRightsFinishedMsg))
	arr = append(arr, new(datastruct.FTRegisterMsg))
	arr = append(arr, new(datastruct.FTOrderRefundMsg))
	var count int64
	count = 0
	for _, v := range arr {
		tmp, err := engine.Where(where_str, ft_id).Count(v)
		if err != nil {
			log.Error("DBHandler->GetFtUnReadMsgCount err:", err.Error())
			return nil, datastruct.GetDataFailed
		}
		count += tmp
	}
	return count, datastruct.NULLError
}

func (handle *DBHandler) GetUserUnReadMsgCount(user_id int) (interface{}, datastruct.CodeType) {

	engine := handle.mysqlEngine
	where_str := "user_read = 0 and user_id = ?"
	arr := make([]interface{}, 0)
	arr = append(arr, new(datastruct.OrderInfoMsg))
	arr = append(arr, new(datastruct.OrderRightsFinishedMsg))
	arr = append(arr, new(datastruct.UserRegisterMsg))
	arr = append(arr, new(datastruct.UserOrderRefundMsg))

	var count int64
	count = 0
	for _, v := range arr {
		tmp, err := engine.Where(where_str, user_id).Count(v)
		if err != nil {
			log.Error("DBHandler->GetUserUnReadMsgCount err:", err.Error())
			return nil, datastruct.GetDataFailed
		}
		count += tmp
	}
	return count, datastruct.NULLError
}

func (handle *DBHandler) GetFtSystemMsg(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	//type 1 f_t_register_msg
	//type 2 f_t_order_refund_msg
	//type 3 order_info_msg
	//type 4 order_rights_finished_msg
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize

	sql := "select * from (SELECT -1 as id,ftr.f_t_read,ftr.created_at,ftr.f_t_id,1 as type,-1 as user_nick_name,-1 as product_name,-1 as handle from f_t_register_msg ftr union all select ftor.id,ftor.f_t_read,ftor.created_at,ftor.f_t_id,2 as type,user_nick_name,product_name,refund_type as handle from f_t_order_refund_msg ftor union all select oi.id,oi.f_t_read,oi.created_at,oi.f_t_id,3 as type,user_nick_name,product_name,-1 as handle from order_info_msg oi union all select orf.id,orf.f_t_read,orf.created_at,orf.f_t_id,4 as type,user_nick_name,product_name,is_agree as handle from order_rights_finished_msg orf) as tmp_msg where f_t_id = ? ORDER BY created_at DESC LIMIT ?,?"
	results, err := engine.Query(sql, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetFtSystemMsg err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	arr := make([]interface{}, 0, len(results))
	var rs interface{}
	for _, v := range results {
		tableType := tools.StringToInt(string(v["type"][:]))
		createdAt := tools.StringToInt64(string(v["created_at"][:]))
		isRead := tools.StringToBool(string(v["f_t_read"][:]))
		switch tableType {
		case 1:
			rrm := new(datastruct.RespRegisterMsg)
			rrm.CreatedAt = createdAt
			rrm.Type = tableType
			rrm.GZH_Name = important.WX_GZH_NAME
			rrm.QRCode = important.WX_GZH_QRCode
			rs = rrm
			if !isRead {
				ftrm := new(datastruct.FTRegisterMsg)
				ftrm.FTRead = true
				engine.Where("f_t_id=?", ft_id).Cols("f_t_read").Update(ftrm)
			}
		case 2:
			rrfm := new(datastruct.RespRefundFTMsg)
			rrfm.CreatedAt = createdAt
			rrfm.OrderId = tools.StringToInt64(string(v["id"][:]))
			rrfm.NickName = string(v["user_nick_name"][:])
			rrfm.ProductName = string(v["product_name"][:])
			rrfm.RefundType = datastruct.UserOrderRefundType(tools.StringToInt(string(v["handle"][:])))
			rrfm.Type = tableType
			rs = rrfm
			if !isRead {
				ftor := new(datastruct.FTOrderRefundMsg)
				ftor.FTRead = true
				engine.Where("id=?", rrfm.OrderId).Cols("f_t_read").Update(ftor)
			}
		case 3:
			toi := new(datastruct.TmpOrderInfoFT)
			toi.CreatedAt = createdAt
			toi.OrderId = tools.StringToInt64(string(v["id"][:]))
			toi.NickName = string(v["user_nick_name"][:])
			toi.ProductName = string(v["product_name"][:])
			toi.Type = tableType
			rs = toi
			if !isRead {
				oi := new(datastruct.OrderInfoMsg)
				oi.FTRead = true
				engine.Where("id=?", toi.OrderId).Cols("f_t_read").Update(oi)
			}
		case 4:
			rfm := new(datastruct.RespRightsFinishedFTMsg)
			rfm.CreatedAt = createdAt
			rfm.OrderId = tools.StringToInt64(string(v["id"][:]))
			rfm.NickName = string(v["user_nick_name"][:])
			rfm.ProductName = string(v["product_name"][:])
			rfm.Type = tableType
			rfm.IsAgree = tools.StringToBool(string(v["handle"][:]))
			rs = rfm
			if !isRead {
				orf := new(datastruct.OrderRightsFinishedMsg)
				orf.FTRead = true
				engine.Where("id=?", rfm.OrderId).Cols("f_t_read").Update(orf)
			}
		}
		arr = append(arr, rs)
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) GetUserSystemMsg(user_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	//type 1 register_msg
	//type 2 user_order_refund_msg
	//type 3 order_info_msg
	//type 4 order_rights_finished_msg
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize

	results, err := engine.Query(datastruct.AllMsgSqlForFt, user_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetUserSystemMsg err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	arr := make([]interface{}, 0, len(results))
	var rs interface{}
	for _, v := range results {
		tableType := tools.StringToInt(string(v["type"][:]))
		isRead := tools.StringToBool(string(v["user_read"][:]))
		switch tableType {
		case 1:
			rrm := new(datastruct.RespRegisterMsg)
			rrm.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rrm.Type = tableType
			rrm.GZH_Name = important.WX_GZH_NAME
			rrm.QRCode = important.WX_GZH_QRCode
			rs = rrm
			if !isRead {
				urm := new(datastruct.UserRegisterMsg)
				urm.UserRead = true
				engine.Where("user_id=?", user_id).Cols("user_read").Update(urm)
			}
		case 2:
			rrfu := new(datastruct.RespRefundUserMsg)
			rrfu.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rrfu.OrderId = tools.StringToInt64(string(v["id"][:]))
			rrfu.NickName = string(v["f_t_nick_name"][:])
			rrfu.ProductName = string(v["product_name"][:])
			rrfu.RefundResultType = datastruct.UserOrderRefundResultType(tools.StringToInt(string(v["handle"][:])))
			rrfu.Type = tableType
			rs = rrfu
			if !isRead {
				uor := new(datastruct.UserOrderRefundMsg)
				uor.UserRead = true
				engine.Where("id=?", rrfu.OrderId).Cols("user_read").Update(uor)
			}
		case 3:
			toi := new(datastruct.TmpOrderInfoFT)
			toi.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			toi.OrderId = tools.StringToInt64(string(v["id"][:]))
			toi.NickName = string(v["f_t_nick_name"][:])
			toi.ProductName = string(v["product_name"][:])
			toi.Type = tableType
			rs = toi
			if !isRead {
				oi := new(datastruct.OrderInfoMsg)
				oi.UserRead = true
				engine.Where("id=?", toi.OrderId).Cols("user_read").Update(oi)
			}
		case 4:
			rfm := new(datastruct.RespRightsFinishedFTMsg)
			rfm.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
			rfm.OrderId = tools.StringToInt64(string(v["id"][:]))
			rfm.NickName = string(v["f_t_nick_name"][:])
			rfm.ProductName = string(v["product_name"][:])
			rfm.Type = tableType
			rfm.IsAgree = tools.StringToBool(string(v["handle"][:]))
			rs = rfm
			if !isRead {
				orf := new(datastruct.OrderRightsFinishedMsg)
				orf.UserRead = true
				engine.Where("id=?", rfm.OrderId).Cols("user_read").Update(orf)
			}
		}
		arr = append(arr, rs)
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) GetFtDndList(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select dnd.id,dnd.created_at,u.avatar,u.nick_name from f_t_dnd_list dnd join cold_user_info u on dnd.user_id = u.id where dnd.f_t_id = ? order by dnd.created_at desc LIMIT ?,?"
	results, err := engine.Query(sql, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetFtDndList err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	rs := make([]*datastruct.RespDndList, 0, len(results))
	for _, v := range results {
		dnd := new(datastruct.RespDndList)
		dnd.Id = tools.StringToInt(string(v["id"][:]))
		dnd.Avatar = string(v["avatar"][:])
		dnd.CreatedAt = tools.StringToInt64(string(v["avatar"][:]))
		dnd.NickName = string(v["nick_name"][:])
		rs = append(rs, dnd)
	}
	return rs, datastruct.NULLError
}

func (handle *DBHandler) RemoveFtDndList(id int, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	_, err := engine.Where("id=? and f_t_id=?", id, ft_id).Delete(new(datastruct.FTDndList))
	if err != nil {
		log.Error("DBHandler->RemoveFtDndList err: %s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func getShopIdWithSession(session *xorm.Session, ft_id int) (int, datastruct.CodeType) {
	sql := "select id from shop_info where f_t_id = ?"
	results, err := session.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		str := fmt.Sprintf("DBHandler->getShopIdWithSession err")
		rollbackError(str, session)
		return -1, datastruct.UpdateDataFailed
	}
	shop_id := tools.StringToInt(string(results[0]["id"][:]))
	return shop_id, datastruct.NULLError
}

func getShopId(engine *xorm.Engine, ft_id int) (int, datastruct.CodeType) {
	sql := "select id from shop_info where f_t_id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->getShopId err")
		return -1, datastruct.UpdateDataFailed
	}
	shop_id := tools.StringToInt(string(results[0]["id"][:]))
	return shop_id, datastruct.NULLError
}

func (handle *DBHandler) EditProduct(body *datastruct.EditProductBody, ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine

	sql := "select word from sensitive_word order by id asc"
	results, err := engine.Query(sql)
	if err != nil {
		log.Error("EditProduct get words err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}

	for _, v := range results {
		tmp := string(v["word"][:])
		if strings.Contains(body.ProductName, tmp) || strings.Contains(body.ProductDesc, tmp) {
			return tmp, datastruct.Sensitive
		}
	}

	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return nil, code
	}

	count, err := engine.Where("shop_id=?", shop_id).Count(new(datastruct.ProductInfo))
	if err != nil {
		return nil, datastruct.GetDataFailed
	}
	if count >= datastruct.MAX_PRODUCT_COUNT {
		return nil, datastruct.MaxCreateCount
	}

	isAdd := true
	if body.Id > 0 {
		isAdd = false
	}
	nowTime := time.Now().Unix()
	pro := new(datastruct.ProductInfo)
	pro.UpdatedAt = nowTime
	pro.IsHidden = body.IsHidden
	pro.Price = body.Price
	pro.ProductDesc = body.ProductDesc
	pro.ProductName = body.ProductName
	pro.SortId = 0
	if isAdd {
		pro.CreatedAt = nowTime
		pro.ShopId = shop_id
		_, err := engine.InsertOne(pro)
		if err != nil {
			log.Error("DBHandler->EditProduct insert err:%v", err.Error())
			return nil, datastruct.UpdateDataFailed
		}
	} else {
		_, err := engine.Where("id = ? and shop_id = ?", body.Id, shop_id).Cols("product_name", "product_desc", "price", "is_hidden", "updated_at", "sort_id").Update(pro)
		if err != nil {
			log.Error("DBHandler->EditProduct update err:%v", err.Error())
			return nil, datastruct.UpdateDataFailed
		}
	}
	return nil, datastruct.NULLError
}

func (handle *DBHandler) RemoveProduct(id int, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return code
	}
	_, err := engine.Where("id=? and shop_id=?", id, shop_id).Delete(new(datastruct.ProductInfo))
	if err != nil {
		log.Error("DBHandler->RemoveProduct err: %s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetProduct(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return nil, code
	}
	resp := new(datastruct.RespProductInfo)
	onsale_count, onsale, code := getSaleInfo(shop_id, true, engine)
	if code != datastruct.NULLError {
		return nil, code
	}
	offsale_count, offsale, code := getSaleInfo(shop_id, false, engine)
	if code != datastruct.NULLError {
		return nil, code
	}
	resp.OnSaleCount = onsale_count
	resp.OnSale = onsale
	resp.OffSaleCount = offsale_count
	resp.OffSale = offsale
	resp.Total = onsale_count + offsale_count
	return resp, datastruct.NULLError
}

func getSaleInfo(shop_id int, isOnSale bool, engine *xorm.Engine) (int, []*datastruct.TmpProductInfo, datastruct.CodeType) {
	isHidden := 1
	if isOnSale {
		isHidden = 0
	}
	sql := "select id,product_name,product_desc,price from product_info where shop_id = ? and is_hidden = ? order by sort_id asc,created_at asc"
	results, err := engine.Query(sql, shop_id, isHidden)
	if err != nil {
		log.Error("DBHandler->GetProduct err0: %s", err.Error())
		return -1, nil, datastruct.GetDataFailed
	}
	sale_count := len(results)
	sale := make([]*datastruct.TmpProductInfo, 0, sale_count)
	for _, v := range results {
		tmp := new(datastruct.TmpProductInfo)
		tmp.Id = tools.StringToInt(string(v["id"][:]))
		tmp.Price = tools.StringToFloat64(string(v["price"][:]))
		tmp.ProductDesc = string(v["product_desc"][:])
		tmp.ProductName = string(v["product_name"][:])
		sale = append(sale, tmp)
	}
	return sale_count, sale, datastruct.NULLError
}

func (handle *DBHandler) SortProducts(pids []int) datastruct.CodeType {
	engine := handle.mysqlEngine
	front_sql := "insert into product_info(id,sort_id) values "
	values := ""
	back_sql := " on duplicate key update sort_id=values(sort_id)"
	tmp := ""
	for i := 0; i < len(pids); i++ {
		if i == 0 {
			tmp = fmt.Sprintf("(%d,%d)", pids[i], i)
		} else {
			tmp = fmt.Sprintf(",(%d,%d)", pids[i], i)
		}
		values += tmp
	}
	sql := front_sql + values + back_sql
	_, err := engine.Exec(sql)
	if err != nil {
		log.Error("DBHandler->SortProducts err:%s", err.Error())
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetAllFtOrder(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {

	engine := handle.mysqlEngine
	start := (pageIndex - 1) * pageSize
	limit := pageSize

	results, err := engine.Query(datastruct.AllOrderSqlForFt, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetAllFtOrder err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	arr := make([]interface{}, 0, len(results))
	var rs interface{}
	for _, v := range results {
		dataType := tools.StringToInt(string(v["type"][:]))
		createdAt := tools.StringToInt64(string(v["created_at"][:]))
		orderId := tools.StringToInt64(string(v["orderid"][:]))
		productName := string(v["product_name"][:])
		productDesc := string(v["product_desc"][:])
		avatar := string(v["avatar"][:])
		nickName := string(v["nick_name"][:])
		ft_order := new(datastruct.RespOrderForFt)
		ft_order.Avatar = avatar
		ft_order.CreatedAt = createdAt
		ft_order.DataType = dataType
		ft_order.NickName = nickName
		ft_order.OrderId = orderId
		ft_order.ProductDesc = productDesc
		ft_order.ProductName = productName
		switch dataType {
		case 1:
			//购买成功,评价与未评价的订单
			ft_p_order := new(datastruct.RespPurchaseOrderForFt)
			ft_p_order.Order = ft_order
			ft_p_order.IsAppraised = tools.StringToBool(string(v["handle"][:]))
			rs = ft_p_order
		case 2:
			//退款中
			ft_rg_order := new(datastruct.RespRefundingOrderForFt)
			ft_rg_order.Order = ft_order
			rs = ft_rg_order
		case 3:
			//退款结果
			ft_rf_order := new(datastruct.RespRefundFinishedOrderForFt)
			ft_rf_order.Order = ft_order
			ft_rf_order.RefundType = datastruct.UserOrderRefundType(tools.StringToInt(string(v["handle"][:])))
			rs = ft_rf_order
		case 4:
			//维权中
			ft_rtg_order := new(datastruct.RespRightingOrderForFt)
			ft_rtg_order.Order = ft_order
			rs = ft_rtg_order
		case 5:
			//维权结果
			ft_rtf_order := new(datastruct.RespRightFinishedOrderForFt)
			ft_rtf_order.Order = ft_order
			ft_rtf_order.IsAgree = tools.StringToBool(string(v["handle"][:]))
			rs = ft_rtf_order
		}
		arr = append(arr, rs)
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) CreateFakeAppraised(body *datastruct.FakeAppraisedBody, ft_id int) datastruct.CodeType {
	engine := handle.mysqlEngine
	sql := "select account from hot_f_t_info where f_t_id = ?"
	results, err := engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->CreateFakeAppraised err0")
		return datastruct.UpdateDataFailed
	}
	account := tools.StringToInt64(string(results[0]["account"][:]))
	var pay int64
	pay = 2
	if account < pay {
		return datastruct.AccountLess
	}

	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	ap_info := new(datastruct.AppraisedInfo)
	ap_info.AppraisedType = datastruct.Char
	ap_info.CreatedAt = body.Time
	ap_info.Desc = body.Desc
	ap_info.ProductId = body.Id
	ap_info.IsAnonym = true
	ap_info.IsFake = true
	ap_info.Mark = body.Mark
	ap_info.Score = body.Score

	_, err = session.InsertOne(ap_info)
	if err != nil {
		str := fmt.Sprintf("DBHandler->CreateFakeAppraised Insert AppraisedInfo err:%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	sql = "update hot_f_t_info set account = account - ? where f_t_id = ?"
	_, err = session.Exec(sql, pay, ft_id)
	if err != nil {
		str := fmt.Sprintf("DBHandler->CreateFakeAppraised update account :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->CreateFakeAppraised Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	return datastruct.NULLError
}

func (handle *DBHandler) IsAgreeRefund(body *datastruct.IsAgreeRefundBody) datastruct.CodeType {
	nowTime := time.Now().Unix()
	engine := handle.mysqlEngine
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	rs, err := session.Where("id=?", body.Id).Delete(new(datastruct.UserOrderRefunding))
	if err != nil || rs <= 0 {
		str := fmt.Sprintf("DBHandler->IsAgreeRefund delete err:%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	var tmp interface{}
	if body.IsAgree {
		rf_order := new(datastruct.UserOrderRefundFinished)
		rf_order.CreatedAt = nowTime
		rf_order.Id = body.Id
		rf_order.RefundType = datastruct.Apply
		tmp = rf_order
	} else {
		rtg_order := new(datastruct.UserOrderRighting)
		rtg_order.CreatedAt = nowTime
		rtg_order.Id = body.Id
		tmp = rtg_order
	}

	_, err = session.InsertOne(tmp)
	if err != nil {
		str := fmt.Sprintf("DBHandler->IsAgreeRefund insert err:%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}

	if body.IsAgree {
		sql := "select pro.price,uoi.user_id from user_order_info uoi join product_info pro on pro.id = uoi.product_id where uoi.id = ?"
		results, err := engine.Query(sql, body.Id)
		if err != nil || len(results) <= 0 {
			rollbackError("DBHandler->IsAgreeRefund get product price err", session)
			return datastruct.UpdateDataFailed
		}
		price := tools.StringToFloat64(string(results[0]["price"][:]))
		user_id := tools.StringToInt64(string(results[0]["user_id"][:]))
		gold_count := int64(price)
		sql = "update hot_user_info set gold_count = gold_count + ? where user_id = ?"
		_, err = session.Exec(sql, gold_count, user_id)
		if err != nil {
			str := fmt.Sprintf("DBHandler->IsAgreeRefund update user account err:%s", err.Error())
			rollbackError(str, session)
			return datastruct.UpdateDataFailed
		}
		g_change := new(datastruct.UserGoldChange)
		g_change.ChangeType = datastruct.Refund
		g_change.CreatedAt = nowTime
		g_change.UserId = user_id
		g_change.VarGold = gold_count
		_, err = session.InsertOne(g_change)
		if err != nil {
			str := fmt.Sprintf("DBHandler->IsAgreeRefund insert UserGoldChange err:%s", err.Error())
			rollbackError(str, session)
			return datastruct.UpdateDataFailed
		}
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("DBHandler->IsAgreeRefund Commit :%s", err.Error())
		rollbackError(str, session)
		return datastruct.UpdateDataFailed
	}
	return datastruct.NULLError
}

func (handle *DBHandler) GetFinance(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return nil, code
	}
	sql := "select * from (select sum(pro.price) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = 0 union all select sum(pro.price) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = 1) as tmp_sum"

	results, err := engine.Query(sql, shop_id, shop_id)

	if err != nil {
		log.Error("DBHandler->getFinance err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	notCheckAmount := tools.StringToFloat64(string(results[0]["amount"][:])) //未结算订单总金额
	checkedAmount := tools.StringToFloat64(string(results[1]["amount"][:]))  //已结算订单总金额
	totalAmount := notCheckAmount + checkedAmount                            //总订单金额

	sql = "select balance_total,balance,income_per from hot_f_t_info where f_t_id = ?"
	results, err = engine.Query(sql, ft_id)
	checkedBalance := tools.StringToFloat64(string(results[0]["balance_total"][:])) //已结算的收益
	balance := tools.StringToFloat64(string(results[0]["balance"][:]))              //可提现的金额
	IncomePer := tools.StringToFloat64(string(results[0]["income_per"][:]))         //收益提成百分比

	notCheckBalance := notCheckAmount * IncomePer / 100.0 //未结算的收益
	totalBalance := notCheckBalance + checkedBalance      //总收益
	rs := new(datastruct.RespFtFinance)
	rs.Balance = balance
	rs.CheckedAmount = checkedAmount
	rs.CheckedBalance = checkedBalance
	rs.NotCheckAmount = notCheckAmount
	rs.NotCheckBalance = notCheckBalance
	rs.TotalAmount = totalAmount
	rs.TotalBalance = totalBalance
	return rs, datastruct.NULLError
}

func (handle *DBHandler) GetProducts(ft_id int) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return nil, code
	}
	sql := "select id,product_name from product_info where shop_id = ? order by updated_at desc"
	results, err := engine.Query(sql, shop_id)
	if err != nil {
		log.Error("DBHandler->GetProducts err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	arr := make([]*datastruct.RespProducts, 0, len(results))
	for _, v := range results {
		pro := new(datastruct.RespProducts)
		pro.Id = tools.StringToInt(string(v["id"][:]))
		pro.ProductName = string(v["product_name"][:])
		arr = append(arr, pro)
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) GetAmountList(datatype int, pageIndex int, pageSize int, ft_id int) (interface{}, datastruct.CodeType) {
	start := (pageIndex - 1) * pageSize
	limit := pageSize

	engine := handle.mysqlEngine
	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return nil, code
	}
	var sql string
	var count float64
	var totalAmount float64
	var err error
	var results []map[string][]byte
	switch datatype {
	case 0:
		fallthrough
	case 1:
		sql_count := "select * from (select sum(pro.price) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = ? union all select count(uoi.id) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = ?) as tmp_sum"
		results, err = engine.Query(sql_count, shop_id, datatype, shop_id, datatype)
		if err != nil {
			log.Error("DBHandler->GetAmountList get amount err: %s", err.Error())
			return nil, datastruct.GetDataFailed
		}
		totalAmount = tools.StringToFloat64(string(results[0]["amount"][:]))
		count = tools.StringToFloat64(string(results[1]["amount"][:]))
		sql = "select uoc.is_checked as checked,cui.nick_name,cui.avatar,pro.product_name,pro.price,uoc.updated_at as createdat from user_order_info uoi join user_order_check uoc on uoc.id = uoi.id join product_info pro on pro.id=uoi.product_id join cold_user_info cui on cui.id=uoi.user_id where uoc.is_checked = ? and pro.shop_id = ? ORDER BY uoc.updated_at DESC LIMIT ?,?"
		results, err = engine.Query(sql, datatype, shop_id, start, limit)
	case 2:
		sql_count := "select * from (select sum(pro.price) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = 0 union all select sum(pro.price) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = 1 union all select count(uoi.id) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ?) as tmp_sum"
		results, err = engine.Query(sql_count, shop_id, shop_id, shop_id)
		if err != nil {
			log.Error("DBHandler->GetAmountList get amount err: %s", err.Error())
			return nil, datastruct.GetDataFailed
		}
		notCheckAmount := tools.StringToFloat64(string(results[0]["amount"][:])) //未结算订单总金额
		checkedAmount := tools.StringToFloat64(string(results[1]["amount"][:]))  //已结算订单总金额
		count = tools.StringToFloat64(string(results[2]["amount"][:]))           //总数量
		totalAmount = notCheckAmount + checkedAmount                             //总订单金额
		sql = "select tor.checked,cui.nick_name,cui.avatar,pro.product_name,pro.price,tor.createdat from (select 0 as checked,uoi.user_id as uid,uoi.product_id as pid,uoc.updated_at as createdat from user_order_info uoi join user_order_check uoc on uoc.id = uoi.id where uoc.is_checked = 0 union all select 1 as checked,uoi.user_id as uid,uoi.product_id as pid,uoc.updated_at as createdat from user_order_info uoi join user_order_check uoc on uoc.id = uoi.id where uoc.is_checked = 1 ) as tor join product_info pro on pro.id=tor.pid join cold_user_info cui on cui.id=tor.uid where pro.shop_id = ? ORDER BY tor.createdat DESC LIMIT ?,?"
		results, err = engine.Query(sql, shop_id, start, limit)
	}

	rs := new(datastruct.RespOrderList)
	rs.Count = int(count)
	rs.TotalAmount = totalAmount

	if err != nil {
		log.Error("DBHandler->GetAmountList get list err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}

	arr := make([]*datastruct.RespOrderInfo, 0, len(results))
	for _, v := range results {
		rs := new(datastruct.RespOrderInfo)
		rs.Avatar = string(v["avatar"][:])
		rs.NickName = string(v["nick_name"][:])
		rs.CreatedAt = tools.StringToInt64(string(v["createdat"][:]))
		rs.IsChecked = tools.StringToBool(string(v["checked"][:]))
		rs.Price = tools.StringToFloat64(string(v["price"][:]))
		rs.ProductName = string(v["product_name"][:])
		arr = append(arr, rs)
	}
	rs.List = arr
	return rs, datastruct.NULLError
}

func (handle *DBHandler) GetIncomeList(datatype int, pageIndex int, pageSize int, ft_id int) (interface{}, datastruct.CodeType) {
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	engine := handle.mysqlEngine
	shop_id, code := getShopId(engine, ft_id)
	if code != datastruct.NULLError {
		return nil, code
	}
	var sql string
	var count float64
	var totalAmount float64
	var err error
	var results []map[string][]byte
	switch datatype {
	case 0:
		var income_per float64
		totalAmount, count, code, income_per = computeAmount(datatype, shop_id, ft_id, engine)
		if code != datastruct.NULLError {
			return nil, code
		}
		sql = "select uoc.is_checked as checked,cui.nick_name,cui.avatar,pro.product_name,pro.price * %f as price,uoc.updated_at as createdat from user_order_info uoi join user_order_check uoc on uoc.id = uoi.id join product_info pro on pro.id=uoi.product_id join cold_user_info cui on cui.id=uoi.user_id where uoc.is_checked = 1 and pro.shop_id = ? ORDER BY uoc.updated_at DESC LIMIT ?,?"
		sql = fmt.Sprintf(sql, income_per/100.0)
		results, err = engine.Query(sql, shop_id, start, limit)
	case 1:
		totalAmount, count, code, _ = computeAmount(datatype, shop_id, ft_id, engine)
		if code != datastruct.NULLError {
			return nil, code
		}
		sql = "select uoc.is_checked as checked,cui.nick_name,cui.avatar,pro.product_name,uoc.checked_income as price,uoc.updated_at as createdat from user_order_info uoi join user_order_check uoc on uoc.id = uoi.id join product_info pro on pro.id=uoi.product_id join cold_user_info cui on cui.id=uoi.user_id where uoc.is_checked = 1 and pro.shop_id = ? ORDER BY uoc.updated_at DESC LIMIT ?,?"
		results, err = engine.Query(sql, shop_id, start, limit)
	case 2:
		notCheckedAmount, notCheckedCount, code, income_per := computeAmount(0, shop_id, ft_id, engine)
		if code != datastruct.NULLError {
			return nil, code
		}
		checkedAmount, checkedCount, code, _ := computeAmount(1, shop_id, ft_id, engine)
		if code != datastruct.NULLError {
			return nil, code
		}

		totalAmount = notCheckedAmount + checkedAmount
		count = notCheckedCount + checkedCount
		sql = "select tor.checked,cui.nick_name,cui.avatar,tor.product_name,tor.price,tor.createdat from (select 0 as checked,uoi.user_id as uid,pro.product_name,pro.price * %f as price,uoc.updated_at as createdat,pro.shop_id from user_order_info uoi join user_order_check uoc on uoc.id = uoi.id join product_info pro on pro.id=uoi.product_id where uoc.is_checked = 0 union all select 1 as checked,uoi.user_id as uid,pro.product_name,uoc.checked_income as price,uoc.updated_at as createdat,pro.shop_id from user_order_info uoi join user_order_check uoc on uoc.id = uoi.id join product_info pro on pro.id=uoi.product_id where uoc.is_checked = 1 ) as tor join cold_user_info cui on cui.id=tor.uid where tor.shop_id = ? ORDER BY tor.createdat DESC LIMIT ?,?"
		sql = fmt.Sprintf(sql, income_per/100.0)
		results, err = engine.Query(sql, shop_id, start, limit)
	}

	rs := new(datastruct.RespOrderList)
	rs.Count = int(count)
	rs.TotalAmount = totalAmount

	if err != nil {
		log.Error("DBHandler->GetAmountList get list err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}

	arr := make([]*datastruct.RespOrderInfo, 0, len(results))
	for _, v := range results {
		rs := new(datastruct.RespOrderInfo)
		rs.Avatar = string(v["avatar"][:])
		rs.NickName = string(v["nick_name"][:])
		rs.CreatedAt = tools.StringToInt64(string(v["createdat"][:]))
		rs.IsChecked = tools.StringToBool(string(v["checked"][:]))
		rs.Price = tools.StringToFloat64(string(v["price"][:]))
		rs.ProductName = string(v["product_name"][:])
		arr = append(arr, rs)
	}
	rs.List = arr
	return rs, datastruct.NULLError
}

func computeAmount(datatype int, shop_id int, ft_id int, engine *xorm.Engine) (float64, float64, datastruct.CodeType, float64) {
	var count float64
	var totalAmount float64
	var income_per float64
	if datatype == 0 {
		sql := "select * from (select sum(pro.price) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = 0 union all select count(uoi.id) as amount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = 0) as tmp_sum"
		results, err := engine.Query(sql, shop_id, shop_id)
		if err != nil {
			log.Error("DBHandler->computeAmount get amount err: %s", err.Error())
			return -1, -1, datastruct.GetDataFailed, -1
		}
		totalPrice := tools.StringToFloat64(string(results[0]["amount"][:]))
		count = tools.StringToFloat64(string(results[1]["amount"][:]))
		sql = "select income_per from hot_f_t_info where f_t_id = ?"
		results, err = engine.Query(sql, ft_id)
		if err != nil {
			log.Error("DBHandler->computeAmount get hot_f_t_info err: %s", err.Error())
			return -1, -1, datastruct.GetDataFailed, -1
		}
		income_per = tools.StringToFloat64(string(results[0]["income_per"][:]))
		totalAmount = totalPrice * income_per / 100.0
	} else {
		sql := "select balance_total from hot_f_t_info where f_t_id = ?"
		results, err := engine.Query(sql, ft_id)
		if err != nil {
			log.Error("DBHandler->computeAmount get hot_f_t_info err: %s", err.Error())
			return -1, -1, datastruct.GetDataFailed, -1
		}
		totalAmount = tools.StringToFloat64(string(results[0]["balance_total"][:]))
		sql = "select count(uoi.id) as checkedcount from user_order_info uoi join product_info pro on pro.id=uoi.product_id join user_order_check uoc on uoc.id = uoi.id where pro.shop_id = ? and uoc.is_checked = 1"
		results, err = engine.Query(sql, shop_id)
		if err != nil {
			log.Error("DBHandler->computeAmount get count err: %s", err.Error())
			return -1, -1, datastruct.GetDataFailed, -1
		}
		count = tools.StringToFloat64(string(results[0]["checkedcount"][:]))
	}
	return totalAmount, count, datastruct.NULLError, income_per
}

func (handle *DBHandler) GetDrawCashParams(params_type datastruct.DrawCashParamsType) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	params := new(datastruct.DrawCashParams)
	has, err := engine.Where("params_type=?", params_type).Get(params)
	if err != nil || !has {
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.RespDrawCashParams)
	resp.MaxDrawCount = params.MaxDrawCount
	resp.MinCharge = params.MinCharge
	resp.MinPoundage = params.MinPoundage
	resp.PoundagePer = params.PoundagePer
	resp.RequireVerify = params.RequireVerify
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetFtDrawCashInfo(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select origin,poundage,charge,state,arrival_type,created_at from f_t_draw_cash_info where f_t_id = ? ORDER BY created_at DESC LIMIT ?,?"
	engine := handle.mysqlEngine
	results, err := engine.Query(sql, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetFtDrawCashInfo err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	arr := make([]*datastruct.RespFtDrawInfo, 0, len(results))
	for _, v := range results {
		ft_draw := new(datastruct.RespFtDrawInfo)
		ft_draw.Charge = tools.StringToFloat64(string(v["charge"][:]))
		ft_draw.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		ft_draw.Origin = tools.StringToFloat64(string(v["origin"][:]))
		ft_draw.Poundage = tools.StringToFloat64(string(v["poundage"][:]))
		state := tools.StringToInt(string(v["state"][:]))
		arrivalType := tools.StringToInt(string(v["arrival_type"][:]))
		ft_draw.State = tools.DrawCashStateToString(datastruct.DrawCashState(state))
		ft_draw.ArrivalType = tools.DrawCashArrivalTypeToString(datastruct.DrawCashArrivalType(arrivalType))
		arr = append(arr, ft_draw)
	}
	return arr, datastruct.NULLError
}

func (handle *DBHandler) GetFtAccountChangeInfo(ft_id int, pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select created_at,var_account,change_type from f_t_account_change where f_t_id = ? ORDER BY created_at DESC LIMIT ?,?"
	engine := handle.mysqlEngine
	results, err := engine.Query(sql, ft_id, start, limit)
	if err != nil {
		log.Error("DBHandler->GetFtAccountChangeInfo err: %s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.RespFtAccountChangeInfo)
	arr := make([]*datastruct.RespFtAccountChange, 0, len(results))
	for _, v := range results {
		tmp := new(datastruct.RespFtAccountChange)
		tmp.ChangeCount = tools.StringToInt64(string(v["var_account"][:]))
		tmp.CreatedAt = tools.StringToInt64(string(v["created_at"][:]))
		changeType := tools.StringToInt(string(v["change_type"][:]))
		tmp.ChangeType = tools.ScoreChangeTypeToString(datastruct.ScoreChangeType(changeType))
		arr = append(arr, tmp)
	}
	sql = "select account from hot_f_t_info where f_t_id=?"
	results, err = engine.Query(sql, ft_id)
	if err != nil || len(results) <= 0 {
		log.Error("DBHandler->GetFtAccountChangeInfo get hot_f_t_info err")
		return nil, datastruct.GetDataFailed
	}
	resp.Account = tools.StringToInt64(string(results[0]["account"][:]))
	resp.List = arr
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UserRegister(body *datastruct.UserRegisterBody) (*datastruct.RespUserLogin, datastruct.CodeType) {
	return userRegister(handle.mysqlEngine, body, nil)
}

func (handle *DBHandler) UserRegisterWithDetail(body *datastruct.UserRegisterDetailBody) (*datastruct.RespUserLogin, datastruct.CodeType) {
	return userRegister(handle.mysqlEngine, nil, body)
}

func userRegister(engine *xorm.Engine, body *datastruct.UserRegisterBody, detail *datastruct.UserRegisterDetailBody) (*datastruct.RespUserLogin, datastruct.CodeType) {
	nowTime := time.Now().Unix()
	cold_user := new(datastruct.ColdUserInfo)
	cold_user.CreatedAt = nowTime
	if body != nil {
		cold_user.Phone = body.Phone
		cold_user.Pwd = body.Pwd
		cold_user.Registration = body.Platform
		cold_user.Sex = datastruct.Male
		cold_user.Avatar = datastruct.DefaultUserAvatar
	} else {
		cold_user.Phone = detail.Phone
		cold_user.Pwd = detail.Pwd
		cold_user.Registration = detail.Platform
		cold_user.NickName = detail.NickName
		cold_user.Avatar = detail.Avatar
		cold_user.Sex = detail.Sex
		cold_user.ActualName = detail.ActualName
		cold_user.DateBirth = detail.Birthday
		cold_user.BirthPlace = detail.City
	}
	session := engine.NewSession()
	defer session.Close()
	session.Begin()

	_, err := session.InsertOne(cold_user)
	if err != nil {
		log.Error("userRegister err:%s", err.Error())
		return nil, datastruct.UpdateDataFailed
	}
	if body != nil {
		cold_user.NickName = fmt.Sprintf("用户%d", cold_user.Id)
		_, err = session.Where("id=?", cold_user.Id).Cols("nick_name").Update(cold_user)
		if err != nil {
			log.Error("userRegister update err:%s", err.Error())
			return nil, datastruct.UpdateDataFailed
		}
	}

	im_id := fmt.Sprintf("user_%d", cold_user.Id)
	im_privatekey, code := tools.AccountGenForIM(im_id, important.IM_SDK_APPID)
	if code != datastruct.NULLError {
		return nil, code
	}

	token := commondata.UniqueId()
	hot_user := new(datastruct.HotUserInfo)
	hot_user.UserId = cold_user.Id
	hot_user.Token = token
	hot_user.LoginTime = nowTime
	hot_user.IMPrivateKey = im_privatekey

	_, err = session.InsertOne(hot_user)
	if err != nil {
		log.Error("userRegister err:%s", err.Error())
		return nil, datastruct.UpdateDataFailed
	}

	err = session.Commit()
	if err != nil {
		str := fmt.Sprintf("userRegister Commit :%s", err.Error())
		rollbackError(str, session)
		return nil, datastruct.UpdateDataFailed
	}
	resp := new(datastruct.RespUserLogin)
	resp.AccountState = hot_user.AccountState
	resp.IMPrivateKey = im_privatekey
	resp.Id = cold_user.Id
	resp.Token = hot_user.Token
	return resp, datastruct.NULLError
}

func (handle *DBHandler) UserLoginWithPwd(body *datastruct.UserLoginWithPwdBody) (*datastruct.RespUserLogin, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select cui.id,hui.account_state from cold_user_info cui join hot_user_info hui on hui.user_id=cui.id where phone=? and pwd=?"
	results, err := engine.Query(sql, body.Phone, body.Pwd)
	if err != nil {
		log.Error("UserLoginWithPwd query err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	if len(results) <= 0 {
		return nil, datastruct.LoginFailed
	}
	userId := tools.StringToInt64(string(results[0]["id"][:]))
	accountState := tools.StringToAccountState(string(results[0]["account_state"][:]))
	im_id := fmt.Sprintf("user_%d", userId)
	im_privatekey, code := tools.AccountGenForIM(im_id, important.IM_SDK_APPID)
	if code != datastruct.NULLError {
		return nil, code
	}

	hot_user := new(datastruct.HotUserInfo)
	hot_user.LoginTime = time.Now().Unix()
	hot_user.IMPrivateKey = im_privatekey
	hot_user.Token = commondata.UniqueId()
	_, err = engine.Where("user_id=?", userId).Cols("token", "login_time", "i_m_private_key").Update(hot_user)
	if err != nil {
		log.Error("UserLoginWithPwd update err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.RespUserLogin)
	resp.IMPrivateKey = im_privatekey
	resp.Id = userId
	resp.Token = hot_user.Token
	resp.AccountState = accountState
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetHomeData(platform datastruct.Platform) (interface{}, datastruct.CodeType) {
	engine := handle.mysqlEngine
	sql := "select img_url,is_jump_to,jump_to from ad_info where is_hidden = 0 and (platform = ? or platform = ?) order by sort_id desc"
	results, err := engine.Query(sql, platform, datastruct.PC+1)
	if err != nil {
		log.Error("GetHomeData query ad err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	resp := new(datastruct.RespHomeData)
	ad_list := make([]*datastruct.RespAdInfo, 0, len(results))
	for _, v := range results {
		ad := new(datastruct.RespAdInfo)
		ad.ImgUrl = string(v["img_url"][:])
		isJumpTo := tools.StringToBool(string(v["is_jump_to"][:]))
		ad.IsJumpTo = isJumpTo
		if isJumpTo {
			ad.JumpTo = string(v["jump_to"][:])
		}
		ad_list = append(ad_list, ad)
	}

	sql = "select cft.id,cft.nick_name,cft.avatar from cold_f_t_info cft join shop_info shi on shi.f_t_id=cft.id where shi.is_recommended = 1"
	results, err = engine.Query(sql)
	if err != nil {
		log.Error("GetHomeData query ft err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	ft_list := make([]*datastruct.RespFtData, 0, len(results))
	for _, v := range results {
		ft_data := new(datastruct.RespFtData)
		ft_data.Avatar = string(v["avatar"][:])
		ft_data.NickName = string(v["nick_name"][:])
		ft_data.Id = tools.StringToInt(string(v["id"][:]))
		ft_list = append(ft_list, ft_data)
	}

	sql = "select img_nav.img_url,img_nav.key from img_nav where is_hidden = 0 order by sort_id desc"
	results, err = engine.Query(sql)
	if err != nil {
		log.Error("GetHomeData query bottom err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	banner_list := make([]*datastruct.RespBottomBanner, 0, len(results))
	for _, v := range results {
		banner := new(datastruct.RespBottomBanner)
		banner.ImgUrl = string(v["img_url"][:])
		banner.Key = string(v["key"][:])
		banner_list = append(banner_list, banner)
	}
	resp.AdInfo = ad_list
	resp.FtCount = datastruct.TmpFtCount
	resp.SolveCount = datastruct.TmpSolveCount
	resp.Commend = ft_list
	resp.BottomBanner = banner_list
	return resp, datastruct.NULLError
}

func (handle *DBHandler) GetHomeAppraise(pageIndex int, pageSize int) (interface{}, datastruct.CodeType) {
	start := (pageIndex - 1) * pageSize
	limit := pageSize
	sql := "select ari.desc,cui.nick_name,cui.avatar from appraised_info ari join cold_user_info cui on ari.user_id=cui.id where is_passed = 1 and (appraised_type = ? or appraised_type = ?) order by ari.created_at desc LIMIT ?,?"
	engine := handle.mysqlEngine
	results, err := engine.Query(sql, datastruct.Char, datastruct.CharAndImg, start, limit)
	if err != nil {
		log.Error("GetHomeAppraised err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	rs := make([]*datastruct.RespHomeAppraise, 0, len(results))
	for _, v := range results {
		tmp := new(datastruct.RespHomeAppraise)
		tmp.Avatar = string(v["avatar"][:])
		tmp.NickName = string(v["nick_name"][:])
		tmp.Desc = string(v["desc"][:])
		rs = append(rs, tmp)
	}
	return rs, datastruct.NULLError
}

func (handle *DBHandler) QueryFtIds(key string) ([]int, datastruct.CodeType) {
	sql := "select cft.id from cold_f_t_info cft join hot_f_t_info hft on cft.id = hft.f_t_id"
	where_str := " where cft.nick_name like "
	like_str := "'%" + key + "%'"
	or_str := " or hft.mark like "
	orderby := " ORDER BY hft.sort_id asc,hft.order_score desc"
	if key != "" {
		sql = sql + where_str + like_str + or_str + like_str + orderby
	} else {
		sql = sql + orderby
	}
	engine := handle.mysqlEngine
	results, err := engine.Query(sql)
	if err != nil {
		log.Error("QueryFtIds err:%s", err.Error())
		return nil, datastruct.GetDataFailed
	}
	ids := make([]int, 0, len(results))
	for _, v := range results {
		ft_id := tools.StringToInt(string(v["id"][:]))
		ids = append(ids, ft_id)
	}
	return ids, datastruct.NULLError
}

func (handle *DBHandler) QueryFtInfo(onlineIds []int, offlineIds []int) (interface{}, datastruct.CodeType) {
	online_length := len(onlineIds)
	offline_length := len(offlineIds)
	max := 50
	//engine := handle.mysqlEngine
	if online_length > 0 {
		n := online_length / max
		if n == 0 {
			//onlineIds[:]
		} else {
			var i int
			for i = 1; i <= n; i++ {
				//onlineIds[(i-1)*max : max*i]
			}
			if online_length%max != 0 {
				//onlineIds[(i-1)*max:]
			}
		}
	}
	if offline_length > 0 {
		n := offline_length / max
		if n == 0 {
			//offlineIds[:]
		} else {
			var i int
			for i = 1; i <= n; i++ {
				//offlineIds[(i-1)*max : max*i]
			}
			if offline_length%max != 0 {
				//offlineIds[(i-1)*max:]
			}
		}
	}
	return nil, datastruct.NULLError
}

func queryFtInfo(ids []int, engine *xorm.Engine) (interface{}, datastruct.CodeType) {

	return nil, datastruct.NULLError
}

//string(results[0][column_name][:])
