package svsdk

import "koudai-box/conf"

//获取sn
func GetSN() string {
	SN := conf.GetConf().DefaultSN
	return SN
}
