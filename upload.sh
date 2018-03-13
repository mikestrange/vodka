#!/bin/bash
#scp -r ./* root@120.77.149.74:~/home/godark


list=`ls |grep -v 'pkg\|bin\|upload.sh\|debug.log\|README.md'`
echo $list

for data in ${list[@]}
do
scp -r ${data} root@120.77.149.74:~/home/godark
done

exit