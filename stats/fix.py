lines = [list(map(float, line.split()))
         for line in open('./reqlat.ds.1', 'r').readlines()]

Ts = range(4)
Cs = range(20, 501, 20)

with open('./reqlat', 'w') as fo:
    for i in range(0, len(lines), len(Cs)):
        for j in range(len(Cs)):
            f, l = lines[i + j]
            fo.write(f'{(f / Cs[j]):.6f} {l:.6f}\n')
