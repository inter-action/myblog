# LDFLAGS, 
# 	-X main.commitHash=$(GIT_HASH), set main packge 下面的commitHash variable 的值为 git hash.
# 	-x, 删除元信息, https://zhuanlan.zhihu.com/p/26733683
# 	LDFLAGS, link flags, https://golang.org/cmd/link/

# IMPORT_PATH
# 	return the relative path to $GOPATH/src

APP_NAME      = $(shell pwd | sed 's:.*/::')
BIN           = $(GOPATH)/bin
GIT_HASH      = $(shell git rev-parse HEAD)
#todo, filter-out
GO_FILES      = $(filter-out ./server/bindata.go, $(shell find ./server  -type f -name "*.go"))
IMPORT_PATH   = $(shell pwd | sed "s|^$(GOPATH)/src/||g")
LDFLAGS       = -w -X main.commitHash=$(GIT_HASH)
ON            = $(BIN)/on
PID           = .pid
TARGET        = build/$(APP_NAME)
DEP         := $(shell command -v dep 2> /dev/null)

pre-copmile:

$(ON):
	go install $(IMPORT_PATH)/vendor/github.com/olebedev/on

clean:
	rm -rf $(TARGET)	

serve: $(ON) clean restart
	@$(ON) -m 2 $(GO_FILES) | xargs -n1 -I{} make restart || make kill

kill:
	@kill `cat $(PID)` || true

# $$是当前bash进程的pid 等同于 $BASHPID
restart: LDFLAGS += -X main.debug=true
restart: kill $(TARGET)
	@echo restart the app...
	@$(TARGET) & echo $$! > $(PID)


# $@ symbol: https://stackoverflow.com/questions/3220277/what-do-the-makefile-symbols-and-mean
$(TARGET): clean
	@go build -ldflags '$(LDFLAGS)' -o $@ $(IMPORT_PATH)/server

build: clean $(TARGET)

csbuild:
	env GOOS=linux GOARCH=amd64 $(MAKE) build


install:

ifdef GLIDE		# :todo, this ifdef directive
	@dep ensure
else
	$(warning "Skipping installation of Go dependencies: glide is not installed")
endif


# test
test:
	go test -timeout 30s ./server/...