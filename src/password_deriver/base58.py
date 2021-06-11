DIGIT_REPRESENTATION = b"123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"


def b58encode(s):
    leading_zeros = 0
    for byte in s:
        if byte != 0:
            break
        leading_zeros += 1

    x = int.from_bytes(s, byteorder="big")

    result = []
    while x > 0:
        x, digit = divmod(x, 58)
        result.insert(0, DIGIT_REPRESENTATION[digit])

    for _ in range(leading_zeros):
        result.insert(0, DIGIT_REPRESENTATION[0])

    return bytes(result)
