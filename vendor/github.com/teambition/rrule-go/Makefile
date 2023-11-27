test:
	go test --race

cover:
	rm -f *.coverprofile
	go test -coverprofile=rrule.coverprofile
	go tool cover -html=rrule.coverprofile
	rm -f *.coverprofile

.PHONY: test cover
