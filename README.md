##Twilio Endpoint Load tester

The goal of the script is to load test your twilio endpoint, just like a simple `ab` tester.
You can pass json data to the endpoint (both valid and invalid) and see how it utilizes resources and get an **estimate of costs** of using twilio with high load. 
For e.g. You can send 5000 invalid requests and see how your server utilizes its resources (does it flood the log files, emails, database, costs of receiving/replying etc).


In order to run the script, make a copy of the sample ini files and then run

```
$ go build tw_load.go
$ ./tw_load -f ./tw_load.ini -d ./tw_data.ini
```

Once completed, it will print the results of the load testing along with costs of sending/receiving messages (sample response below)

Benchmarking http://127.0.0.1:5000/sms (be patient).....
Done!
Document Path:        http://127.0.0.1:5000/sms
Concurrency Level:    5
Users:                200
Completed requests:   1000
Failed requests:      0
Total time:           2.15083s
Time per request:     10.691935ms
-----------------------------------
Total recv:       1000
Total sent        1000
Approx cost:      $15.00

```

###Config INI
Contains information such as account token, load factors

###Data INI
Contains data that needs to be sent to the endpoint. 


###TODO
- Remove the data ini file and instead choose a simple file format (txt) by having each line as a json data row. We can then distribute data across the requests
- Expand unit tests to cover edge cases