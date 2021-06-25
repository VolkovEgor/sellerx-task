## Нагрузочное тестирование

```
$ ab -c 10 -n 10000 localhost:9000/chats/get
This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)


Server Software:
Server Hostname:        localhost
Server Port:            9000

Document Path:          /chats/get
Document Length:        33 bytes

Concurrency Level:      10
Time taken for tests:   166.873 seconds
Complete requests:      10000
Failed requests:        0
Non-2xx responses:      10000
Total transferred:      1720000 bytes
HTML transferred:       330000 bytes
Requests per second:    59.93 [#/sec] (mean)
Time per request:       166.873 [ms] (mean)
Time per request:       16.687 [ms] (mean, across all concurrent requests)
Transfer rate:          10.07 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      0       1
Processing:    16  166  23.8    160     313
Waiting:       16  166  23.8    160     313
Total:         16  167  23.8    161     313

Percentage of the requests served within a certain time (ms)
  50%    161
  66%    170
  75%    178
  80%    182
  90%    198
  95%    213
  98%    236
  99%    248
 100%    313 (longest request)
```
