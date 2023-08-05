
default: test

test: 
	go test -v

# GOPROXY=proxy.golang.org go list -m github.com/isaac-weisberg/go-jason@v0.1.1