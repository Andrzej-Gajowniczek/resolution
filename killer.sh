kill -9 `ps aux |grep main |grep go|awk '{print$2}'|tr '\n' ' '`
