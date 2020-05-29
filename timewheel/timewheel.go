package timewheel

import (
	"container/list"
	"errors"
	"log"
	"time"
)

// Job implementing Interface to a timer task handle function
type JobHandle interface {
	HandleMessage(interface{})
}

// TimeWheel
type TimeWheel struct {
	interval       time.Duration       // how long the point go next slot
	ticker         *time.Ticker        //ticker every interval
	slots          []*list.List        // time wheel slots
	timer          map[interface{}]int // key: timer unique key value: slot where timer in
	currentPos     int                 // current point to which slot
	slotN          int                 // slot num
	addTaskChan    chan *Task          // add task channel
	removeTaskChan chan interface{}    // delete task channel
	stopFlag       chan struct{}       // stop timer channel
}

// Task timer task
type Task struct {
	timing    time.Duration // time
	circle    int           // TimeWheel need to go circle
	key       interface{}   // timer unique key
	params    interface{}   // params
	jobHandle JobHandle
}

// New new a time wheel
func New(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:       interval,
		slots:          make([]*list.List, slotNum),
		timer:          make(map[interface{}]int),
		currentPos:     0,
		slotN:          slotNum,
		addTaskChan:    make(chan *Task),
		removeTaskChan: make(chan interface{}),
		stopFlag:       make(chan struct{}),
	}

	for i := 0; i < tw.slotN; i++ {
		tw.slots[i] = list.New()
	}

	return tw
}

//NewTask new a timer task
func NewTask(timing time.Duration, key interface{}, params interface{}, jobHandle JobHandle) *Task {
	if timing <= 0 {
		log.Fatal("task timing must greater than 0")
		return nil
	}
	return &Task{
		timing:    timing,
		key:       key,
		params:    params,
		jobHandle: jobHandle,
	}
}

//Dispatch start time wheel
func (tw *TimeWheel) Dispatch() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.dispatch()
}

//Stop stop time wheel
func (tw *TimeWheel) Stop() {
	tw.stopFlag <- struct{}{}
}

//AddTimerTask add new timer task
func (tw *TimeWheel) AddTimerTask(job *Task) error {
	return tw.addTimerTask(job)
}

//addTimerTask add new timer task
func (tw *TimeWheel) addTimerTask(job *Task) error {
	if job.timing <= 0 || job.jobHandle == nil {
		return errors.New("job is wrong")
	}
	tw.addTaskChan <- job
	return nil
}

//RemoveTimer remove timer by unipue key
func (tw *TimeWheel) RemoveTimer(key interface{}) {
	if key == nil {
		return
	}
	tw.removeTaskChan <- key
}

//dispatch really dispatch to time wheel
func (tw *TimeWheel) dispatch() {
	log.Println("time wheel start...")
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChan:
			tw.addTask(task)
		case key := <-tw.removeTaskChan:
			tw.removeTask(key)
		case <-tw.stopFlag:
			tw.ticker.Stop()
			log.Println("time wheel stop...")
			return
		}
	}
}

//tickHandler ticker timer for time wheel
func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAndRunTask(l)
	if tw.currentPos == tw.slotN-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

//scanAndRunTask scan list where timer is over and do handle function
func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}

		go task.jobHandle.HandleMessage(task.params)
		next := e.Next()
		l.Remove(e)
		if task.key != nil {
			delete(tw.timer, task.key)
		}
		e = next
	}
}

//addTask add new timer task to list
func (tw *TimeWheel) addTask(task *Task) {
	pos, circle := tw.getPositionAndCircle(task.timing)
	task.circle = circle

	tw.slots[pos].PushBack(task)

	if task.key != nil {
		tw.timer[task.key] = pos
	}
}

//getPositionAndCircle  get timer slot position and circle it should go
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle = int(delaySeconds / intervalSeconds / tw.slotN)
	pos = int(tw.currentPos+delaySeconds/intervalSeconds) % tw.slotN

	return
}

//removeTask delete task by timer unique key
func (tw *TimeWheel) removeTask(key interface{}) {
	// get the timer slot
	position, ok := tw.timer[key]
	if !ok {
		return
	}
	// get slot list
	l := tw.slots[position]
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.key == key {
			delete(tw.timer, task.key)
			l.Remove(e)
		}

		e = e.Next()
	}
}
