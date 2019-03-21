package datastruct

type UserLoginWithPwdBody struct {
	Phone string `json:"phone"`
	Pwd   string `json:"pwd"`
}

type UserRegisterBody struct {
	Phone    string   `json:"phone"`
	Pwd      string   `json:"pwd"`
	Platform Platform `json:"platform"`
}

type UserRegisterDetailBody struct {
	UserRegisterBody
	NickName   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	Sex        Sex    `json:"sex"`
	ActualName string `json:"actualname"`
	Birthday   string `json:"birthday"`
	City       string `json:"city"`
}

type FTRegisterBody struct {
	Phone    string   `json:"phone"`
	Pwd      string   `json:"pwd"`
	NickName string   `json:"nickname"`
	Avatar   string   `json:"avatar"`
	Mark     string   `json:"mark"` //逗号分隔
	Desc     string   `json:"desc"` //个人介绍
	Platform Platform `json:"platform"`
}

type FTRegisterWithIDBody struct {
	FTRegisterBody
	FtIdentity
}

type FtLoginBody struct {
	Phone string `json:"phone"`
	Pwd   string `json:"pwd"`
}

type UpdateFtInfoBody struct {
	Avatar   string `json:"avatar"`
	NickName string `json:"nickname"`
}

type UpdateFtMarkBody struct {
	Mark string `json:"mark"`
}

type UpdateFtIntroductionBody struct {
	Imgs []string `json:"imgs"`
	Desc string   `json:"desc"`
}

type UpdateFtAutoReplyBody struct {
	AutoReply  string   `json:"autoreply"`
	QuickReply []string `json:"quickreply"`
}

type FtIdentity struct {
	Identity      string `json:"identity"`
	ActualName    string `json:"actualname"`
	IdFrontCover  string `json:"front"`
	IdBehindCover string `json:"behind"`
}

type RemoveWithIdBody struct {
	Id int `json:"id"`
}

type RemoveWithIdsBody struct {
	Ids []int `json:"ids"`
}

type EditProductBody struct {
	Id          int     `json:"id"`
	ProductName string  `json:"name"`
	ProductDesc string  `json:"desc"`
	Price       float64 `json:"price"`
	IsHidden    bool    `json:"ishidden"`
}

type FakeAppraisedBody struct {
	Id    int     `json:"id"`    //产品ID
	Time  int64   `json:"time"`  //评价时间戳
	Score float64 `json:"score"` //评分数
	Mark  string  `json:"mark"`  //标签
	Desc  string  `json:"desc"`  //评价内容
}

type IsAgreeRefundBody struct {
	Id      int64 `json:"id"` //产品ID
	IsAgree bool  `json:"isagree"`
}

type QueryFtListBody struct {
	Key string `json:"key"`
}
