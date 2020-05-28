Interview Questions
===================

Question 1
----------

signature algorithm md5 with rsa encryption is depricated because is no longer secure. due to MD5 being weak against collisions attacker could recreate the certificate with their own common name if they can get the same MD5 hash from different inputs.

RSA key of 1024 bit is also considered to be no longer secure. As with the technology today is likely to be broken relatively quickly

Could also be that CA Certificate(s) not included in pem bundle browser may not be able to verify authenticity, though this could just be because the CA cerificates were not included in this question.

Another issue could be that Common name might not match service hostname.

I saved as pem file and viewed the attributes of the Certificate in MacOS. Then researched the technology used to create the certificate.

Question 2
----------

[Deployment Spec](./kubernetes/deployment.yaml)

Question 3
----------

`go run santander_cycles/main.go`

[Available Santander Cycles Golang](./santander_cycles/main.go)

Execute tests:

`go test -v -cover ./santander_cycles`


# Santander Cycles Service

[Available Santander Cycles Web Service Golang](./santander_cycles_service/main.go)

Execute tests:

`go test -v -cover ./santander_cycles_service`

Start a local service:

`go run ./santander_cycles_service`

Start docker local service:

`docker-compose up`
