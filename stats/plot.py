from matplotlib import pyplot as plt
import seaborn as sns
import numpy as np

REQLAT_PATH = './50c1p15t'


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


if __name__ == "__main__":
    for i in range(3):
        plot(i)
