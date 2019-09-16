/**
 * @Author: DollarKiller
 * @Description: 项目初始化
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 11:47 2019-09-16
 */
package initialization

import (
	"easydevops/config"
	"fmt"
	"github.com/dollarkillerx/easyutils"
	"github.com/dollarkillerx/easyutils/clog"
	"io/ioutil"
	"strings"
)

func Initialization() {
	// 如果用户填写了配置文件 就生成 devops.sh
	createDevops()
}

func createDevops() {
	// 创建sh 脚本存放目录
	err := easyutils.DirPing("./sh")
	if err != nil {
		clog.PrintWa("脚本存放目录 创建失败")
		panic(err)
	}

	// 遍历整个 devops > node
	for _, k := range config.Basis.Devops.Node {
		// 脚本名称
		name := strings.Replace(k.FullName, "/", "1", -1)
		shName := name + k.Port + k.Branch + ".sh"
		shPath := "./sh/" + shName

		// 获取git目录地址
		gitpath := getgitpath(k.FullName)
		path := k.Dirpath
		if string(k.Dirpath[len(k.Dirpath)-1]) == "/" {
			path = k.Dirpath[:len(k.Dirpath)-1]
		}
		runpath := path + "/" + gitpath

		// 如果没有二级目录
		if k.Secondarydirectory == "" {
			sh := fmt.Sprintf(sh1, k.Port, runpath, k.Dirpath, k.Branch, k.Giturl, gitpath, k.Runname, runpath, k.Runname)
			err := ioutil.WriteFile(shPath, []byte(sh), 00666)
			if err != nil {
				clog.PrintWa("sh文件创建失败")
				panic(err)
			}
		} else {
			// 如果存在二级目录
			sh := fmt.Sprintf(sh2, k.Port, runpath, k.Dirpath, k.Branch, k.Giturl, gitpath, k.Secondarydirectory, k.Runname, runpath, k.Secondarydirectory, k.Runname)
			err := ioutil.WriteFile(shPath, []byte(sh), 00666)
			if err != nil {
				clog.PrintWa("sh文件创建失败")
				panic(err)
			}
		}
	}
	clog.Println("sh脚本初始化完毕")
}

func getgitpath(full string) string {
	index := strings.Index(full, "/")

	return full[(index + 1):]
}

var sh1 = `
#! /bin/sh

lsof -i :%s | awk '{print $2}'> tmp
pid=$(awk 'NR==2{print}' tmp);
kill -9 $pid

# 如果文件不存在 git clone 
if [ ! -d "%s" ];then
	cd %s
	git clone -b %s  %s
	cd %s
	./%s &	
else
# 如果文件存在
	cd %s
	git pull
	./%s &

fi
`

var sh2 = `
#! /bin/sh

lsof -i :%s | awk '{print $2}'> tmp
pid=$(awk 'NR==2{print}' tmp);
kill -9 $pid

# 如果文件不存在 git clone 
if [ ! -d "%s" ];then
	cd %s
	git clone -b %s  %s
	cd %s
	cd %s
	./%s &	
else
# 如果文件存在
	cd %s
	git pull
	cd %s
	./%s &

fi
`
