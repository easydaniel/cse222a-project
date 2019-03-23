from django.shortcuts import render, HttpResponse
import time
from random import randint
from .models import Client
import requests
from django.core.files.storage import FileSystemStorage
import uuid, os
from django.test import Client
from .models import Upload
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


def populateData(request):
    amount = request.GET.get("amount", "100")
    age = request.GET.get("age", "15")
    ClientList = []
    for i in range(int(amount)):
        name = "FakeName" + str(randint(0, 100))
        # age = randint(10, 30)
        ClientList.append((Client(Name=name, Age=age)))

    Client.objects.bulk_create(ClientList)
    return HttpResponse("Populate Done")


def removeData(request):
    Client.objects.all().delete()
    return HttpResponse("Remove Done")


def queryData(request):
    age = request.GET.get("age", "20")
    amount = Client.objects.filter(Age__lte=age).count()
    return HttpResponse("Query Done. Age(Smaller): " + str(age) + ", Amount: " + str(amount))


def uploadFile(request):
    size = request.GET.get("size", "100")
    filename = str(uuid.uuid4())
    c = Client()
    with open(filename, "wb") as fp:
        fp.write(os.urandom(int(size)))
    fp.close()
    with open(filename, "rb") as fp:
        # print(c.get("/main/uploadFileHandler"))
        response = c.post("/main/uploadFileHandler", {"name": filename, "fileData": fp})
    fp.close()
    os.remove(fp.name)
    print(response)
    return HttpResponse("Upload Done. Size: " + size)


def uploadFileHandler(request):
    if request.method == "POST":
        name = request.POST["name"]
        fileData = request.FILES["fileData"]
        upload = Upload(name=name, file=fileData)
        upload.save()
    elif request.method == "GET":
        print("GET")
    return HttpResponse("Return uploadFileHandler")


def pingOther(request):
    ip = request.GET.get("ip")
    result = requests.get("http://" + ip + "/main/")
    return HttpResponse(result)


def notServerOtherPeriod(request):
    host = request.get_host()
    server = request.GET.get("server")
    period_s = int(request.GET.get('period_s', '5')) / 1000
    period_not_s = int(request.GET.get('period_not_s', '5')) / 1000
    if host == server:
        time.sleep(float(period_s))
        return HttpResponse("Sleep done(s). Time: " + time.ctime() + " Period = " + str(period_s))
    else:
        time.sleep(float(period_not_s))
        return HttpResponse("Sleep done(not_s). Time: " + time.ctime() + " Period = " + str(period_not_s))
