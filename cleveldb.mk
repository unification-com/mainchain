CLEVELDB_ENABLED ?= false

ifeq ($(CLEVELDB_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for cleveldb support, please install or set CLEVELDB_ENABLED=false)
    else
      build_tags += cleveldb
      build_tags += gcc
      ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
    endif
  else
    build_tags += cleveldb
    build_tags += gcc
    ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
  endif
endif
