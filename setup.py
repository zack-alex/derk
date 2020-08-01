from setuptools import find_packages
from setuptools import setup


setup(
    name="password-deriver",
    version="0.0.1",
    url="https://github.com/harius/password-deriver-py",
    packages=find_packages("src"),
    package_dir={"": "src"},
    zip_safe=False,
    python_requires=">=3.6",
    entry_points={
        "console_scripts": ["password-deriver = password_deriver.__main__:main"]
    },
)
