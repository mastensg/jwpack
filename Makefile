GOSRC = http.go main.go pack.go

all: jwpack

gen_bindata.go: Makefile gen/bindata.go bindata/*
	go run gen/bindata.go

jwpack: Makefile $(GOSRC) gen_bindata.go
	go build

.PHONY: clean dist sources

clean:
	$(RM) jwpack gen_bindata.go

sources:
	@ls -1d Makefile
	@ls -1d $(GOSRC)
	@ls -1d gen/bindata.go
	@find bindata/
