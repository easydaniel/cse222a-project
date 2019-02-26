from django.db import models
from storages.backends.s3boto3 import S3Boto3Storage

# Create your models here.

class FileStorage(S3Boto3Storage):
    location = 'file'
    default_acl = 'public-read'
    file_overwrite = False

class Client(models.Model):
    Name = models.CharField(max_length=200, null=False)
    Age = models.PositiveSmallIntegerField(null=False)

class Upload(models.Model):
    name = models.CharField(max_length=200, null=False)
    file = models.FileField(storage=FileStorage(), null=False)



