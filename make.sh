#!/bin/sh

# ps aux | grep godark

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
    echo "no args"
    exit
fi

##
dir=$(cd `dirname $0`; pwd)
echo $dir
export GOPATH=$dir
#go build $dir/src/main.go
#编辑进程名称
go build -o $NAME $dir/src/main.go
cd $dir
./$NAME $@ > debug.log &
echo "start server:"$@
#上面已经进入后台
#bg %1

exit