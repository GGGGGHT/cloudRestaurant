package param

type LoginParam struct {
	// 用户名
	Name string `json:"name"`
	// 密码
	Password string `json:"pwd"`
	// 验证码ID
	Id string `json:"id"`
	// 验证码
	Value string `json:"value"`
}
