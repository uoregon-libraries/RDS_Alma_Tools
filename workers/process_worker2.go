package main

import(
  worker "github.com/contribsys/faktory_worker_go"
  "rds_alma_tools/withdraw"
)

func main(){
  mgr := worker.NewManager()
  mgr.Concurrency = 1
  mgr.ProcessStrictPriorityQueues("process2")
  mgr.Register("HelloJob", withdraw.HelloJobWorker)
  mgr.Register("ProcessJob", withdraw.CheckJobWorker)
  mgr.Run()
}
