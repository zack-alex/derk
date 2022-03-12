import pathlib
import json
import secrets


def config_path():
    return pathlib.Path.home() / ".config" / "password-deriver" / "config.json"


def get_salt():
    with open(config_path()) as f:
        config = json.load(f)
    return config["salt"]


def set_salt(salt):
    try:
        with open(config_path()) as f:
            config = json.load(f)
    except FileNotFoundError:
        config = {}
    config["salt"] = salt
    config_path().parent.mkdir(parents=True, exist_ok=True)
    with open(config_path(), "w") as f:
        json.dump(config, f)


def init_salt():
    return secrets.token_bytes(16).hex()


def get_or_init_salt():
    try:
        return get_salt()
    except (FileNotFoundError, KeyError):
        salt = init_salt()
        set_salt(salt)
        return salt


def get_user():
    with open(config_path()) as f:
        config = json.load(f)
    return config["user"]


def get_master_password_hash():
    with open(config_path()) as f:
        config = json.load(f)
    res = config.get("master_password_hash")
    if res is None:
        return None
    return bytes.fromhex(res)


def set_master_password_hash(h):
    h = h.hex()
    try:
        with open(config_path()) as f:
            config = json.load(f)
    except FileNotFoundError:
        config = {}
    config["master_password_hash"] = h
    config_path().parent.mkdir(parents=True, exist_ok=True)
    with open(config_path(), "w") as f:
        json.dump(config, f)
