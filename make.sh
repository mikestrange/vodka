#!/bin/sh

#查找进程
# ps aux | grep godark
#直接构建
#go build $dir/src/main.go
#构建进程名称
#go build -o $NAME $dir/src/main.go

echo "=======================BEGIN======================="
NAME="godark"
#这里关闭当前进程
#ID=`ps -ef | grep "$NAME" | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
ID=`ps aux | grep "$NAME" | grep -v "grep" | awk '{print $2}'`
echo "---------------"
for id in $ID
do
kill -9 $id
echo "killed $id"
done
echo "---------------"

if [ $# == 0 ]; then
    echo "[no args kill all]"
    echo "========================OVER========================"
    exit
fi

##注入GOPATH
dir=$(cd `dirname $0`; pwd)
echo $dir
export GOPATH=$dir
##构建环境和目录
go build -o $NAME $dir/src/main.go
cd $dir
#判断启动方式
if [[ $1 == "disnets" ]]; then
##test 本地运行
./$NAME ser gate > history/gate_log.log &
./$NAME ser world > history/world_log.log &
./$NAME ser login > history/login_log.log &
./$NAME ser chat > history/chat_log.log &
else
./$NAME $@ > debug.log &
fi

#上面已经进入后台
#bg %1
echo "========================OVER========================"
exit