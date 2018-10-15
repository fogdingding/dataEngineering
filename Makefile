# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
RSORTLIB = libraries/sort
RSORTEXE = rsort

build:
	$(GOBUILD) $(RSORTLIB)/rsort.go $(RSORTLIB)/mergeSort.go $(RSORTLIB)/externalMergeSort.go

test:
	./$(RSORTEXE)
	./$(RSORTEXE) --parallel
	./$(RSORTEXE) -k @B
	./$(RSORTEXE) -k @B --parallel
	./$(RSORTEXE) -s -r
	./$(RSORTEXE) -s -r --parallel

external:
	./$(RSORTEXE) --external /tmp/input1.rec --chunk 3 --parallel

clean:
	$(GOCLEAN)
	rm rsort
	rm dataset/tmpFile/*