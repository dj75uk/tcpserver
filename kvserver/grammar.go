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
	}
}
