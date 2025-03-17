all:
	go build -o main

test:
	bisect FEAT1=PATTERN FEAT2=PATTERN FEAT3=PATTERN ./main

clean: rm -rf main

.PHONY: clean test
