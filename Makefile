TARG=iexprename
GOFILES=iexprename.go

$(TARG): $(GOFILES)
	go build -o $(TARG) $(GOFILES)

clean:
	rm -f $(TARG) mon.out
