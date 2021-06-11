import string
import argparse
import subprocess
import getpass
import hashlib
import random
import pathlib
import json

from . import base58, config


FINGERPRINT_CHARS = [" ", "░", "▒", "▓", "█"]


def derive_secret_key(master_password, domain, username, counter):
    salt = f"{domain}{username}{counter:x}".encode("utf-8")
    return hashlib.pbkdf2_hmac(
        "sha256", master_password, salt=salt, iterations=100_000, dklen=32
    )


def format_password_hex(secret_key):
    return base58.b58encode(secret_key[-12:]).decode("ascii") + "-"


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


def write_to_clipboard(data):
    process = subprocess.Popen("pbcopy", stdin=subprocess.PIPE)
    process.communicate(data.encode("ascii"))


def graphical_fingerprint(data):
    rng = random.Random(data)
    image = [[0 for _ in range(8)] for _ in range(8)]
    for _ in range(4):
        x1, y1 = (rng.randrange(8) for _ in range(2))
        x2, y2 = (rng.randrange(z, 8) for z in (x1, y1))
        for x in range(x1, x2 + 1):
            for y in range(y1, y2 + 1):
                image[x][y] += 1
    return "\n".join(
        ["".join([FINGERPRINT_CHARS[value] for value in row]) for row in image]
    )


def password_hash(password, salt):
    return hashlib.scrypt(password, salt=salt, n=2**15, r=8, p=1, maxmem=64 * 1024 ** 2)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("domain")
    parser.add_argument("--user", default="")
    parser.add_argument("--counter", type=int, default=1)
    args = parser.parse_args()

    print("Domain:", args.domain)
    print("User:", args.user)

    master_password = getpass.getpass("Enter the master passphrase: ").encode("utf-8")

    salt = config.get_or_init_salt().encode("ascii")
    print("Master passphrase fingerprint:")
    print(graphical_fingerprint(password_hash(master_password, salt)))

    secret_key = derive_secret_key(
        master_password, args.domain, args.user, args.counter
    )
    password = format_password_hex(secret_key)

    write_to_clipboard(password)
    print("The password is copied to the clipboard")


if __name__ == "__main__":
    main()
