from matplotlib import pyplot as plt
import seaborn as sns
import numpy as np
import math

REQLAT_PATH = './256c1p15t'


def plot(tp=0):
    sns.set_style('whitegrid')
    plt.figure(figsize=(12, 6))
    plt.xlabel('# of connections')
    ylabel = ['Failure rate', 'Latency (ms)', 'Throughput(req/s)']
    plt.ylabel(ylabel[tp])

    dataset = [list(map(float, v.split()))
               for v in open(REQLAT_PATH, 'r').readlines()]
    dataset = np.array(dataset)

    xticks = list(range(20, 501, 20))
    plt.xticks(xticks)

    labels = ['Normal Requests (0ms)', 'Short Requests (100ms)',
              'Long Requests (500ms)', 'Mixed']

    for i in range(4):
        sns.lineplot(x=xticks, y=dataset[:, tp]
                     [i*len(xticks): (i+1)*len(xticks)], label=labels[i])
    fname = ['failrate.png', 'reqlatency.png', 'throughput.png']
    plt.savefig(f'{REQLAT_PATH[2:]}-{fname[tp]}')


def perf(tp=0, sch='rr'):
    sns.set_style('whitegrid')

    plt.figure(figsize=(12, 6))
    plt.xlabel('# of connections')
    ylabel = ['Failure rate', 'Latency (ms)', 'Throughput(req/s)']
    plt.ylabel(ylabel[tp])

    bad = [list(map(float, v.split()))
           for v in open(f'{sch}-bad', 'r').readlines()]
    bad = np.array(bad)
    normal = [list(map(float, v.split()))
              for v in open(f'{sch}-normal', 'r').readlines()]
    normal = np.array(normal)

    xticks = list(range(20, 501, 20))
    plt.xticks(xticks)
    sns.lineplot(x=xticks, y=bad[:, tp], marker="o", label='normal')
    sns.lineplot(x=xticks, y=normal[:, tp], marker="o", label='good')
    fname = ['failrate.png', 'reqlatency.png', 'throughput.png']
    plt.savefig(f'{sch}-{fname[tp]}')


def lllt(tp=0):
    sns.set_style('whitegrid')

    plt.figure(figsize=(12, 6))
    plt.xlabel('# of connections')
    ylabel = ['Failure rate', 'Latency (ms)', 'Throughput(req/s)']
    plt.ylabel(ylabel[tp])

    xticks = list(range(20, 501, 20))
    plt.xticks(xticks)

    for t in ['ll', 'lt']:
        d = [list(map(float, v.split()))
             for v in open(f'{t}-normal', 'r').readlines()]
        d = np.array(d)
        sns.lineplot(x=xticks, y=d[:, tp], marker="o", label=f'{t}')
    fname = ['failrate.png', 'reqlatency.png', 'throughput.png']
    plt.savefig(f'lt-ll-{fname[tp]}')


def compare(tp=0):
    sns.set_style('whitegrid')

    plt.figure(figsize=(12, 6))
    plt.xlabel('# of connections')
    ylabel = ['Failure rate', 'Latency (ms)', 'Throughput(req/s)']
    plt.ylabel(ylabel[tp])

    xticks = list(range(20, 501, 20))
    plt.xticks(xticks)
    for t, name in zip(['rr', 'll', 'pa'], ['round-robin', 'least-load', 'pairing']):
        d = [list(map(float, v.split()))
             for v in open(f'{t}-epbad', 'r').readlines()]
        d = np.array(d)
        sns.lineplot(x=xticks, y=d[:, tp], marker="o", label=name)
    fname = ['failrate.png', 'reqlatency.png', 'throughput.png']
    plt.savefig(f'comp-{fname[tp]}')


def wrr(tp=0):
    sns.set_style('whitegrid')

    plt.figure(figsize=(12, 6))
    plt.xlabel('# of connections')
    ylabel = ['Failure rate', 'Latency (ms)', 'Throughput(req/s)']
    plt.ylabel(ylabel[tp])

    xticks = list(range(20, 501, 20))
    plt.xticks(xticks)
    for r in [11, 15, 125]:
        d = [list(map(float, v.split()))
             for v in open(f'wrr-pattern{r}', 'r').readlines()]
        d = np.array(d)
        sns.lineplot(x=xticks, y=d[:, tp], marker="o",
                     label=f'weighted-rr {r}')
    fname = ['failrate.png', 'reqlatency.png', 'throughput.png']
    plt.savefig(f'wrr-pattern-{fname[tp]}')


def wrrp(tp=0):
    sns.set_style('whitegrid')

    plt.figure(figsize=(12, 6))
    plt.xlabel('# of connections')
    ylabel = ['Failure rate', 'Latency (ms)', 'Throughput(req/s)']
    plt.ylabel(ylabel[tp])

    xticks = list(range(20, 501, 20))
    plt.xticks(xticks)
    for r in [11, 23, 15]:
        d = [list(map(float, v.split()))
             for v in open(f'wrrp-normal{r}', 'r').readlines()]
        d = np.array(d)
        sns.lineplot(x=xticks, y=d[:, tp], marker="o",
                     label=f'weighted-rr {r}')
    fname = ['failrate.png', 'reqlatency.png', 'throughput.png']
    plt.savefig(f'wrrp-{fname[tp]}')


def websites():
    sns.set_style('whitegrid')

    plt.figure(figsize=(12, 6))
    plt.ylabel('# of websites')
    d = [list(map(int, v.split()))
         for v in open(f'websites', 'r').readlines()]
    d = np.array(d)
    plt.xticks(d[:, 0])
    plt.yticks([i * 500000000 for i in range(5)])
    plt.ylim(0, 500000000 * 4)
    ax = sns.barplot(x=d[:, 0], y=d[:, 1], color='steelblue')
    ax.set_xticklabels(ax.get_xticklabels(), rotation=45)
    ax.set_yticklabels(['0', '500,000,000', '1,000,000,000',
                        '1,500,000,000', '2,000,000,000'])
    plt.savefig(f'websites.png')


if __name__ == "__main__":
    # websites()
    for i in range(3):
        # lllt(i)
        wrr(i)
        # plot(i)
        # perf(i, sch='pa')
        # compare(i)
