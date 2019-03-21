package datastruct

import "sync"

type ServerData struct {
	RWMutex    *sync.RWMutex //读写互斥量
	Version    string        //当前服务端版本号
	IsMaintain bool          //是否维护
}

const DefaultId = 1
const MAX_PRODUCT_COUNT = 50
const Reidis_IdField = "Id"
const Reidis_AccountStateField = "AccountState"
const DefaultUserAvatar = "https://shensuanzi.oss-cn-shenzhen.aliyuncs.com/user_avatar/defaultavatar.png"
const RedisExpireTime = 300 //Redis过期时间300秒

type CodeType int //错误码
const (
	NULLError                   CodeType = iota //无错误
	ParamError                                  //参数错误,数据为空或者类型不对等
	LoginFailed                                 //登录失败,如无此账号或者密码错误等
	JsonParseFailedFromPostBody                 //来自post请求中的Body解析json失败
	GetDataFailed                               //获取数据失败
	UpdateDataFailed                            //修改数据失败
	TokenError                                  //没有Token或者值为空,或者不存在此Token
	JsonParseFailedFromPutBody                  //来自put请求中的Body解析json失败
	WXCodeInvalid                               //无效的微信code
	PlatformInvalid                             //无效的平台参数
	NickNamePhoneIsUsed                         //昵称或手机号已被使用
	NotRegisterPhone                            //手机号未注册
	PwdError                                    //密码错误
	AuthingCode                                 //账号审核中
	AuthFailedCode                              //账号审核失败
	Redirect                                    //重定向
	Sensitive                                   //敏感词
	MaxCreateCount                              //最大创建数量
	AccountLess                                 //余额不足
	HeaderParamError                            //header参数不足或参数错误
)

type Platform int //平台
const (
	APP Platform = iota
	H5
	PC
)

type PayPlatform int //付费平台
const (
	WXPay PayPlatform = iota
)

type Sex int

const (
	Female Sex = iota
	Male
	Secret //保密
)

type Calendar int

const (
	GongLi Calendar = iota //公历
	Nongli                 //农历
)

type AccountState int

const (
	Normal    AccountState = iota //正常
	Freeze                        //冻结
	BlackList                     //黑名单
)

type AppraisedType int

const (
	Char       AppraisedType = iota //只有文字
	Img                             //只有图片
	CharAndImg                      //图文
)

type GoldChangeType int

const (
	Deposit GoldChangeType = iota //充值
	Cost                          //消费
	Refund                        //退款
)

type DrawCashState int

const (
	Review  DrawCashState = iota //审核中
	Succeed                      //提现成功
	Failed                       //提现失败
)

type DrawCashArrivalType int //到账类型
const (
	ArrivalWX  DrawCashArrivalType = iota //微信钱包
	ArrivalZFB                            //支付宝
)

type DrawCashParamsType int //提现参数类型
const (
	User DrawCashParamsType = iota //用户
	FT                             //命理师
)

type UserOrderRefundType int //退款类型
const (
	Apply UserOrderRefundType = iota //用户申请退款
	Auto                             //系统自动退款
)

type OrderRightsType int //维权订单类型
const (
	NotAgree OrderRightsType = iota //命理师不同意退款
)

type UserOrderRefundResultType int //退款结果类型
const (
	AgreeResult    UserOrderRefundResultType = iota //命理师同意退款
	NotAgreeResult                                  //命理师不同意退款
	AutoResult                                      //系统自动退款
)

type AuthState int //账号审核状态
const (
	Authing     AuthState = iota //审核中
	AuthFailed                   //审核失败
	AuthSucceed                  //审核成功
)

type ScoreChangeType int //积分变化类型

const (
	DepositScore ScoreChangeType = iota //充值
	CostScore                           //消费
)

type HttpStatusCode int //错误码
const (
	Maintenance  CodeType = 900 //服务器维护中
	VersionError          = 901 //客户端与服务器版本不一致
)

// type FTOnlineState int //命理师在线状态
// const (
// 	AutoAdjust    FTOnlineState = iota //根据在线时间自动调整状态
// 	ManualOffline                      //手动设置离线
// 	ManualBusy                         //手动设置忙碌
// )

type IdCardState int //审核身份证状态
const (
	IdCardNotSubmit IdCardState = iota //未提交
	IdCardSubmited                     //已提交
	IdCardNotPass                      //未通过
	IdCardPassed                       //已通过
)

type OnlineUIState int //在线UI状态
const (
	OnlineUI  OnlineUIState = iota //在线
	BusyUI                         //忙碌
	OfflineUI                      //离线
)

type FtRedisData struct {
	FtId         int
	Token        string
	AccountState AccountState
}

type UserRedisData struct {
	UserId       int64
	Token        string
	AccountState AccountState
}
