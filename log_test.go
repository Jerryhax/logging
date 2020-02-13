/*
@Group : Jerryhax
@Author : Jerryhax
@DateTime : 1/11/20 10:58 AM
@File : log_test
@Version : v0.0.1.0
*/

// @Description: TODO
package logging

import (
	"testing"
)

func TestLog(t *testing.T) {

	var (
		testStr = "this is a test string"
		d       = 5.5
	)

	Setup("./", "logging")
	//SetOutputLevel(Ldebug)

	Debugf("Debug: string:%s,float64:%2.2f", testStr, d)
	Debug("Debug: foo")

	Infof("Info:  string:%s,float64:%2.2f", testStr, d)
	Info("Info: foo")

	//Warnf("Warn:  string:%s,float64:%2.2f", testStr, d)
	//Warn("Warn: foo")
	//
	//Errorf("Error:  string:%s,float64:%2.2f", testStr, d)
	//Error("Error: foo")

	//SetOutputLevel(Linfo)

	Debugf("Debug:  string:%s,float64:%2.2f", testStr, d)
	Debug("Debug: foo")

	Infof("Info:  string:%s,float64:%2.2f", testStr, d)
	Info("Info: foo")

	Warnf("Warn:  string:%s,float64:%2.2f", testStr, d)
	Warn("Warn: foo")

	//Errorf("Error:  string:%s,float64:%2.2f", testStr, d)
	//Error("Error: foo")

	//Fatalf("Error:  string:%s,float64:%2.2f", testStr, d)
	//Fatal("Error: foo")
}

