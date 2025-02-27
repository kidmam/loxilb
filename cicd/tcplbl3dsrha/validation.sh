#!/bin/bash
source ../common.sh
echo SCENARIO-tcplbl3dsr
$hexec ep1 socat -v -T0.05 tcp-l:8080,reuseaddr,fork system:"echo 'server1'; cat" >/dev/null 2>&1 &
$hexec ep2 socat -v -T0.05 tcp-l:8080,reuseaddr,fork system:"echo 'server2'; cat" >/dev/null 2>&1 &
$hexec ep3 socat -v -T0.05 tcp-l:8080,reuseaddr,fork system:"echo 'server3'; cat" >/dev/null 2>&1 &

sleep 5
code=0
servArr=( "server1" "server2" "server3" )
ep=( "56.56.56.1" "57.57.57.1" "58.58.58.1" )
j=0
waitCount=0
while [ $j -le 2 ]
do
    #res=$($hexec user curl --max-time 10 -s ${ep[j]}:8080)
    res=$($hexec user socat -T60 - TCP:${ep[j]}:8080,reuseaddr)
    #echo $res
    if [[ $res == "${servArr[j]}" ]]
    then
        echo "$res UP"
        j=$(( $j + 1 ))
    else
        echo "Waiting for ${servArr[j]}(${ep[j]})"
        waitCount=$(( $waitCount + 1 ))
        if [[ $waitCount == 10 ]];
        then
            echo "All Servers are not UP"
            echo SCENARIO-tcplbl3dsr [FAILED]
            exit 1
        fi
    fi
    sleep 1
done

exp=""
for i in {1..4}
do
for j in {0..2}
do
    #res=$($hexec user curl --local-port 55001 --max-time 10 -s 20.20.20.1:2020)
    #res=$($hexec user socat -T60 - TCP:20.20.20.1:8080,sp=15402,reuseaddr)
    res=$($hexec user socat -T60 - TCP:20.20.20.1:8080)
    echo $res
    if [[ $exp == "" ]]
    then
      exp=$res
    fi
    if [[ $exp != $res ]]
    then
      code=1
    fi
    sleep 1
done
done
if [[ $code == 0 ]]
then
    echo SCENARIO-tcplbl3dsr [OK]
else
    echo SCENARIO-tcplbl3dsr [FAILED]
fi
sudo pkill -9 socat 2>&1 > /dev/null
exit $code

