#!/bin/sh

trap '' 1 3 15

function GetPid()
{
  if [ -f $SVC_PID ]; then
    read PROC_ID < $SVC_PID
  else
    PROC_ID=0
  fi

  if [ -f $MONITOR_PID ]; then
    read SAFE_ID < $MONITOR_PID
  else
    SAFE_ID=0
  fi
}

function Mon()
{
  while true
  do
    sleep 10
    GetPid
    if test "$PROC_ID" == "0"; then
       echo "$(date) $SVC_BIN exit, safe down"  >> $SVC_LOG
       exit
    else
      kill -0 $PROC_ID > /dev/null 2>&1
      if test "$?" != "0"; then
	check_stop
	if test "$?" != "0"; then
		exit
	fi
        echo "$(date) $SVC_BIN has gone, will restart"  >> $SVC_LOG
        Boot
      fi
    fi
  done
}

function Boot()
{
 # echo $SVC_CMD >>Boot.log
  $SVC_CMD 2> /dev/null
  if test "$?" != "0"; then
    echo "$(date) $SVC_BIN start failed"  >> $SVC_LOG
    exit
  else
    sleep 2
    GetPid
    echo "$(date) $SVC_BIN started, pid: $PROC_ID"  >> $SVC_LOG
    echo $$ > $MONITOR_PID
  fi
}

function Start()
{
  if [ ! -f "check_start.txt" ]
        then
        touch check_start.txt
        echo "core_num:0">>check_start.txt
  fi
  GetPid
  if test "$PROC_ID" == "0"; then
    Boot
    Mon
  else
    kill -0 $PROC_ID > /dev/null 2>&1
    if test "$?" != "0"; then
      Boot
      Mon
    else
      echo "$(date) $SVC_BIN already running, pid: $PROC_ID"
    fi
  fi
}

#SIGINT to stop
function Stop()
{
  GetPid
  if test "$PROC_ID" == "0"; then
    echo "$(date) $SVC_BIN not running"
  else
    kill -2 $PROC_ID > /dev/null 2>&1
    echo "$(date) kill -2 $PROC_ID"  >> $SVC_LOG
    rm $SVC_PID > /dev/null 2>&1
  fi

  if test "$SAFE_ID" == "0"; then
    echo "$(date) monitor not running"
  else
    kill -9 $SAFE_ID
    echo "$(date) kill -9 $SAFE_ID"
    rm $MONITOR_PID > /dev/null 2>&1
  fi
}

#SIGQUIT to safe stop
function SafeStop()
{
  GetPid
  if test "$PROC_ID" == "0"; then
    echo "$(date) $SVC_BIN not running"
  else
    echo "$(date) kill -3 $PROC_ID"
    kill -3 $PROC_ID > /dev/null 2>&1
	#检查程序是否退出
	while true
	do
		sleep 1
		kill -0 $PROC_ID > /dev/null 2>&1
		if test "$?" != "0"; then
			echo "$(date) $SVC_BIN has gone, safe stop finish"  >> $SVC_LOG
			break;
		fi
	done
    rm $SVC_PID > /dev/null 2>&1
  fi

  if test "$SAFE_ID" == "0"; then
    echo "$(date) monitor not running"
  else
    kill -9 $SAFE_ID
    echo "$(date) kill -9 $SAFE_ID"
    rm $MONITOR_PID > /dev/null 2>&1
  fi
}

function Stat()
{
  GetPid
  if test "$PROC_ID" == "0"; then
    echo "$(date) $SVC_BIN not running"
  else
    kill -0 $PROC_ID > /dev/null 2>&1
    if test "$?" != "0"; then
      echo "$(date) $SVC_BIN not running"
    else
      echo "$(date) $SVC_BIN ok, pid: $PROC_ID"
    fi
  fi
}

function Usage()
{
  echo
  echo "usage: safe_svrd.sh SVRD (start|stop|stat|safestop)"
  echo "   eg: safe_svrd.sh commsvrd start"
  echo
}

function Clear()
{
  echo 'clear()'
  for id in ./*;
  do
  if [ "${id##*.}" = "shmid" ];
  then
  	if [[ "${id}" =~ "\c*COMMPIPE\c*" ]]
	then
		continue
	elif [[ "${id}" =~ "\c*PKGSTAT\c*" ]];
	then
		continue
	else
	ipcrm -m $(cat $id)
	fi
  fi
  done

  echo $SVC_EXEC
}

function find()
{
	for id in ./*/*
	do
	if [ "${id##*.}" = "shmid" ];
	then
		if [[ $SVC_EXEC -eq $(cat $id) ]];
		then
		echo $id
		break
		fi
	fi
	done

}

function check_stop()
{
  if [ ! -f "check_start.txt" ]
        then
        touch check_start.txt
        echo "core_num:0">>check_start.txt
  fi

  DATE=`date "+%Y-%m-%d %H:%M:%S"`
  seconds=`date -d "$DATE" +%s`
  cur_num=`cat check_start.txt |grep core_num:|cut -d ':' -f 2`
  start_index=`expr $cur_num - 3`
  if [ $start_index -lt 0 ]
        then
        time_start=0
  else
        time_start=`cat check_start.txt |grep "$start_index:"|cut -d ':' -f 2`
  fi
  time_accept=`expr $time_start + 60`
  if [ "$time_start" = "" -o $time_accept -lt $seconds ]
        then
        return_val=0
  else
        return_val=1
  fi
  echo "`expr $cur_num + 1`:$seconds">>check_start.txt
  mod_cur_num=`expr $cur_num + 1`
  sed -i "s/core_num:$cur_num/core_num:$mod_cur_num/g" check_start.txt;
  return $return_val
}

if [ "$#" -lt "2" ]
then
  Usage
  exit
fi

SVC_EXEC=`echo $1|awk '{print $1}'`
#SVC_EXEC=$1    //changed by kavin. to support arglist.
SVC_BIN="$SVC_EXEC"
SVC_PID="$SVC_BIN".pid
SVC_CMD=$1     
#SVC_CMD=$SVC_BIN //changed by kavin. 
SVC_LOG=safe_svrd."$SVC_EXEC".log

MONITOR_PID=safe_svrd."$SVC_EXEC".pid
CURRENT=$(cd `dirname $0`; pwd)

export LD_LIBRARY_PATH=$CURRENT:$LD_LIBRARY_PATH
export PATH=$CURRENT:$PATH

if test "$2" == "start"; then
  ulimit -c unlimited
  Start
elif test "$2" == "stop"; then
  Stop
elif test "$2" == "safestop"; then
  SafeStop
elif test "$2" == "stat"; then
  Stat
elif test "$2" == "restart"; then
  Stop
  Start
elif test "$2" == "clear"; then
  Clear
elif test "$2" == "find"; then
  find
else
  Usage
fi

exit
