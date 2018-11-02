package accounts

import (
	"testing"
)

func TestURLParsing(t *testing.T) {
	url, err := parseURL("https://xchain-go.org")
	if err != nil {
		t.Errorf("解析地址出错: %v", err)
	}
	if url.Scheme != "https" {
		t.Errorf("https与%v不一致", url.Scheme)
	}
	if url.Path != "xchain-go.org" {
		t.Errorf("xchain-go.org与%v不同", url.Path)
	}
	_, err = parseURL("xchain-go.org")
	if err == nil {
		t.Errorf("未发生错误提示")
	}
}

func TestString(t *testing.T) {
	url := URL{Scheme: "Wallet", Path: "workspace/src"}
	if url.String() != "Wallet://workspace/src" {
		t.Errorf("解析的结果不正确")
	}
	url = URL{Scheme: "", Path: "workspace/src"}
	if url.String() != "workspace/src" {
		t.Errorf("解析结果不正确")
	}
}

func TestMarshalJSON(t *testing.T) {
	url := URL{Scheme: "Wallet", Path: "workspace/src"}
	json, err := url.MarshalJSON()
	if err != nil {
		t.Errorf("编码成json文件出错%v", err)
	}
	if string(json) != "\"Wallet://workspace/src\"" {
		t.Errorf("json编码结果出错%v应该为\"Wallet://workspace/src\"", string(json))
	}
}
