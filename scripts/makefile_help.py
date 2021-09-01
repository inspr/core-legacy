from sys import argv
import re


def print_helper(helperh, commandd, maxlen):
    helperh = "".join([(maxlen-len(commandd)) * " " + "\t\t-- "+helperh.split(
        "\n")[0]]+list(map(lambda x: maxlen*" " + "\t\t"+x, helperh.split("\n")[1:])))
    print(commandd, helperh)
    helperh = ""


title = re.compile("# *[a-zA-Z]+ *{")
endtitle = re.compile("# *}")

with open(argv[1]) as f:
    l = []
    command = ""
    helper = ""
    level = 1
    titles = []
    hs = 0
    for line in f:
        if title.match(line):
            titles.append((line.strip()[2:-1], level, hs))
            level += 1
        if endtitle.match(line):
            level -= 1
        if line[:2] == '##':
            if len(argv) >= 3 or helper == "":
                helper += line[2:].strip(" ")
            continue
        if helper != "":
            helper = helper.strip()
            command = line[:line.find(":")]
            l.append((helper, command))
            helper = ""
            command = ""
            hs += 1

    if len(argv) > 2:
        for i, (h, c) in enumerate(l):
            if c == argv[2]:
                print(c)
                print(h)

    else:
        m = len(max(l, key=lambda x: len(x[1]))[1])
        n = len(max(l, key=lambda x: len(x[0]))[0])
        k = 0
        for i, (h, c) in enumerate(l):
            while k < len(titles) and titles[k][2] <= i:
                if (titles[k][1] == 1):
                    print("\n" + (m+n)*"=")
                    print("\n" + (titles[k][0]).center(m + n) + "\n")

                else:
                    print("\n" + titles[k][1] * "#" + " " + titles[k][0])
                k += 1
            print_helper(h, c, m)
