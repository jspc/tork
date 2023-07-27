package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/tork/task"
	"github.com/tork/worker"
)

func main() {
	w, err := worker.NewWorker()
	if err != nil {
		panic(err)
	}
	t := task.Task{
		ID:    uuid.New().String(),
		State: task.Pending,
		Name:  "test-container-1",
		Image: "postgres:13",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	w.EnqueueTask(t)

	err = w.RunTask()
	if err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)

	err = w.StopTask(t)
	if err != nil {
		panic(err)
	}

	// 	te := task.TaskEvent{
	// 		ID:        uuid.New(),
	// 		State:     task.Pending,
	// 		Timestamp: time.Now(),
	// 		Task:      t,
	// 	}

	// 	fmt.Printf("task: %v\n", t)
	// 	fmt.Printf("task event: %v\n", te)

	// 	w := worker.Worker{
	// 		Queue: *queue.New(),
	// 		DB:    make(map[uuid.UUID]task.Task),
	// 	}
	// 	fmt.Printf("worker: %v\n", w)
	// 	w.CollectStats()
	// 	w.RunTask()
	// 	w.StartTask(t)
	// 	w.StopTask(t)

	// 	m := manager.Manager{
	// 		Pending: *queue.New(),
	// 		TaskDB:  make(map[string][]task.Task),
	// 		EventDB: make(map[string][]task.TaskEvent),
	// 		Workers: []string{w.Name},
	// 	}

	// 	fmt.Printf("manager: %v\n", m)
	// 	m.SelectWorker()
	// 	m.UpdateTasks()
	// 	m.SendWork()

	// 	n := node.Node{
	// 		Name:   "Node-1",
	// 		IP:     "192.168.1.1",
	// 		Cores:  4,
	// 		Memory: 1024,
	// 		Disk:   25,
	// 		Role:   "worker",
	// 	}

	// 	fmt.Printf("node: %v\n", n)

	// 	fmt.Printf("create a test container\n")
	// 	dockerTask, createResult := createContainer()
	// 	if createResult.Error != nil {
	// 		fmt.Println(createResult.Error)
	// 		os.Exit(1)
	// 	}

	// 	time.Sleep(time.Second * 2)
	// 	fmt.Printf("stopping container %s\n", createResult.ContainerID)
	// 	_ = stopContainer(dockerTask)
	// }

	// func createContainer() (*task.Docker, *task.DockerResult) {
	// 	c := task.Config{
	// 		Name:  "test-container-1",
	// 		Image: "postgres:13",
	// 		Env: []string{
	// 			"POSTGRES_USER=cube",
	// 			"POSTGRES_PASSWORD=secret",
	// 		},
	// 	}

	// 	dc, _ := client.NewClientWithOpts(client.FromEnv)
	// 	d := &task.Docker{
	// 		Client: dc,
	// 		Config: c,
	// 	}

	// 	result := d.Run()
	// 	if result.Error != nil {
	// 		fmt.Printf("%v\n", result.Error)
	// 		return nil, nil
	// 	}

	// 	fmt.Printf(
	// 		"Container %s is running with config %v\n", result.ContainerID, c)
	// 	return d, &result
	// }

	// func stopContainer(d *task.Docker) *task.DockerResult {
	// 	result := d.Stop(d.ContainerID)
	// 	if result.Error != nil {
	// 		fmt.Printf("%v\n", result.Error)
	// 		return nil
	// 	}

	// 	fmt.Printf(
	// 		"Container %s has been stopped and removed\n", result.ContainerID)
	// 	return &result
}