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


def derive_and_format(master_password, spec):
    domain = spec["domain"]
    username = spec["username"]
    method = spec["method"]
    if method == "v1-count3" or method == "v1-shorter-count3":
        counter = 3
    elif (
        method == "v1-count2"
        or method == "v1-shorter-count2"
        or method == "v1-with-bang-count2"
    ):
        counter = 2
    else:
        counter = 1
    secret_key = derive_secret_key(master_password, domain, username, counter)
    if method == "v1":
        return format_password_hex(secret_key)
    elif method == "v1-count2" or method == "v1-count3":
        return format_password_hex(secret_key)
    elif method == "v1-wo-tail":
        return format_password_hex(secret_key)[:-1]
    elif method == "v1-with-bang" or method == "v1-with-bang-count2":
        return format_password_hex(secret_key)[:-1] + "!"
    elif (
        method == "v1-shorter"
        or method == "v1-shorter-count2"
        or method == "v1-shorter-count3"
    ):
        return format_password_hex(secret_key)[:-3]
    elif method == "v1-shorter-with-dash":
        return format_password_hex(secret_key)[:-3] + "-"
    elif method == "ethereum":
        return format_ethereum_private_key(secret_key)
    elif method == "none":
        return ""
    else:
        raise ValueError(f"Unknown method: {method}.")
