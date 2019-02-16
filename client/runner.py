from subprocess import call
from time import sleep

for T in range(4):
    for C in range(10, 501, 10):
        print(f'1000 {C} {T}')
        call(
            f'./benchmark -n 10000 -c {C} -t {T} http://dev-env.jvrb83sqrz.us-west-2.elasticbeanstalk.com/main/time ', shell=True)
        sleep(30)
