package kvserver

import "kvsapp/parsing"

func getStandardGrammar() map[string]parsing.ParserGrammar {
	return map[string]parsing.ParserGrammar{
		"die": {ExpectedArguments: 0},
		"bye": {ExpectedArguments: 0},
		"get": {ExpectedArguments: 1},
		"del": {ExpectedArguments: 1},
		"put": {ExpectedArguments: 2},
		"hed": {ExpectedArguments: 2, Arg2LengthIsValue: true},
		"hst": {ExpectedArguments: 2},
		"sdl": {ExpectedArguments: 1},
		"spt": {ExpectedArguments: 2},
		"sgt": {ExpectedArguments: 1},
		"chk": {ExpectedArguments: 1},
		"nop": {ExpectedArguments: 0},
	}
}
