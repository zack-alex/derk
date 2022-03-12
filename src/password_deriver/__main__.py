import string
import argparse
import subprocess
import getpass
import hashlib
import random
import pathlib
import json
import sys

from . import config, algorithms


FINGERPRINT_CHARS = [" ", "░", "▒", "▓", "█"]


if sys.platform.startswith("linux"):
    def write_to_clipboard(data):
        process = subprocess.Popen("wl-copy", stdin=subprocess.PIPE)
        process.communicate(data.encode("ascii"))
elif sys.platform.startswith("darwin"):
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
    return hashlib.scrypt(
        password, salt=salt, n=2 ** 15, r=8, p=1, maxmem=64 * 1024 ** 2
    )


def get_master_password():
    salt = config.get_or_init_salt().encode("ascii")
    master_password = getpass.getpass("Enter the master passphrase: ").encode("utf-8")

    print("Master passphrase fingerprint:")
    print(graphical_fingerprint(password_hash(master_password, salt)))

    return master_password


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("domain")
    parser.add_argument("--user")
    parser.add_argument("--counter", type=int, default=1)
    parser.add_argument("--type", choices=["password", "ethereum"], default="password")
    args = parser.parse_args()
    args.user = args.user or config.get_user()
    print("Domain:", args.domain)
    print("User:", args.user)
    fn = {"password": algorithms.format_password_hex, "ethereum": algorithms.format_ethereum_private_key}[args.type]
    master_password = get_master_password()

    secret_key = algorithms.derive_secret_key(
        master_password, args.domain, args.user, args.counter
    )
    password = fn(secret_key)

    write_to_clipboard(password)
    print("The password is copied to the clipboard")


if __name__ == "__main__":
    main()
