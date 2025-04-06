# Distributed-Logs

# Features

1. Distributed logs with async replication.  
2. Explicit acknowledgement of data read.  
3. Allow users to choose in-memory or on-disk log storage.  

# Design
  
1. Data read and split into chunks on disc to ensure no data is loss.  
2. On disc data are replicated.  

HTTP main server?  
 - easiest  
  
gRPC?  
 - need to allocate memory more than http? 

Make a custom protocol for the distributed streaming service?  


### Practice Server interactions
  
```
$ echo "hello hello hello" | curl --data-binary @- 'http://localhost:8080/write';echo
$ curl 'http://localhost:8080/read';echo
```

sending "hello hello hello" through the terminal to /write.  
/read sending it back  

