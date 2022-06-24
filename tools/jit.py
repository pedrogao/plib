import ctypes
import dis

import mj


def get_codeobj(function):
    return function.__code__


# SSA
# operation dest src


class Compiler(object):
    """Compiles Python bytecode to intermediate representation (IR)."""

    def __init__(self, bytecode, constants):
        self.bytecode = bytecode
        self.constants = constants
        self.index = 0

    def fetch(self):
        byte = self.bytecode[self.index]
        self.index += 1
        return byte

    def decode(self):
        opcode = self.fetch()
        opname = dis.opname[opcode]
        # 单指令
        if opname.startswith(("UNARY", "BINARY", "INPLACE", "RETURN")):
            argument = None
            self.fetch()
        else:
            argument = self.fetch()

        return opname, argument

    def variable(self, number):
        # AMD64 argument passing order for our purposes.
        order = ("rdi", "rsi", "rdx", "rcx")
        return order[number]

    def compile(self):
        while self.index < len(self.bytecode):
            op, arg = self.decode()  # 指令、参数
            print(op, arg)

            if op == "LOAD_FAST":
                yield "push", self.variable(arg), None

            elif op == "STORE_FAST":
                yield "pop", "rax", None
                yield "mov", self.variable(arg), "rax"

            elif op == "LOAD_CONST":
                yield "immediate", "rax", self.constants[arg]
                yield "push", "rax", None

            elif op == "BINARY_MULTIPLY":
                yield "pop", "rax", None
                yield "pop", "rbx", None
                yield "imul", "rax", "rbx"
                yield "push", "rax", None

            elif op in ("BINARY_ADD", "INPLACE_ADD"):
                yield "pop", "rax", None
                yield "pop", "rbx", None
                yield "add", "rax", "rbx"
                yield "push", "rax", None

            elif op in ("BINARY_SUBTRACT", "INPLACE_SUBTRACT"):
                yield "pop", "rbx", None
                yield "pop", "rax", None
                yield "sub", "rax", "rbx"
                yield "push", "rax", None

            elif op == "UNARY_NEGATIVE":
                yield "pop", "rax", None
                yield "neg", "rax", None
                yield "push", "rax", None

            elif op == "RETURN_VALUE":
                yield "pop", "rax", None
                yield "ret", None, None

            else:
                raise NotImplementedError(op)


def optimize(ir):
    def fetch(n):
        if n < len(ir):
            return ir[n]
        else:
            return None, None, None

    index = 0
    while index < len(ir):
        op1, a1, b1 = fetch(index)
        op2, a2, b2 = fetch(index + 1)
        op3, a3, b3 = fetch(index + 2)

        if op1 == "mov" and a1 == b1:
            index += 1
            continue

        if op1 == op2 == "mov" and a1 == b2:
            index += 2
            yield "mov", a2, b1
            continue

        if op1 == "push" and op2 == "pop":
            index += 2
            yield "mov", a2, a1
            continue

        if op1 == "push" and op3 == "pop" and op2 not in ("push", "pop"):
            if a2 != a3:
                index += 3
                yield "mov", a3, a1
                yield op2, a2, b2
                continue

        index += 1
        yield op1, a1, b1


class Assembler(object):
    """An x86-64 assembler."""

    def __init__(self, size):
        self.block = mj.create_block(size)
        self.index = 0
        self.size = size

    @property
    def raw(self):
        """Returns machine code as a raw string."""
        return bytes(self.block[:self.index])

    @property
    def address(self):
        """Returns address of block in memory."""
        return ctypes.cast(self.block, ctypes.c_void_p).value

    def little_endian(self, n):
        """Converts 64-bit number to little-endian format."""
        if n is None:
            n = 0
        return [(n & (0xff << (i * 8))) >> (i * 8) for i in range(8)]

    def registers(self, a, b=None):
        """Encodes one or two registers for machine code instructions."""
        order = ("rax", "rcx", "rdx", "rbx", "rsp", "rbp", "rsi", "rdi")
        enc = order.index(a)
        if b is not None:
            enc = enc << 3 | order.index(b)
        return enc

    def emit(self, *args):
        """Writes machine code to memory block."""
        print("emit: ", args)
        for code in args:
            self.block[self.index] = code
            self.index += 1

    def ret(self, a, b):
        self.emit(0xc3)

    def push(self, a, _):
        self.emit(0x50 | self.registers(a))

    def pop(self, a, _):
        self.emit(0x58 | self.registers(a))

    def imul(self, a, b):
        self.emit(0x48, 0x0f, 0xaf, 0xc0 | self.registers(a, b))

    def add(self, a, b):
        self.emit(0x48, 0x01, 0xc0 | self.registers(b, a))

    def sub(self, a, b):
        self.emit(0x48, 0x29, 0xc0 | self.registers(b, a))

    def neg(self, a, _):
        self.emit(0x48, 0xf7, 0xd8 | self.registers(a))

    def mov(self, a, b):
        self.emit(0x48, 0x89, 0xc0 | self.registers(b, a))

    def immediate(self, a, number):
        self.emit(0x48, 0xb8 | self.registers(a), *self.little_endian(number))


def foo(x, y):
    return x * x - y * y


if __name__ == "__main__":
    bytecode = foo.__code__.co_code
    constants = foo.__code__.co_consts

    c = Compiler(bytecode, constants)
    ir = c.compile()
    ir = list(ir)
    print(ir)
    ir = list(optimize(ir))
    # print(ir)
    ir = list(optimize(ir))
    # print(ir)
    ir = list(optimize(ir))
    print(ir)

    assembler = Assembler(mj.PAGESIZE)
    for name, a, b in ir:
        emit = getattr(assembler, name)
        emit(a, b)
    print("-----" * 10)
    print(assembler.index)
    print(assembler.raw)
    mj.make_executable(assembler.block, assembler.size)
    argcount = foo.__code__.co_argcount
    signature = ctypes.CFUNCTYPE(*[ctypes.c_int64] * argcount)
    signature.restype = ctypes.c_int64

    native_foo = signature(assembler.address)
    print(native_foo(2, 3))
