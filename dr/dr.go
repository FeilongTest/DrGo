package dr

import (
	"fmt"
	"github.com/levigross/grequests"
	"log"
	"strconv"
	"strings"
)

//Login 登录
func Login(account, pwd string) bool {
	resp, err := grequests.Post(
		"http://10.10.10.5",
		&grequests.RequestOptions{
			Data: map[string]string{
				"DDDDD":  account,
				"upass":  pwd,
				"m1":     "000000000000",
				"0MKKey": "0123456789",
				"ver":    "1.5.100.202101112.G.L.A.D",
				"sim_sp": "undefine",
				"cver1":  "1",
				"cver2":  "0010501000&R6=1",
			},
			Headers: map[string]string{
				"Uip":     "va5=1.2.3.4.7d90b22d5534741c6057abd51f0317a090824c79",
				"Charset": "utf-8",
			},
			UserAgent: "DrCOM-Http",
		})
	if err != nil {
		log.Println("err", err)
	} else {
		if strings.Index(resp.String(), "<!--Dr.COMWebLoginID_3.htm-->") != -1 {
			return true
		} else {
			return false
		}
	}
	return false
}

//Logout 注销
func Logout() bool {
	resp, err := grequests.Get(
		"http://10.10.10.5/F.html",
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Charset": "utf-8",
			},
			UserAgent: "DrCOM-HttpClient",
		})
	if err != nil {
		log.Println("err", err)
	} else {
		if strings.Index(resp.String(), "<!--Dr.COMWebLoginID_2.htm-->") != -1 {
			return true
		} else {
			return false
		}
	}
	return false
}

//GetBalance 获取账户余额
func GetBalance() string {
	resp, err := grequests.Get("http://10.10.10.5/", &grequests.RequestOptions{
		UserAgent: "DrCOM-Http",
	})
	if err != nil {
		log.Println("getBalance 失败", err)
	} else {
		balanceStr := GetBetweenStr(resp.String(), "fsele=1;fee='", "';xsele=0")
		balanceStr = strings.TrimSpace(balanceStr)
		balance, _ := strconv.ParseFloat(balanceStr, 64)
		balanceStr = fmt.Sprintf("%.2f", balance/10000)
		return balanceStr
	}
	return "UnKnow"
}

//GetBetweenStr 取中间文本
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start)
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
