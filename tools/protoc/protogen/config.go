package main

type Config struct {
	proto, genpath, dproto                                                          string
	include                                                                         string
	useEnumPlugin, useGateWayPlugin, useValidatorsOutPlugin, useGqlPlugin, stdPatch bool
}

var config = Config{}
