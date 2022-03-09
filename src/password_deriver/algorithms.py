import string
import argparse
import subprocess
import getpass
import hashlib
import random
import pathlib
import json
import sys

from . import base58


def derive_secret_key(master_password, domain, username, counter):
    salt = f"{domain}{username}{counter:x}".encode("utf-8")
    return hashlib.pbkdf2_hmac(
        "sha256", master_password, salt=salt, iterations=100_000, dklen=32
    )


def format_password_hex(secret_key):
    return base58.b58encode(secret_key[-12:]).decode("ascii") + "-"


def format_ethereum_private_key(secret_key):
    return secret_key.hex()


def format_password_lesspass(secret_key):
    main_alphabet = (
        string.ascii_lowercase
        + string.ascii_uppercase
        + string.digits
        + string.punctuation
    )
    seed = int.from_bytes(secret_key, byteorder="big")

    password = []
    for _ in range(12):
        seed, index = divmod(seed, len(main_alphabet))
        password.append(main_alphabet[index])

    extra_chars = []
    for alphabet in [
        string.ascii_lowercase,
        string.ascii_uppercase,
        string.digits,
        string.punctuation,
    ]:
        seed, index = divmod(seed, len(alphabet))
        extra_chars.append(alphabet[index])

    for extra_char in extra_chars:
        seed, index = divmod(seed, len(password))
        password.insert(index, extra_char)

    return "".join(password)
