import pytest

from password_deriver.base58 import b58encode


@pytest.mark.parametrize(
    ("given", "expect"),
    [
        (b"Hello World!", b"2NEpo7TZRRrLZSi2U"),
        (b"", b""),
        (b"\x00", b"1"),
        (b"\x00\x00\x00\x00\x00", b"11111"),
        (b"\x00\x00\x01", b"112"),
        (b"\x01\x00\x00", b"LUw"),
    ],
)
def test_b58encode(given, expect):
    assert b58encode(given) == expect
