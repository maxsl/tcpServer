#!/bin/bash  
   
SERVICE="/root/talkServer/bin/talkServer"  
SERVICE_LOG="/root/talkServer/bin/server.log"  
SERVICE_ERROR="/root/talkServer/bin/server.err"  
ERRORLOG="/root/talkServer/bin/errorlog"  
  
  
while true  
do  
    PID_COUNT=`ps aux|grep "$SERVICE" |grep -v grep |wc -l`  
      
    if [ $PID_COUNT -eq 0 ]  
    then  
        [ ! -e $SERVICE ] && echo "ERROR: $SERVICE not exists." >> $ERRORLOG && exit  
        nohup $SERVICE  >>$SERVICE_LOG  &  
    fi  
    sleep 3 
done  
  
exit 
