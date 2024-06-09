package main

type Config struct {
	proto, genpath, dproto                                                         string
	include                                                                        string
	useEnumPlugin, useGateWayPlugin, useValidatorOutPlugin, useGqlPlugin, stdPatch bool
	apidocDir                                                                      string
}

var config = Config{}
