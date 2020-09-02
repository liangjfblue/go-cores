# compile and lib parameter
CC      := gcc
LIBS    :=
LDFLAGS :=
DEFINES :=
INCLUDE := -I./include
CFLAGS  :=
CXXFLAGS:=

# link parameter
LIB := libcal.so

#link
$(LIB):cal.o
	$(CC) -shared -o -fPIC -o $@ $^

#compile
cal.o:cal.c
	$(CC) -c -fPIC $^ -o $@

# clean
clean:
	rm -fr *.o