# cmd version

A golang cmd App

Build
===

```
go build cmd.go
```


Usage
===

```
./cmd --help 
-Usage of ./cmd:
  -items string
        A string which specifies a list of shopping items. (default "AP1,AP1,OM1,AP1")
  -set string
        A string which specifies a list of product and price paires. (default "CH1,3.11,AP1,6.00,CF1,11.23,MK1,4.75,OM1,3.69")


./cmd -set="AP1,6.00,CF1,11.23" -items="AP1,CF1"
```

Docker Image Install
===

```
1. 

2. docker build -t my-cmd .
```

Run app in Docker
===

```
docker run -it my-cmd

```


Modify the docker CMD option
===
1. Open the dockerfile.
2. Modify the line of CMD, replace ""-set=CF1,1.33", "-items=CF1,CF1"" with other values like ""-set=AP1,6.00,MK1,4.75", "-items=OM1,CF1,CH1"".
3. Build the docker image: docker build -t my-cmd . 
4. Run the docker image: docker run -it my-cmd
