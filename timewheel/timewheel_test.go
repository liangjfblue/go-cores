package timewheel

import (
	"log"
	"testing"
	"time"
)

type Job1 struct {
}

func (j *Job1) HandleMessage(params interface{}) {
	log.Println("Job1 time is over")
}

type Job2 struct {
}

func (j *Job2) HandleMessage(params interface{}) {
	if paramss, ok := params.(map[string]interface{}); ok {
		log.Println("Job2 time is over")
		for k, v := range paramss {
			log.Println(k, " ", v)
		}
	}
}

type Job3 struct {
}

func (j *Job3) HandleMessage(params interface{}) {
	log.Println("Job3 time is over")
}

func TestTimeWheel(t *testing.T) {
	ws := New(1*time.Second, 10)

	ws.Dispatch()
	defer ws.Stop()
	job1 := Job1{}
	task1 := NewTask(2*time.Second, "test1", nil, &job1)
	if err := ws.AddTimerTask(task1); err != nil {
		log.Fatal(err.Error())
		return
	}

	job2 := Job2{}
	task2 := NewTask(5*time.Second, "test2", map[string]interface{}{"param1": 1, "param2": "hello"}, &job2)
	if err := ws.AddTimerTask(task2); err != nil {
		log.Fatal(err.Error())
		return
	}

	job3 := Job3{}
	task3 := NewTask(12*time.Second, "test3", nil, &job3)
	if err := ws.AddTimerTask(task3); err != nil {
		log.Fatal(err.Error())
		return
	}

	time.Sleep(time.Second * 20)

}
