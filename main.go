/*
* xim scheduler
* Copyright (C)  Alex Gaetano Padula
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
	cronparser "xim/cron"
)

type Job struct {
	Entry      string `yaml:"entry"`
	stopSignal chan struct{}
}

type Cron struct {
	jobs     []*Job
	mu       sync.Mutex
	stopChan chan os.Signal
}

func NewCron() *Cron {
	return &Cron{
		stopChan: make(chan os.Signal, 1),
	}
}

func (c *Cron) AddJob(job *Job) {
	c.mu.Lock()
	defer c.mu.Unlock()

	job.stopSignal = make(chan struct{})
	c.jobs = append(c.jobs, job)
}

func (c *Cron) Start() {
	for _, job := range c.jobs {
		go func(job *Job) {
			for {
				parsedEntry, err := cronparser.ParseCronJob(job.Entry)
				if err != nil {
					log.Println(err.Error())
				}

				nextScheduleTime, err := cronparser.GetNextScheduledTime(parsedEntry)
				if err != nil {
					log.Println(err.Error())
				}
				duration := nextScheduleTime.Sub(time.Now())
				select {
				case <-time.After(duration):
					err := executeCommand(parsedEntry.Command)
					if err != nil {
						log.Println(err.Error())
					}

				case <-job.stopSignal:
					return
				}
				time.Sleep(time.Millisecond)
			}
		}(job)
	}
}

func (c *Cron) Stop() {
	c.stopChan <- syscall.SIGTERM
}

func (c *Cron) WaitAndSaveJobsOnShutdown() {
	signal.Notify(c.stopChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-c.stopChan

	// Save jobs to YAML file on shutdown
	if err := saveJobs(c.jobs); err != nil {
		fmt.Printf("Error saving jobs: %v\n", err)
	}

	fmt.Println("xim stopped.")
	os.Exit(0)
}

func main() {
	cron := NewCron()

	// Load jobs
	jobs, err := loadJobs("ximtab")
	if err != nil {
		fmt.Printf("Error loading jobs: %v\n", err)
	}

	// Add loaded jobs to the cron system
	for _, job := range jobs {
		cron.AddJob(job)
	}

	// Start the cron system
	go cron.Start()

	// Wait for shutdown signal and save jobs on exit
	cron.WaitAndSaveJobsOnShutdown()
}

func executeCommand(command string) error {
	fmt.Printf("Executing command: %s\n", command)

	//output, err := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...).Output()
	//if err != nil {
	//	return err
	//}
	//fmt.Println(string(output))

	go func() {
		log.Println("running")
		cmd := exec.Command(strings.Split(command, " ")[0], strings.Split(command, " ")[1:]...)
		cmd.Start()
		cmd.Wait()
		log.Println("done")
	}()
	return nil
}

func loadJobs(filename string) ([]*Job, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var jobs []*Job

	for _, j := range strings.Split(string(data), "\n") {
		jobs = append(jobs, &Job{
			Entry:      strings.TrimSpace(j),
			stopSignal: make(chan struct{}),
		})
	}

	return jobs, nil
}

func saveJobs(jobs []*Job) error {
	var data string

	for i, j := range jobs {
		if i == len(jobs)-1 {
			data += fmt.Sprintf("%s", j.Entry)
			continue
		}
		data += fmt.Sprintf("%s\n", j.Entry)
	}

	err := os.WriteFile("ximtab", []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}
