from django.shortcuts import render, HttpResponse
import time

# Create your views here.


def main(request):
    return render(request, 'main/main.html')


def timeToSleep(request):
    period = int(request.GET.get('period', '5')) / 1000
    print("Period = " + str(period))
    print("===== Start sleep =====" + time.ctime())
    time.sleep(float(period))
    print("===== End sleep =====" + time.ctime())
    print()
    print()
    return HttpResponse("Sleep done. Time: " + time.ctime() + " Period = " + str(period))
