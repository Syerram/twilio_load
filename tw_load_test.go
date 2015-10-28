package main

import (
	"testing"
)

func TestLoad(t *testing.T) {
	t.Log("Populate config data") 
	{
		//pass a temp file with INI data and verify the config object
		config_file := "./tw_load_test.ini"
		config := Config{}
		config.populateConfig(config_file)
		//validate each option
		if config.Token != "12345566778" {
			t.Fatal("Not the same token as expected")
		}
		if config.Users != 200 {
			t.Fatal("Not the same number of users as expected", config.Users, "200")
		}
		if config.Concurrency != 5 {
			t.Fatal("Not the same token as expected", config.Token, "5")
		}

		uri := "http://127.0.0.1:5000/sms"
		if config.URI != uri {
			t.Fatal("Not the same URL as expected", config.URI, uri)
		}
		if config.Unit != 0.0075 {
			t.Fatal("Not the same unit cost as expected", config.Unit, "0.0075")
		}
	}
	
	t.Log("Given the need to test downloading content.")
	{
		config_file := "./tw_load_test.ini"
		config := Config{}
		config.populateConfig(config_file)
	
		data_file := "./tw_data_test.ini"
		request := Request{}
		request.unmarshal(data_file, config)
		//validate each option and signature
		if len(request.Params) != 2 {
			t.Fatal("Not the same request params as expected", len(request.Params), "2")
		}
		hashed := "9/sJicrE1hXxBW3H8W21KkdNo3I=" //if this fails, then recompute the hash with the params
		if request.Hashed != hashed {
			t.Fatal("Not the same hashed string as expected", len(request.Hashed), hashed)
		}
	}
	
	t.Log("Post requests") 
	{
		config_file := "./tw_load_test.ini"
		config := Config{}
		config.populateConfig(config_file)
		
		data_file := "./tw_data_test.ini"
		request := Request{}
		request.unmarshal(data_file, config)
		
		resp := request.post(config)
		if resp.Status != "200 OK" {
			t.Fatal("Not the same status response", resp.Status, "200 OK")
		}
	}
}