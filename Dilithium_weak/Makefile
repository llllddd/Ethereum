CC = /usr/bin/cc
CFLAGS = -Wall -Wextra -march=native -mtune=native -O3 -fomit-frame-pointer
NISTFLAGS = -O3 -fomit-frame-pointer -march=native -fPIC
SOURCES = sign.c polyvec.c packing.c poly.c reduce.c ntt.c rounding.c fips202.c randombytes.c
HEADERS = params.h sign.h polyvec.h packing.h poly.h reduce.h ntt.h randombytes.h\
  rounding.h fips202.h

all: example1

example1: $(SOURCES) example1.c $(HEADERS)
	$(CC) $(SOURCES) example1.c -o example1

