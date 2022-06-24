import ctypes
import ctypes.util
import mmap as MMAP
import os

# Load the C standard library
libc = ctypes.CDLL(ctypes.util.find_library("c"))

# A few constants
MAP_FAILED = -1  # voidptr actually

# Set up strerror
strerror = libc.strerror
strerror.argtypes = [ctypes.c_int]
strerror.restype = ctypes.c_char_p

# Get pagesize
PAGESIZE = os.sysconf(os.sysconf_names["SC_PAGESIZE"])

# 8-bit unsigned pointer type
c_uint8_p = ctypes.POINTER(ctypes.c_uint8)

# Setup mmap
mmap = libc.mmap
mmap.argtypes = [ctypes.c_void_p,
                 ctypes.c_size_t,
                 ctypes.c_int,
                 ctypes.c_int,
                 ctypes.c_int,
                 # Below is actually off_t, which is 64-bit on macOS
                 ctypes.c_int64]
mmap.restype = c_uint8_p

# Setup munmap
munmap = libc.munmap
munmap.argtypes = [ctypes.c_void_p, ctypes.c_size_t]
munmap.restype = ctypes.c_int

# Set mprotect
mprotect = libc.mprotect
mprotect.argtypes = [ctypes.c_void_p, ctypes.c_size_t, ctypes.c_int]
mprotect.restype = ctypes.c_int


def create_block(size):
    """Allocated a block of memory using mmap."""
    ptr = mmap(0, size, MMAP.PROT_WRITE | MMAP.PROT_READ,
               MMAP.MAP_PRIVATE | MMAP.MAP_ANONYMOUS, 0, 0)

    if ptr == MAP_FAILED:
        raise RuntimeError(strerror(ctypes.get_errno()))

    return ptr


def make_executable(block, size):
    """Marks mmap'ed memory block as read-only and executable."""
    if mprotect(block, size, MMAP.PROT_READ | MMAP.PROT_EXEC) != 0:
        raise RuntimeError(strerror(ctypes.get_errno()))


def destroy_block(block, size):
    """Deallocated previously mmapped block."""
    if munmap(block, size) == -1:
        raise RuntimeError(strerror(ctypes.get_errno()))
    del block
