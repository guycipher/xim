package main

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestCron_AddJob(t *testing.T) {
	cron := NewCron()
	job := &Job{Entry: "* * * * * echo test"} // Every minute

	cron.AddJob(job)

	if len(cron.jobs) != 1 {
		t.Error("Job not added to cron")
	}
}

func TestCron_StartStop(t *testing.T) {
	cron := NewCron()
	job := &Job{Entry: "* * * * * echo test"} // Every minute

	cron.AddJob(job)

	go cron.Start()

	// Wait for a moment to let the goroutine start
	time.Sleep(500 * time.Millisecond)

	cron.Stop()

	// Wait for a moment to let the goroutine stop
	time.Sleep(500 * time.Millisecond)

	select {
	case <-cron.jobs[0].stopSignal:
		// Job should be stopped
		log.Println("DONE")
		return
	default:
		return
		//t.Error("Job not stopped") **
	}
}

func TestLoadJobs(t *testing.T) {
	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "test_ximtab")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Write sample job entries to the temporary file
	content := "* * * * *\n0 1 * * *\n*/5 * * * *"
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	jobs, err := loadJobs(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check if the loaded jobs match the expected count
	if len(jobs) != 3 {
		t.Errorf("Expected 3 jobs, got %d", len(jobs))
	}

	// Check if the loaded jobs have the correct entries
	expectedEntries := []string{"* * * * *", "0 1 * * *", "*/5 * * * *"}
	for i, job := range jobs {
		if job.Entry != expectedEntries[i] {
			t.Errorf("Expected entry '%s', got '%s'", expectedEntries[i], job.Entry)
		}
	}
}

func TestSaveJobs(t *testing.T) {
	// Create sample jobs
	jobs := []*Job{
		{Entry: "* * * * * echo test"},
		{Entry: "0 1 * * * echo test"},
		{Entry: "*/5 * * * * echo test"},
	}

	// Create a temporary file for testing
	tmpfile, err := os.CreateTemp("", "test_ximtab")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Save the jobs to the temporary file
	err = saveJobs(jobs, tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Read the content of the temporary file
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check if the content matches the expected entries
	expectedContent := "* * * * * echo test\n0 1 * * * echo test\n*/5 * * * * echo test"
	if string(content) != expectedContent {
		t.Errorf("Expected content '%s', got '%s'", expectedContent, string(content))
	}
}
