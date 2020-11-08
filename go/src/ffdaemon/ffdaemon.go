package ffdaemon
import(
    "os"
    "os/exec"
    "fmt"
    "syscall"
)

var depth_value int = 0

func NewSysProcAttr() *syscall.SysProcAttr {
    return &syscall.SysProcAttr{
        HideWindow: true,
    }
}

func Daemon()(*exec.Cmd){
    env_name := "ENVNAME"
    env_value := "ENVVALUE"

    val := os.Getenv(env_name)
    if val == env_value {
        fmt.Println("child:", os.Getpid())
        return nil
    }

    cmd := &exec.Cmd{
            Path: os.Args[0],
            Args: os.Args,
            Env:  os.Environ(),
            Stderr: os.Stderr,
            Stdout: os.Stdout,
		    SysProcAttr: NewSysProcAttr(),
    }

    cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", env_name, env_value))
    err := cmd.Start()
    if err != nil{
        return nil
    }

    os.Exit(0)
    return cmd
}
