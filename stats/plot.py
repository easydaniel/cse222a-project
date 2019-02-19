from matplotlib import pyplot as plt
import seaborn as sns

REQLAT_PATH = './reqlat'


def plot():
    sns.set_style('whitegrid')
    plt.figure(figsize=(12, 6))
    plt.xlabel('# of connections')
    plt.ylabel('Failure rate')
    # plt.ylabel('Latency (ms)')
    dataset = [list(map(float, v.split()))
               for v in open(REQLAT_PATH, 'r').readlines()]
    frate, lat = [], []
    xticks = list(range(20, 501, 20))
    plt.xticks(xticks)

    idx = 0
    for _ in range(4):
        tf, tl = [], []
        for _ in range(20, 501, 20):
            tf.append(dataset[idx][0])
            tl.append(dataset[idx][1])
            idx += 1
        frate.append(tf)
        lat.append(tl)
    labels = ['Normal Requests (0ms)', 'Short Requests (100ms)',
              'Long Requests (500ms)', 'Mixed']
    for i in range(4):
        sns.lineplot(x=xticks, y=frate[i], label=labels[i])

    plt.savefig('failrate.png')


if __name__ == "__main__":
    plot()
