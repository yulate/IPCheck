package core

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"regexp"
)

func Check(path string) {

}

// WBCheck 使用微步接口对ip进行扫描，需要token
func WBCheck(url string) {
	//url := "https://x.threatbook.com/v5/ip/8.8.8.8?source=top&isRange=true"
	cookie := "rememberme=25f073452327f0186b2772f64ec12077f988a98e|a69119a8753a4e2390f990c4d7076950|1686494126955|public|w; xx-csrf=25f073452327f0186b2772f64ec12077f988a98e; csrfToken=mWRFtEcTh-R8UxyMZ-UQP5x5; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22a69119a8753a4e2390f990c4d7076950%22%2C%22first_id%22%3A%22182e3aeda66423-044bfe089185ef4-72422e2e-1693734-182e3aeda6714a0%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E8%87%AA%E7%84%B6%E6%90%9C%E7%B4%A2%E6%B5%81%E9%87%8F%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC%22%2C%22%24latest_referrer%22%3A%22https%3A%2F%2Fwww.google.com.hk%2F%22%7D%2C%22identities%22%3A%22eyIkaWRlbnRpdHlfY29va2llX2lkIjoiMTgyZTNhZWRhNjY0MjMtMDQ0YmZlMDg5MTg1ZWY0LTcyNDIyZTJlLTE2OTM3MzQtMTgyZTNhZWRhNjcxNGEwIiwiJGlkZW50aXR5X2xvZ2luX2lkIjoiYTY5MTE5YTg3NTNhNGUyMzkwZjk5MGM0ZDcwNzY5NTAifQ%3D%3D%22%2C%22history_login_id%22%3A%7B%22name%22%3A%22%24identity_login_id%22%2C%22value%22%3A%22a69119a8753a4e2390f990c4d7076950%22%7D%2C%22%24device_id%22%3A%22182e3aeda66423-044bfe089185ef4-72422e2e-1693734-182e3aeda6714a0%22%7D; day_first_activity=true; day_first=true"
	client := resty.New()
	resp, err := client.R().EnableTrace().
		SetHeader("Cookie", cookie).
		Get(url)
	if err != nil {
		panic(err)
	}

	//fmt.Println("resp:", resp)
	// 处理出关键数据部分
	pattern := `window\.__INITIAL_STATE__\s*=\s*({.*?});`
	re := regexp.MustCompile(pattern)

	match := re.FindStringSubmatch(resp.String())
	if len(match) > 1 {
		// 提取JSON数据
		jsonData := match[1]
		fmt.Println(jsonData)
	} else {
		fmt.Println("未找到匹配的JSON数据")
	}
}
