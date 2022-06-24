#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import capstone

md = capstone.Cs(capstone.CS_ARCH_X86, capstone.CS_MODE_64)


def hexbytes(raw):
    return "".join("%02x " % b for b in raw)


code = bytearray.fromhex('4889f848f7d84889c34889f8480fafc3504889f34889f0480fafc34889c3584829d8c3')

for i in md.disasm(code, 0xc0000a4018):
    print("0x%x %-15s%s %s\n" % (i.address, hexbytes(i.bytes), i.mnemonic, i.op_str))
    # if i.mnemonic == "ret":
    #     break
