from os import getenv

if __name__ == "__main__":
    version = getenv("VERSION")
    kind = getenv("KIND")
    major, minor, patch = version.split(".")
    major = int(major)
    minor = int(minor)
    patch = int(patch)

    if kind == "major":
        major += 1
        minor = 0
        patch = 0
    elif kind == "minor":
        minor += 1
        patch = 0
    elif kind == "patch":
        patch += 1

    print(f"{major}.{minor}.{patch}", end='')
