package Pool

import (
	"fmt"
	"math/rand"
)

type Job struct {
	Id      int //id
	RandNum int //需要计算的随机数
}

type Result struct {
	job *Job //对象实例
	sum int  //求和
}

func Demo() {

	//需要2个管道
	//1.job管道
	jobChan:=make(chan *Job,128)
	//2.结果管道
	resultChan:=make(chan *Result,128)
	//3.创建工作池
	createPool(64,jobChan,resultChan)
	//4.开个打印的协程
	go func(resultChan chan *Result) {
		//遍历结果管道打印
		for result:=range resultChan{
			fmt.Printf("job id:%v randnum:%v result:%d\n",result.job.Id,result.job.RandNum,result.sum)
		}
	}(resultChan)

	var id int
	//循环创建job，输入到管道
	for {
		id++
		//生成随机数
		rNum:=rand.Int()
		job:=&Job{
			Id:id,
			RandNum: rNum,
		}
		jobChan<-job
	}
}

//创建工作池
//num开启协诚的个数
func createPool(num int, jobChan chan *Job, resultChan chan *Result) {
	//根据开协诚的个数，运行程序
	for i := 0; i < num; i++ {
		go func(jobChan chan *Job, resultChan chan *Result) {
			//执行运算
			//遍历Job管道所有数据并相加
			for job := range jobChan {
				//接收随机数
				rNum := job.RandNum
				//随机数每一位相加
				//定义返回值
				var sum int
				for rNum != 0 {
					temp := rNum % 10
					sum += temp
					rNum /= 10
				}

				//想要的结果Result
				r := &Result{
					job: job,
					sum: sum,
				}

				//运算结果扔到管道
				resultChan <- r
			}
		}(jobChan, resultChan)
	}
}
