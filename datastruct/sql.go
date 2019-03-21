package datastruct

const AllOrderSqlForFt = "select tmp.type,tmp.orderid,tmp.created_at,pro.product_name,pro.product_desc,cui.avatar,cui.nick_name,tmp.handle from (SELECT 1 as type,uoi.id as orderid,uoc.updated_at as created_at,uoi.product_id,uoi.user_id,uoc.is_appraised as handle from user_order_info uoi join user_order_check uoc on uoi.id = uoc.id union all SELECT 2 as type,uoi.id as orderid,uorg.created_at,uoi.product_id,uoi.user_id,-1 as handle from user_order_info uoi join user_order_refunding uorg on uoi.id = uorg.id union all SELECT 3 as type,uoi.id as orderid,uorf.created_at,uoi.product_id,uoi.user_id,uorf.refund_type as handle from user_order_info uoi join user_order_refund_finished uorf on uoi.id = uorf.id union all SELECT 4 as type,uoi.id as orderid,uortg.created_at,uoi.product_id,uoi.user_id,-1 as handle from user_order_info uoi join user_order_righting uortg on uoi.id = uortg.id union all SELECT 5 as type,uoi.id as orderid,uortf.created_at,uoi.product_id,uoi.user_id,uortf.is_agree as handle from user_order_info uoi join user_order_rights_finished uortf on uoi.id = uortf.id) as tmp join product_info pro on pro.id = tmp.product_id join shop_info si on si.id = pro.shop_id join cold_user_info cui on cui.id = tmp.user_id where si.f_t_id = ? ORDER BY tmp.created_at DESC LIMIT ?,?"

const AllMsgSqlForFt = "select * from (select -1 as id,ur.user_read,ur.created_at,ur.user_id,1 as type,-1 as f_t_nick_name,-1 as product_name,-1 as handle from user_register_msg ur union all select uor.id,uor.user_read,uor.created_at,uor.user_id,2 as type,f_t_nick_name,product_name,refund_result_type as handle from user_order_refund_msg uor union all select oi.id,oi.user_read,oi.created_at,oi.user_id,3 as type,f_t_nick_name,product_name,-1 as handle from order_info_msg oi union all select orf.id,orf.user_read,orf.created_at,orf.user_id,4 as type,f_t_nick_name,product_name,is_agree as handle from order_rights_finished_msg orf) as tmp_msg where user_id = ? ORDER BY created_at DESC LIMIT ?,?"

/*用户信息冷数据*/
type ColdUserInfo struct {
	Id           int64    `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	Phone        string   `xorm:"VARCHAR(20) not null unique COMMENT('手机号')"`
	Pwd          string   `xorm:"VARCHAR(30) null COMMENT('密码')"`
	NickName     string   `xorm:"VARCHAR(100) null COMMENT('昵称')"`
	Avatar       string   `xorm:"VARCHAR(500) null COMMENT('头像')"`
	Sex          Sex      `xorm:"TINYINT(1) null COMMENT('性别')"`
	ActualName   string   `xorm:"VARCHAR(100) null COMMENT('真实姓名')"`
	DateBirth    string   `xorm:"VARCHAR(16) null COMMENT('出生日期')"`
	BirthPlace   string   `xorm:"VARCHAR(100) null COMMENT('出生地')"`
	Calendar     Calendar `xorm:"TINYINT(1) null COMMENT('0是公历,1是农历')"`
	Registration Platform `xorm:"TINYINT(1) not null COMMENT('注册平台')"`
	CreatedAt    int64    `xorm:"bigint not null COMMENT('注册时间')"`
}

/*用户信息热数据*/
type HotUserInfo struct {
	UserId        int64        `xorm:"bigint not null pk COMMENT('用户Id')"`
	Token         string       `xorm:"CHAR(40) not null COMMENT('令牌值')"`
	LoginTime     int64        `xorm:"bigint not null COMMENT('登录时间')"`
	GoldCount     int64        `xorm:"bigint not null default 0 COMMENT('当前金币数')"`
	DepositTotal  float64      `xorm:"decimal(16,2) not null default 0 COMMENT('总充值额')"`
	Balance       float64      `xorm:"decimal(16,2) not null default 0 COMMENT('当前佣金额')"`
	BalanceTotal  float64      `xorm:"decimal(16,2) not null default 0 COMMENT('总佣金额')"`
	PayGoldTotal  float64      `xorm:"decimal(16,2) not null default 0 COMMENT('金币消费总额')"`
	PurchaseTotal float64      `xorm:"decimal(16,2) not null default 0 COMMENT('直接购买消费总额')"`
	AccountState  AccountState `xorm:"TINYINT(1) not null default 0 COMMENT('账户状态,0是正常,1是黑名单,2是冻结')"`
	IMPrivateKey  string       `xorm:"CHAR(255) not null COMMENT('IM私钥')"`
}

/*微信用户信息*/
type WXUserInfo struct {
	UserId           int64  `xorm:"bigint not null pk COMMENT('用户Id')"`
	UUID             string `xorm:"CHAR(50) not null COMMENT('微信UUID')"`
	PayOpenidForGzh  string `xorm:"CHAR(50) null COMMENT('微信公众号OpenId')"`
	PayOpenidForKfpt string `xorm:"CHAR(50) null COMMENT('微信开放平台OpenId')"`
	PayeeOpenid      string `xorm:"CHAR(50) null COMMENT('提现OpenId')"`
	GzhAppid         string `xorm:"CHAR(22) not null COMMENT('微信公众号Appid')"`
	KfptAppid        string `xorm:"CHAR(22) not null COMMENT('微信开放平台Appid')"`
}

/*微信命理师信息*/
type WXFTInfo struct {
	FTId             int    `xorm:"INT(11) not null pk COMMENT('大师Id')"`
	UUID             string `xorm:"CHAR(50) not null COMMENT('微信UUID')"`
	PayOpenidForGzh  string `xorm:"CHAR(50) null COMMENT('微信公众号OpenId')"`
	PayOpenidForKfpt string `xorm:"CHAR(50) null COMMENT('微信开放平台OpenId')"`
	PayeeOpenid      string `xorm:"CHAR(50) null COMMENT('提现OpenId')"`
	GzhAppid         string `xorm:"CHAR(22) not null COMMENT('微信公众号Appid')"`
	KfptAppid        string `xorm:"CHAR(22) not null COMMENT('微信开放平台Appid')"`
}

/*广告信息表*/
type AdInfo struct {
	Id       int      `xorm:"not null pk autoincr INT(4) COMMENT('自增编号')"`
	ImgUrl   string   `xorm:"VARCHAR(255) not null COMMENT('广告图片地址')"`
	IsJumpTo bool     `xorm:"TINYINT(1) not null COMMENT('是否跳转')"`
	JumpTo   string   `xorm:"VARCHAR(255) null COMMENT('跳转地址')"`
	SortId   int      `xorm:"not null INT(4) COMMENT('排序Id')"`
	Platform Platform `xorm:"TINYINT(1) not null COMMENT('展示平台')"`
	IsHidden bool     `xorm:"TINYINT(1) not null COMMENT('是否隐藏')"`
}

/*图片导航表*/
type ImgNav struct {
	Id       int    `xorm:"not null pk autoincr INT(4) COMMENT('自增编号')"`
	ImgUrl   string `xorm:"VARCHAR(255) not null COMMENT('图片地址')"`
	Key      string `xorm:"VARCHAR(255) not null COMMENT('查询关键字')"`
	SortId   int    `xorm:"not null INT(4) COMMENT('排序Id')"`
	IsHidden bool   `xorm:"TINYINT(1) not null COMMENT('是否隐藏')"`
}

/*用户评价表*/
type AppraisedInfo struct {
	Id            int           `xorm:"not null pk autoincr INT(11) COMMENT('自增编号')"`
	ProductId     int           `xorm:"not null INT(11) COMMENT('产品Id')"`
	UserId        int64         `xorm:"bigint not null default 0 COMMENT('用户Id')"`
	Score         float64       `xorm:"decimal(4,2) not null COMMENT('评分')"`
	Mark          string        `xorm:"VARCHAR(255) null COMMENT('多个标签内容')"`
	AppraisedType AppraisedType `xorm:"TINYINT(1) not null COMMENT('0为只有文字,1只有图,2图文')"`
	Desc          string        `xorm:"VARCHAR(255) null COMMENT('评价描述')"`
	IsAnonym      bool          `xorm:"TINYINT(1) not null COMMENT('是否匿名')"`
	IsFake        bool          `xorm:"TINYINT(1) not null COMMENT('是否为假数据')"`
	IsPassed      bool          `xorm:"TINYINT(1) not null default 1 COMMENT('是否审核通过')"`
	CreatedAt     int64         `xorm:"bigint not null COMMENT('评价时间')"`
}

/*评价标签信息*/
type AppraisedMark struct {
	Id   int    `xorm:"not null pk autoincr INT(6) COMMENT('自增编号')"`
	Mark string `xorm:"VARCHAR(50) not null COMMENT('标签内容')"`
}

/*评价图片*/
type AppraisedImgs struct {
	Id          int    `xorm:"not null pk autoincr INT(11) COMMENT('自增编号')"`
	AppraisedId int    `xorm:"not null INT(11) COMMENT('评价Id')"`
	OrderId     int    `xorm:"not null INT(2) COMMENT('顺序Id')"`
	ImgUrl      string `xorm:"VARCHAR(255) not null COMMENT('图片地址')"`
}

/*用户关注*/
type UserAttention struct {
	Id           int   `xorm:"not null pk autoincr INT(11) COMMENT('自增编号')"`
	UserId       int64 `xorm:"bigint not null COMMENT('用户Id')"`
	IsFollowedId int   `xorm:"not null INT(11) COMMENT('被关注者Id')"`
}

/*邀请用户注册*/
type InviteUserRegister struct {
	Id        int   `xorm:"not null pk autoincr INT(11) COMMENT('自增编号')"`
	Sender    int64 `xorm:"bigint not null COMMENT('发送者Id')"`
	Receiver  int64 `xorm:"bigint not null COMMENT('被邀请者Id')"`
	CreatedAt int64 `xorm:"bigint not null COMMENT('创建时间')"`
}

/*用户充值记录*/
type UserDepositInfo struct {
	Id          int64       `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	UserId      int64       `xorm:"bigint not null COMMENT('用户Id')"`
	Money       float64     `xorm:"decimal(16,2) not null COMMENT('充值金额')"`
	Platform    Platform    `xorm:"TINYINT(1) not null COMMENT('软件平台,如app,h5,pc')"`
	PayPlatform PayPlatform `xorm:"TINYINT(1) not null COMMENT('第三方付费平台,如微信,支付宝')"`
	CreatedAt   int64       `xorm:"bigint not null COMMENT('创建时间')"`
}

/*命理师充值记录*/
type FTDepositInfo struct {
	Id          int64       `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	FTId        int         `xorm:"INT(11) not null COMMENT('命理师Id')"`
	Money       float64     `xorm:"decimal(16,2) not null COMMENT('充值金额')"`
	Platform    Platform    `xorm:"TINYINT(1) not null COMMENT('软件平台,如app,h5,pc')"`
	PayPlatform PayPlatform `xorm:"TINYINT(1) not null COMMENT('第三方付费平台,如微信,支付宝')"`
	CreatedAt   int64       `xorm:"bigint not null COMMENT('创建时间')"`
}

/*用户金币明细*/
type UserGoldChange struct {
	Id         int64          `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	UserId     int64          `xorm:"bigint not null COMMENT('用户Id')"`
	ChangeType GoldChangeType `xorm:"TINYINT(1) not null COMMENT('0是充值,1是消费,2是退款')"`
	VarGold    int64          `xorm:"bigint not null COMMENT('金币变化的绝对值')"`
	CreatedAt  int64          `xorm:"bigint not null COMMENT('创建时间')"`
}

/*用户佣金收入*/
type BalanceInfo struct {
	Id            int64   `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	UserDepositId int64   `xorm:"bigint not null COMMENT('用户充值Id')"`
	AgencyLevel   int     `xorm:"TINYINT(1) not null default 1 COMMENT('代理级别')"`
	FromUserId    int64   `xorm:"bigint not null COMMENT('充值用户Id')"`
	ToUserId      int64   `xorm:"bigint not null COMMENT('收入用户Id')"`
	EarnBalance   float64 `xorm:"decimal(16,2) not null COMMENT('返现多少佣金')"`
	EarnGold      int64   `xorm:"bigint not null default 0 COMMENT('返多少金币')"`
	CreatedAt     int64   `xorm:"bigint not null COMMENT('创建时间')"`
}

/*用户提现记录*/
type UserDrawCashInfo struct {
	Id          int64               `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	UserId      int64               `xorm:"bigint not null COMMENT('用户Id')"`
	Origin      float64             `xorm:"decimal(16,2) not null COMMENT('发起的提款数目')"`
	Charge      float64             `xorm:"decimal(16,2) not null COMMENT('到账金额')"`
	Poundage    float64             `xorm:"decimal(16,2) not null COMMENT('手续费')"`
	TradeNo     string              `xorm:"VARCHAR(50) not null COMMENT('自定义交易号')"`
	State       DrawCashState       `xorm:"TINYINT(1) not null COMMENT('0为审核中,1为提现成功,2为提现失败')"`
	PaymentNo   string              `xorm:"VARCHAR(50) null COMMENT('第三方生成的交易号')"`
	PaymentTime string              `xorm:"VARCHAR(30) null COMMENT('第三方生成的交易时间')"`
	ArrivalType DrawCashArrivalType `xorm:"TINYINT(1) not null COMMENT('0为到账到微信,1为到账到支付宝')"`
	IpAddr      string              `xorm:"VARCHAR(30) null COMMENT('用户ip地址')"`
	CreatedAt   int64               `xorm:"bigint not null COMMENT('创建时间')"`
}

/*命理师提现记录*/
type FTDrawCashInfo struct {
	Id          int64               `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	FTId        int64               `xorm:"bigint not null COMMENT('用户Id')"`
	Origin      float64             `xorm:"decimal(16,2) not null COMMENT('发起的提款数目')"`
	Charge      float64             `xorm:"decimal(16,2) not null COMMENT('到账金额')"`
	Poundage    float64             `xorm:"decimal(16,2) not null COMMENT('手续费')"`
	TradeNo     string              `xorm:"VARCHAR(50) not null COMMENT('自定义交易号')"`
	State       DrawCashState       `xorm:"TINYINT(1) not null COMMENT('0为审核中,1为提现成功,2为提现失败')"`
	PaymentNo   string              `xorm:"VARCHAR(50) null COMMENT('第三方生成的交易号')"`
	PaymentTime string              `xorm:"VARCHAR(30) null COMMENT('第三方生成的交易时间')"`
	ArrivalType DrawCashArrivalType `xorm:"TINYINT(1) not null COMMENT('0为到账到微信,1为到账到支付宝')"`
	IpAddr      string              `xorm:"VARCHAR(30) null COMMENT('用户ip地址')"`
	CreatedAt   int64               `xorm:"bigint not null COMMENT('创建时间')"`
}

/*提现参数*/
type DrawCashParams struct {
	Id            int64              `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	MinCharge     float64            `xorm:"decimal(16,2) not null COMMENT('最低提现额度')"`
	MinPoundage   float64            `xorm:"decimal(16,2) not null COMMENT('最低提现手续费')"`
	MaxDrawCount  int                `xorm:"INT(11)  not null COMMENT('每日最大提现次数')"`
	PoundagePer   int                `xorm:"INT(11)  not null COMMENT('提现手续费百分比[1-100]')"`
	RequireVerify float64            `xorm:"decimal(16,2) not null COMMENT('超过多少钱需要审核')"`
	ParamsType    DrawCashParamsType `xorm:"TINYINT(1) not null COMMENT('0为用户,1为命理师')"`
}

/*用户已购买的订单*/
type UserOrderInfo struct {
	Id           int64    `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	UserId       int64    `xorm:"bigint not null COMMENT('用户Id')"`
	ProductId    int      `xorm:"not null INT(11) COMMENT('产品Id')"`
	IsPayForGold bool     `xorm:"TINYINT(1) not null COMMENT('是否用金币购买')"`
	Platform     Platform `xorm:"TINYINT(1) not null COMMENT('软件平台,如app,h5,pc')"`
}

/*订单结算*/
type UserOrderCheck struct {
	Id            int64   `xorm:"not null pk bigint COMMENT('订单Id')"`
	IsAppraised   bool    `xorm:"TINYINT(1) not null COMMENT('是否已评价')"`
	IsChecked     bool    `xorm:"TINYINT(1) not null COMMENT('是否已结算')"`
	CheckedIncome float64 `xorm:"decimal(16,2) not null COMMENT('结算后所得收益')"`
	UpdatedAt     int64   `xorm:"bigint not null COMMENT('更新时间')"`
}

/*退款结果*/
type UserOrderRefundFinished struct {
	Id         int64               `xorm:"not null pk bigint COMMENT('订单Id')"`
	RefundType UserOrderRefundType `xorm:"TINYINT(1) not null COMMENT('0为用户申请退款,1为系统自动退款')"`
	CreatedAt  int64               `xorm:"bigint not null COMMENT('创建时间')"`
}

/*退款中*/
type UserOrderRefunding struct {
	Id         int64               `xorm:"not null pk bigint COMMENT('订单Id')"`
	RefundType UserOrderRefundType `xorm:"TINYINT(1) not null default 0 COMMENT('0为用户申请退款,1为系统自动退款')"`
	CreatedAt  int64               `xorm:"bigint not null COMMENT('创建时间')"`
}

/*维权结果*/
type UserOrderRightsFinished struct {
	Id        int64 `xorm:"not null pk bigint COMMENT('订单Id')"`
	IsAgree   bool  `xorm:"TINYINT(1) not null COMMENT('是否同意退款')"`
	CreatedAt int64 `xorm:"bigint not null COMMENT('创建时间')"`
}

/*维权中*/
type UserOrderRighting struct {
	Id         int64           `xorm:"not null pk bigint COMMENT('订单Id')"`
	RightsType OrderRightsType `xorm:"TINYINT(1) not null default 0 COMMENT('0为命理师不同意退款')"`
	CreatedAt  int64           `xorm:"bigint not null COMMENT('创建时间')"`
}

/*用户购买成功通知*/
type OrderInfoMsg struct {
	Id           int64  `xorm:"not null pk bigint COMMENT('订单Id')"`
	UserId       int64  `xorm:"bigint not null COMMENT('用户Id')"`
	UserNickName string `xorm:"VARCHAR(100) null COMMENT('用户昵称')"`
	FTId         int    `xorm:"INT(11) not null pk COMMENT('命理师Id')"`
	FTNickName   string `xorm:"VARCHAR(100) null COMMENT('命理师昵称')"`
	ProductName  string `xorm:"VARCHAR(50) not null COMMENT('产品名称')"`
	UserRead     bool   `xorm:"TINYINT(1) not null default 0 COMMENT('用户是否已读')"`
	FTRead       bool   `xorm:"TINYINT(1) not null default 0 COMMENT('命理师是否已读')"`
	CreatedAt    int64  `xorm:"bigint not null COMMENT('创建时间')"`
}

/*维权结果通知*/
type OrderRightsFinishedMsg struct {
	Id           int64  `xorm:"not null pk bigint COMMENT('订单Id')"`
	UserId       int64  `xorm:"bigint not null COMMENT('用户Id')"`
	UserNickName string `xorm:"VARCHAR(100) null COMMENT('用户昵称')"`
	FTId         int    `xorm:"INT(11) not null pk COMMENT('命理师Id')"`
	FTNickName   string `xorm:"VARCHAR(100) null COMMENT('命理师昵称')"`
	ProductName  string `xorm:"VARCHAR(50) not null COMMENT('产品名称')"`
	IsAgree      bool   `xorm:"TINYINT(1) not null COMMENT('是否同意退款')"`
	UserRead     bool   `xorm:"TINYINT(1) not null default 0 COMMENT('用户是否已读')"`
	FTRead       bool   `xorm:"TINYINT(1) not null default 0 COMMENT('命理师是否已读')"`
	CreatedAt    int64  `xorm:"bigint not null COMMENT('创建时间')"`
}

/*用户收到的退款通知*/
type UserOrderRefundMsg struct {
	Id               int64                     `xorm:"not null pk bigint COMMENT('订单Id')"`
	UserId           int64                     `xorm:"bigint not null COMMENT('用户Id')"`
	FTNickName       string                    `xorm:"VARCHAR(100) null COMMENT('命理师昵称')"`
	ProductName      string                    `xorm:"VARCHAR(50) not null COMMENT('产品名称')"`
	RefundResultType UserOrderRefundResultType `xorm:"TINYINT(1) not null COMMENT('0为命理师同意退款,1为命理师不同意退款,2为系统自动退款')"`
	UserRead         bool                      `xorm:"TINYINT(1) not null default 0 COMMENT('用户是否已读')"`
	CreatedAt        int64                     `xorm:"bigint not null COMMENT('创建时间')"`
}

/*用户注册成功通知*/
type UserRegisterMsg struct {
	UserId    int64 `xorm:"bigint pk not null COMMENT('用户Id')"`
	UserRead  bool  `xorm:"TINYINT(1) not null default 0 COMMENT('用户是否已读')"`
	CreatedAt int64 `xorm:"bigint not null COMMENT('创建时间')"`
}

/*命理师注册成功通知*/
type FTRegisterMsg struct {
	FTId      int   `xorm:"INT(11) not null pk COMMENT('命理师Id')"`
	FTRead    bool  `xorm:"TINYINT(1) not null default 0 COMMENT('命理师是否已读')"`
	CreatedAt int64 `xorm:"bigint not null COMMENT('创建时间')"`
}

/*命理师收到的退款通知*/
type FTOrderRefundMsg struct {
	Id           int64               `xorm:"not null pk bigint COMMENT('订单Id')"`
	FTId         int                 `xorm:"INT(11) not null COMMENT('命理师Id')"`
	UserNickName string              `xorm:"VARCHAR(100) null COMMENT('用户昵称')"`
	ProductName  string              `xorm:"VARCHAR(50) not null COMMENT('产品名称')"`
	RefundType   UserOrderRefundType `xorm:"TINYINT(1) not null default 0 COMMENT('0为用户申请退款')"`
	FTRead       bool                `xorm:"TINYINT(1) not null default 0 COMMENT('命理师是否已读')"`
	CreatedAt    int64               `xorm:"bigint not null COMMENT('创建时间')"`
}

/*命理师信息冷数据*/
type ColdFTInfo struct {
	Id            int       `xorm:"not null pk autoincr INT(11) COMMENT('自增编号')"`
	Phone         string    `xorm:"VARCHAR(20) not null unique COMMENT('手机号')"`
	Pwd           string    `xorm:"VARCHAR(30) null COMMENT('密码')"`
	NickName      string    `xorm:"VARCHAR(100) null unique COMMENT('昵称')"`
	Avatar        string    `xorm:"VARCHAR(500) null COMMENT('头像')"`
	Sex           Sex       `xorm:"TINYINT(1) null COMMENT('性别')"`
	ActualName    string    `xorm:"VARCHAR(100) null COMMENT('真实姓名')"`
	IdentityCard  string    `xorm:"CHAR(20) null COMMENT('身份证')"`
	IdFrontCover  string    `xorm:"VARCHAR(500) null COMMENT('身份证图片正面地址')"`
	IdBehindCover string    `xorm:"VARCHAR(500) null COMMENT('身份证图片反面地址')"`
	Introduction  string    `xorm:"VARCHAR(600) null COMMENT('个人介绍')"`
	Registration  Platform  `xorm:"TINYINT(1) not null COMMENT('注册平台')"`
	AuthState     AuthState `xorm:"TINYINT(1) not null COMMENT('0审核中,1审核失败,2审核成功')"`
	CreatedAt     int64     `xorm:"bigint not null COMMENT('注册时间')"`
}

/*命理师信息热数据*/
type HotFTInfo struct {
	FTId         int          `xorm:"INT(11) not null pk COMMENT('命理师Id')"`
	Token        string       `xorm:"CHAR(40) not null COMMENT('令牌值')"`
	LoginTime    int64        `xorm:"bigint not null COMMENT('登录时间')"`
	Account      int64        `xorm:"bigint not null default 0 COMMENT('当前积分')"`
	DepositTotal float64      `xorm:"decimal(16,2) not null default 0 COMMENT('总充值额')"`
	Balance      float64      `xorm:"decimal(16,2) not null default 0 COMMENT('当前剩下的金额')"`
	BalanceTotal float64      `xorm:"decimal(16,2) not null default 0 COMMENT('赚取的总金额')"`
	PayTotal     float64      `xorm:"decimal(16,2) not null default 0 COMMENT('积分消费总额')"`
	AccountState AccountState `xorm:"TINYINT(1) not null default 0 COMMENT('账户状态,0是正常,1是黑名单,2是冻结')"`
	Mark         string       `xorm:"VARCHAR(255) null COMMENT('多个标签内容')"`
	AutoReply    string       `xorm:"VARCHAR(255) null COMMENT('接待语')"`
	SortId       int          `xorm:"INT(11) not null default 0 COMMENT('指定排序Id')"`
	OrderScore   float64      `xorm:"decimal(16,2) not null default 0 COMMENT('排序得分')"`
	IncomePer    float64      `xorm:"decimal(16,2) not null default 60 COMMENT('收益提成百分比,[0,100]')"`
	IMPrivateKey string       `xorm:"CHAR(255) not null COMMENT('IM私钥')"`
	EnableFree   bool         `xorm:"TINYINT(1) not null default 1 COMMENT('是否开启免费对话')"`
}

/*商铺信息*/
type ShopInfo struct {
	Id            int     `xorm:"not null pk INT(11) autoincr COMMENT('自增编号')"`
	FTId          int     `xorm:"not null INT(11) COMMENT('命理师Id')"`
	Score         float64 `xorm:"decimal(16,2) not null default 0 COMMENT('店铺得分')"`
	ResponseSpeed float64 `xorm:"decimal(16,2) not null default 0 COMMENT('响应速度')"`
	Visits        int64   `xorm:"not null bigint default 0 COMMENT('访问次数')"`
	IsRecommended bool    `xorm:"TINYINT(1) not null default 0 COMMENT('是否被推荐')"`
}

/*商铺图片*/
type ShopImgs struct {
	Id     int    `xorm:"not null pk INT(11) autoincr COMMENT('自增编号')"`
	ShopId int    `xorm:"not null INT(11) COMMENT('店铺Id')"`
	ImgUrl string `xorm:"VARCHAR(255) not null COMMENT('图片地址')"`
}

// /*大师在线时间段*/
// type FTOnlineTime struct {
// 	Id        int    `xorm:"not null pk INT(11) autoincr COMMENT('自增编号')"`
// 	FTId      int    `xorm:"not null INT(11) COMMENT('命理师Id')"`
// 	StartTime string `xorm:"not null CHAR(8) default '00:00' COMMENT('开始时间段,格式:01:05')"`
// 	EndTime   string `xorm:"not null CHAR(8) default '23:59' COMMENT('结束时间段,格式:01:05')"`
// }

/*产品信息*/
type ProductInfo struct {
	Id          int     `xorm:"not null pk INT(11) autoincr COMMENT('自增编号')"`
	ShopId      int     `xorm:"not null INT(11) default 0 COMMENT('所属店铺Id')"`
	ProductName string  `xorm:"VARCHAR(50) not null default '' COMMENT('产品名称')"`
	ProductDesc string  `xorm:"VARCHAR(500) not null default '' COMMENT('产品描述')"`
	Price       float64 `xorm:"decimal(16,2) not null default 0 COMMENT('产品价格')"`
	IsHidden    bool    `xorm:"TINYINT(1) not null default 1 COMMENT('是否隐藏')"`
	SortId      int     `xorm:"not null INT(4) default 0 COMMENT('排序Id')"`
	CreatedAt   int64   `xorm:"bigint not null default 0 COMMENT('创建时间')"`
	UpdatedAt   int64   `xorm:"bigint not null default 0 COMMENT('修改时间')"`
}

/*命理师认证*/
type Authentication struct {
	FTId        int         `xorm:"not null pk INT(11) COMMENT('命理师Id')"`
	IdCardState IdCardState `xorm:"TINYINT(1) not null default 0 COMMENT('0是未提交,1是已提交,2是未通过,3是已通过')"`
	IsCP        bool        `xorm:"TINYINT(1) not null default 0 COMMENT('是否消费保障')"`
	IsHR        bool        `xorm:"TINYINT(1) not null default 0 COMMENT('是否金牌推荐')"`
}

/*命理师缴纳消费保障记录*/
type PayConsumerProtection struct {
	Id        int     `xorm:"not null pk INT(11) autoincr COMMENT('自增编号')"`
	FTId      int     `xorm:"not null INT(11) COMMENT('命理师Id')"`
	Money     float64 `xorm:"decimal(16,2) not null  COMMENT('缴纳金额')"`
	CreatedAt int64   `xorm:"bigint not null COMMENT('创建时间')"`
}

/*命理师的免打扰名单*/
type FTDndList struct {
	Id        int   `xorm:"not null pk INT(11) autoincr COMMENT('自增编号')"`
	FTId      int   `xorm:"not null INT(11) COMMENT('命理师Id')"`
	UserId    int64 `xorm:"bigint not null COMMENT('用户Id')"`
	CreatedAt int64 `xorm:"bigint not null COMMENT('创建时间')"`
}

/*正在进行的聊天列表*/
type ChattingList struct {
	Id        int64 `xorm:"not null pk bigint autoincr COMMENT('自增编号')"`
	FTId      int   `xorm:"not null INT(11) COMMENT('命理师Id')"`
	UserId    int64 `xorm:"bigint not null COMMENT('用户Id')"`
	IsFree    bool  `xorm:"TINYINT(1) not null COMMENT('是否免费对话')"`
	UpdatedAt int64 `xorm:"bigint not null COMMENT('会话创建时间')"`
}

/*结束的聊天列表*/
type ChatEndList struct {
	Id        int64 `xorm:"not null pk bigint autoincr COMMENT('自增编号')"`
	FTId      int   `xorm:"not null INT(11) COMMENT('命理师Id')"`
	UserId    int64 `xorm:"bigint not null COMMENT('用户Id')"`
	IsFree    bool  `xorm:"TINYINT(1) not null COMMENT('是否免费对话')"`
	UpdatedAt int64 `xorm:"bigint not null COMMENT('会话创建时间')"`
}

/*命理师积分明细*/
type FTAccountChange struct {
	Id         int64           `xorm:"not null pk autoincr bigint COMMENT('自增编号')"`
	FTId       int             `xorm:"not null INT(11) COMMENT('命理师Id')"`
	ChangeType ScoreChangeType `xorm:"TINYINT(1) not null COMMENT('0是充值,1是消费')"`
	VarAccount int64           `xorm:"bigint not null COMMENT('积分变化的绝对值')"`
	CreatedAt  int64           `xorm:"bigint not null COMMENT('创建时间')"`
}

/*命理师快捷语*/
type FTQuickReply struct {
	Id    int    `xorm:"not null pk autoincr INT(11) COMMENT('自增编号')"`
	FTId  int    `xorm:"not null INT(11) COMMENT('命理师Id')"`
	Reply string `xorm:"VARCHAR(100) not null COMMENT('快捷回复内容')"`
}

/*命理师标签管理*/
type FTMarkInfo struct {
	Id   int    `xorm:"not null pk autoincr INT(11) COMMENT('自增编号')"`
	Desc string `xorm:"VARCHAR(100) not null COMMENT('标签内容')"`
}

/*域名管理*/
type DomainInfo struct {
	Id                  int    `xorm:"not null pk INT(2) COMMENT('编号')"`
	QrAddr              string `xorm:"VARCHAR(100) not null COMMENT('二维码入口(非炮灰)')"`
	EntryAddr           string `xorm:"VARCHAR(100) not null COMMENT('入口落地页(炮灰)')"`
	AuthAddr            string `xorm:"VARCHAR(100) not null COMMENT('授权中转页(非炮灰)')"`
	AppAddr             string `xorm:"VARCHAR(100) not null COMMENT('app落地页(炮灰)')"`
	FTDepositAddr       string `xorm:"VARCHAR(100) not null COMMENT('大师充值页(非炮灰)')"`
	UserDepositAddr     string `xorm:"VARCHAR(100) not null COMMENT('用户充值页(非炮灰)')"`
	DownLoadAddr        string `xorm:"VARCHAR(100) not null COMMENT('下载中转页(非炮灰)')"`
	DownLoadUIAddr      string `xorm:"VARCHAR(100) not null COMMENT('下载展示页(炮灰)')"`
	IosDownLoadAddr     string `xorm:"VARCHAR(100) not null COMMENT('ios下载页')"`
	AndroidDownLoadAddr string `xorm:"VARCHAR(100) not null COMMENT('android下载页')"`
	IosUrlSchemes       string `xorm:"VARCHAR(100) not null COMMENT('ios urlSchemes')"`
	AndroidUrlSchemes   string `xorm:"VARCHAR(100) not null COMMENT('android urlSchemes')"`
}

/*客服信息*/
type CustomerServiceInfo struct {
	Id        int    `xorm:"not null pk INT(2) COMMENT('编号')"`
	WX        string `xorm:"VARCHAR(30) not null COMMENT('微信号')"`
	QQ        string `xorm:"VARCHAR(30) not null COMMENT('qq号')"`
	QrcodeUrl string `xorm:"VARCHAR(255) not null COMMENT('二维码图片地址')"`
}

/*服务器状态*/
type ServerInfo struct {
	Id         int    `xorm:"not null pk INT(2) COMMENT('编号')"`
	Version    string `xorm:"CHAR(10) not null COMMENT('版本号')"`
	IsMaintain bool   `xorm:"TINYINT(1) not null COMMENT('是否维护')"`
	GzhAppid   string `xorm:"CHAR(50) null COMMENT('公众号Appid')"`
	KfptAppid  string `xorm:"CHAR(50) null COMMENT('开发平台Appid')"`
}

/*敏感词*/
type SensitiveWord struct {
	Id   int64  `xorm:"not null pk bigint COMMENT('编号')"`
	Word string `xorm:"VARCHAR(255) not null COMMENT('敏感词')"`
}
