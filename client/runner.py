from subprocess import call
from time import sleep

RR_EP = 'http://54.200.39.211:3019' 
LL_EP = 'http://52.41.245.156:3019'
PA_EP = 'http://52.42.191.169:3019'
SERV_EP = 'http://dev-env-easy.jvrb83sqrz.us-west-2.elasticbeanstalk.com/main/time'

EP1 = '/main/time'
EP2 = '/main/notServerOtherPeriod'

for T in range(4):
    for C in range(20, 501, 20):
        call(
            './benchmark -n 30 -c %d -t %d -s %d %s' % (C, 3, 2, LL_EP + EP2), shell=True)
        sleep(5)
    break
