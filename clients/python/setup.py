from setuptools import setup

with open("README.md", "r") as fh:
    long_description = fh.read()

setup(
    name="inspr",
    version="0.0.1",
    author="Inspr LLC",
    description="This module define APIs to interact with the Inspr environment",
    long_description=long_description,      # Long description read from the the readme file
    long_description_content_type="text/markdown",
    python_requires='>=3.6',
    packages=[
        "inspr",
    ],
    requires=[
        "requests",
        "flask"
    ],
    install_requires=[
        "requests",
        "flask"
    ]
)
