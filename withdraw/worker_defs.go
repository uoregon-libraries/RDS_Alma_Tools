package withdraw

import(
  "context"
  worker "github.com/contribsys/faktory_worker_go"
  "log"
)

func ProcessWorker(ctx context.Context, args ...interface{}) error{
  help := worker.HelperFor(ctx)
  filename := args[0].(string)
  loc_type := args[1].(string)
  data := args[2].(string)
  Process(filename, loc_type, []byte(data))
  log.Printf("Job %s executed. filename: %s", help.Jid(), filename)
  return nil
}

func VerifyWorker(ctx context.Context, args ...interface{}) error{
  help := worker.HelperFor(ctx)
  filename := args[0].(string)
  data := args[1].(string)
  VerifyList(filename, []byte(data))
  log.Printf("Job %s executed. filename: %s", help.Jid(), filename)
  return nil
}
func RestartWorker(ctx context.Context, args ...interface{}) error{
  help := worker.HelperFor(ctx)
  filename := args[0].(string)
  stage := args[1].(string)
  data := args[2].(string)
  Restart(filename, stage, []byte(data))
  log.Printf("Job %s executed. filename: %s", help.Jid(), filename)
  return nil
}
