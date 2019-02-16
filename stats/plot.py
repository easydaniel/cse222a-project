from matplotlib import pyplot as plt
import seaborn as sns

REQLAT_PATH = './reqlat.ds'


def plot():
    dataset = [list(map(float, v.split()))
               for v in open(REQLAT_PATH, 'r').readlines()]
    frate, lat = [], []
    xlabels = list(range(10, 501, 10))
    idx = 0
    for _ in range(4):
        tf, tl = [], []
        for _ in range(10, 501, 10):
            tf.append(dataset[idx][0])
            tl.append(dataset[idx][1])
            idx += 1
        frate.append(tf)
        lat.append(tl)

    for i in range(4):
        sns.lineplot(x=xlabels, y=lat[i], label=f'{i}')

    plt.show()


if __name__ == "__main__":
    plot()
