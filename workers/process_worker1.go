package main

import(
  worker "github.com/contribsys/faktory_worker_go"
  "rds_alma_tools/withdraw"
)

func main(){
  mgr := worker.NewManager()
  mgr.Concurrency = 1
  mgr.Register("HelloJob", withdraw.HelloJobWorker)
  mgr.Register("ProcessJob", withdraw.ProcessWorker)
  mgr.Register("VerifyJob", withdraw.VerifyWorker)
  mgr.ProcessStrictPriorityQueues("process1")
  mgr.Run()
}
