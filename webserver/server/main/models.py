from django.db import models

# Create your models here.


class Client(models.Model):
    Name = models.CharField(max_length=200, null=False)
    Age = models.PositiveSmallIntegerField(null=False)
