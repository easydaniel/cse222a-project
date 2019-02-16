from subprocess import call
from time import sleep

for T in range(4):
    for C in range(20, 501, 20):
        call(
            f'./benchmark -n 30 -c {C} -t {T} http://dev-env.jvrb83sqrz.us-west-2.elasticbeanstalk.com/main/time ', shell=True)
        sleep(5)
