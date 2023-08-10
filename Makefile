
default: test

test: 
	go test -v ./parser
	go test -v ./tokenizer
