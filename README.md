## Tracer 
a hero in Blizzard Game 《OverWatch》, she can use pulse bombs and spatio-temporal shuttles. 
this project is a go package to upload with kafka and http api. you can accept this message in other project and resolve it by your self 

#### install 

```bash
# govendor
govender fetch github.com/MartinFang/Tracer

# go get

go get github.com/MartinFang/Tracer
```

#### using

```go

package main
import tracer

func main(){
    address := "localhost:9092,localhost:9093,localhost:9094" // the kafka brokers address list, split whit ','
    topic := "test_topic" // producer upload message to this topic
    url := "localhost:8080" // url for uoload
    tracer.Init(address, topic, url) // init config

    msg := &Msg{"test", map[string]string{"test": "test"}, map[string]interface{}{"test": "test", "test2": 123}} // the struct for meesage 
    msg.TracerWithKafka() // upload with fakfa
    msg.TracerWithUrl() // upload with url
}
```