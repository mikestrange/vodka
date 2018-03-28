#!/bin/bash
#scp -r ./* root@120.77.149.74:~/home/ants

#移除不要的
rm -f src/src

list=`ls |grep -v 'pkg\|bin\|upload.sh\|debug.log\|history\|README.md'`
#echo $list

for data in ${list[@]}
do
scp -r ${data} root@hotbeel.com:~/home/ants
done

#直接运行服务器(直接运行分离服务器)
ssh root@hotbeel.com "cd ~/home/ants; ls -a;sh make.sh disnets;exit;"

exit