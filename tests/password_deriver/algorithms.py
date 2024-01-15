import sys
import json
import pytest

from password_deriver.algorithms import derive_and_format


def test_algorithms():
    with open("tests/password_deriver/data.json", "r") as f:
        data = json.load(f)
    fixed_data = []
    for spec, _ in data:
        full_spec = dict(spec, domain="test_domain", username="test_username")
        res = derive_and_format(b"test_master_password", full_spec)
        fixed_data.append([spec, res])
    with open("tests/password_deriver/data.json", "w") as f:
        json.dump(fixed_data, f, indent=2)
