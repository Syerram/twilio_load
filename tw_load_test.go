package main

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Log("Populate config data") 
	{
		//pass a temp file with INI data and verify the config object
		config_file := "./test_tw_load.ini"
		config := Config{}
		config.populateConfig(*config_file)
		//validate each option
		if config.Token != "TestToken" {
			t.Fatal("Not the same token as expected")
		}
		if config.Users != 5 {
			t.Fatal("Not the same number of users as expected", config.Users, "5")
		}
		if config.Concurrency != 5 {
			t.Fatal("Not the same token as expected", config.Token, "TestToken")
		}
		if config.URI != "http://localhost:5000/test" {
			t.Fatal("Not the same URL as expected", config.URI, "http://localhost:5000/test")
		}
		if config.Unit != 1 {
			t.Fatal("Not the same unit cost as expected", config.Unit, "1")
		}
	}
	
	t.Log("Load request data from data file") {
		data_file := "./test_tw_data.ini"
		request := Request{}
		request.unmarshal(*data_file, config)
		//validate each option and signature
		if len(request.Params) != 2 {
			t.Fatal("Not the same request params as expected", len(request.Params), "2")
		}
		if request.Hashed != "_hashstring_" {
			t.Fatal("Not the same hashed string as expected", len(request.Hashed), "_hashstring_")
		}
	}
	
	t.Log("Post requests") {
		config_file := "./test_tw_load.ini"
		config := Config{}
		config.populateConfig(*config_file)
		
		data_file := "./test_tw_data.ini"
		request := Request{}
		request.unmarshal(*data_file, config)
		
		resp = request.post(config)
		if resp.status != "200 OK" {
			t.Fatal("Not the same status response", resp.status, "200 OK")
		}
	}
}