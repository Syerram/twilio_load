##Twilio Endpoint Load tester

The goal of the script is to load test your twilio endpoint, just like a simple `ab` tester.
You can pass json data to the endpoint (both valid and invalid) and see how it utilizes resources and get an **estimate of costs** of using twilio with high load. 
For e.g. You can send 5000 invalid requests and see how your server utilizes its resources (does it flood the log files, emails, database, costs of receiving/replying etc).


In order to run the script, make a copy of the sample ini files and then run

```
go build tw_load.go
./tw_load -f ./tw_load.ini -d ./tw_data.ini
```

Once completed, it will print the results of the load testing along with costs of sending/receiving messages.

###Config INI
Contains information such as account token, load factors
###Data INI
Contains data that needs to be sent to the endpoint. 

###TODO
- Remove the data ini file and instead choose a simple file format (txt) by having each line as a json data row. We can then distribute data across the requests