
default: test

test: 
	go test -v ./parser
	go test -v ./tokenizer

# GOPROXY=proxy.golang.org go list -m github.com/isaac-weisberg/go-jason@v0.1.1