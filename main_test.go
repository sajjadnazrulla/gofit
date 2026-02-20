package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestRecordWorkout_Success(t *testing.T) {
	testFile := "test_record.json"
	defer os.Remove(testFile)

	replaceFilename(testFile)
	defer replaceFilename("workouts.json")

	input := "C001\nrunning\n2024-01-15\n08:30\n30\n5000\n"
	reader := bufio.NewReader(strings.NewReader(input))
	repo := &WorkoutRepository{}

	recordWorkout(reader, repo)

	data, _ := os.ReadFile(testFile)
	var workouts []Workout
	json.Unmarshal(data, &workouts)

	if len(workouts) != 1 {
		t.Fatalf("Expected 1 workout, got %d", len(workouts))
	}

	w := workouts[0]
	if w.CustomerID != "C001" {
		t.Errorf("Expected CustomerID C001, got %s", w.CustomerID)
	}
	if w.Type != Running {
		t.Errorf("Expected Type running, got %s", w.Type)
	}
	if w.Date != "2024-01-15" {
		t.Errorf("Expected Date 2024-01-15, got %s", w.Date)
	}
	if w.Time != "08:30" {
		t.Errorf("Expected Time 08:30, got %s", w.Time)
	}
	if w.Duration != 30 {
		t.Errorf("Expected Duration 30, got %d", w.Duration)
	}
	if w.Distance != 5000 {
		t.Errorf("Expected Distance 5000, got %d", w.Distance)
	}
}

func TestRecordWorkout_InvalidWorkoutType(t *testing.T) {
	testFile := "test_invalid_type.json"
	defer os.Remove(testFile)

	replaceFilename(testFile)
	defer replaceFilename("workouts.json")

	input := "C001\nswimming\n2024-01-15\n08:30\n30\n5000\n"
	reader := bufio.NewReader(strings.NewReader(input))
	repo := &WorkoutRepository{}

	recordWorkout(reader, repo)

	data, _ := os.ReadFile(testFile)
	var workouts []Workout
	if data != nil {
		json.Unmarshal(data, &workouts)
	}

	if len(workouts) != 0 {
		t.Errorf("Expected 0 workouts for invalid type, got %d", len(workouts))
	}
}

func TestRecordWorkout_InvalidDate(t *testing.T) {
	testFile := "test_invalid_date.json"
	defer os.Remove(testFile)

	replaceFilename(testFile)
	defer replaceFilename("workouts.json")

	input := "C001\nrunning\n2024-13-45\n08:30\n30\n5000\n"
	reader := bufio.NewReader(strings.NewReader(input))
	repo := &WorkoutRepository{}

	recordWorkout(reader, repo)

	data, _ := os.ReadFile(testFile)
	var workouts []Workout
	if data != nil {
		json.Unmarshal(data, &workouts)
	}

	if len(workouts) != 0 {
		t.Errorf("Expected 0 workouts for invalid date, got %d", len(workouts))
	}
}

func TestRecordWorkout_InvalidTime(t *testing.T) {
	testFile := "test_invalid_time.json"
	defer os.Remove(testFile)

	replaceFilename(testFile)
	defer replaceFilename("workouts.json")

	input := "C001\nrunning\n2024-01-15\n25:99\n30\n5000\n"
	reader := bufio.NewReader(strings.NewReader(input))
	repo := &WorkoutRepository{}

	recordWorkout(reader, repo)

	data, _ := os.ReadFile(testFile)
	var workouts []Workout
	if data != nil {
		json.Unmarshal(data, &workouts)
	}

	if len(workouts) != 0 {
		t.Errorf("Expected 0 workouts for invalid time, got %d", len(workouts))
	}
}

func TestRecordWorkout_MultipleWorkouts(t *testing.T) {
	testFile := "test_multiple.json"
	defer os.Remove(testFile)

	replaceFilename(testFile)
	defer replaceFilename("workouts.json")

	repo := &WorkoutRepository{}

	input1 := "C001\nwalking\n2024-01-15\n08:00\n20\n2000\n"
	reader1 := bufio.NewReader(strings.NewReader(input1))
	recordWorkout(reader1, repo)

	input2 := "C002\ncycling\n2024-01-16\n09:00\n45\n15000\n"
	reader2 := bufio.NewReader(strings.NewReader(input2))
	recordWorkout(reader2, repo)

	data, _ := os.ReadFile(testFile)
	var workouts []Workout
	json.Unmarshal(data, &workouts)

	if len(workouts) != 2 {
		t.Fatalf("Expected 2 workouts, got %d", len(workouts))
	}
}

func TestListWorkouts_Success(t *testing.T) {
	testFile := "test_list.json"
	defer os.Remove(testFile)

	replaceFilename(testFile)
	defer replaceFilename("workouts.json")

	workouts := []Workout{
		{CustomerID: "C001", Type: Running, Date: "2024-01-15", Time: "08:00", Duration: 30, Distance: 5000},
		{CustomerID: "C002", Type: Walking, Date: "2024-01-16", Time: "09:00", Duration: 20, Distance: 2000},
		{CustomerID: "C001", Type: Cycling, Date: "2024-01-17", Time: "10:00", Duration: 45, Distance: 15000},
	}

	data, _ := json.MarshalIndent(workouts, "", "  ")
	os.WriteFile(testFile, data, 0644)

	input := "C001\n"
	reader := bufio.NewReader(strings.NewReader(input))
	repo := &WorkoutRepository{}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listWorkouts(reader, repo)

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "Running") {
		t.Error("Expected output to contain 'Running'")
	}
	if !strings.Contains(output, "Cycling") {
		t.Error("Expected output to contain 'Cycling'")
	}
	if !strings.Contains(output, "2024-01-15") {
		t.Error("Expected output to contain date 2024-01-15")
	}
}

func TestListWorkouts_NoWorkoutsFound(t *testing.T) {
	testFile := "test_list_empty.json"
	defer os.Remove(testFile)

	replaceFilename(testFile)
	defer replaceFilename("workouts.json")

	workouts := []Workout{
		{CustomerID: "C001", Type: Running, Date: "2024-01-15", Time: "08:00", Duration: 30, Distance: 5000},
	}

	data, _ := json.MarshalIndent(workouts, "", "  ")
	os.WriteFile(testFile, data, 0644)

	input := "C999\n"
	reader := bufio.NewReader(strings.NewReader(input))
	repo := &WorkoutRepository{}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listWorkouts(reader, repo)

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "No workouts found") {
		t.Error("Expected 'No workouts found' message")
	}
}

func TestListWorkouts_FileNotExists(t *testing.T) {
	testFile := "nonexistent_list.json"

	replaceFilename(testFile)
	defer replaceFilename("workouts.json")

	input := "C001\n"
	reader := bufio.NewReader(strings.NewReader(input))
	repo := &WorkoutRepository{}

	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	listWorkouts(reader, repo)

	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "No workouts found") {
		t.Error("Expected 'No workouts found' message when file doesn't exist")
	}
}