package core

// WBResult 微步请求结果
type WBResult struct {
	Csrf  string `json:"csrf"`
	Title string `json:"title"`
	Data  struct {
		UserInfo struct {
			GraphDisable  bool   `json:"graphDisable"`
			UserId        string `json:"userId"`
			IsLogin       bool   `json:"isLogin"`
			Id            string `json:"id"`
			UserImg       string `json:"userImg"`
			NickName      string `json:"nickName"`
			LastLoginTime int64  `json:"lastLoginTime"`
			Email         string `json:"email"`
			PhoneNumber   string `json:"phoneNumber"`
			Company       bool   `json:"company"`
		} `json:"userInfo"`
		Resource      string `json:"resource"`
		ResourceName  string `json:"resourceName"`
		QuerySettings struct {
			Id               string `json:"id"`
			IocSummaryEnable bool   `json:"iocSummaryEnable"`
		} `json:"querySettings"`
		ResourceType string `json:"resourceType"`
		SummaryInfo  struct {
			Judge     int `json:"judge"`
			OpenJudge int `json:"openJudge"`
			UserJudge int `json:"userJudge"`
			Judgments []struct {
				Name string `json:"name"`
				Type int    `json:"type"`
			} `json:"judgments"`
			Events []struct {
				Id       string `json:"id"`
				Name     string `json:"name"`
				Severity string `json:"severity"`
			} `json:"events"`
			Community []interface{} `json:"community"`
			Scene     string        `json:"scene"`
			Type      string        `json:"type"`
			Location  struct {
				Country     string  `json:"country"`
				Lng         float64 `json:"lng"`
				Lat         float64 `json:"lat"`
				Carrier     string  `json:"carrier"`
				CountryCode string  `json:"country_code"`
			} `json:"location"`
			UpdateTime     string `json:"update_time"`
			DomainCount    int    `json:"domain_count"`
			PortCount      int    `json:"port_count"`
			UrlCount       string `json:"url_count"`
			SampleCount    int    `json:"sample_count"`
			WebSampleCount string `json:"web_sample_count"`
			IpClaim        string `json:"ip_claim"`
		} `json:"summaryInfo"`
		Query struct {
			Source  string `json:"source"`
			IsRange string `json:"isRange"`
		} `json:"query"`
	} `json:"data"`
}
