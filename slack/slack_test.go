package slack

import "testing"

func TestPostMsg(t *testing.T) {
	Configure("xoxp-219053016240-220597711318-219638448467-01f44c58422d064107673735b2f4650f", "clear")
	PostMessage("hello, I'm Lnk")
}

