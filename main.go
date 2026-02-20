package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// technically this whole thing should be defined not in main
type WorkoutCategory string

const (
	CategoryDistanceBased  WorkoutCategory = "distance_based"
	CategoryTimeBased WorkoutCategory = "time_based"
)

// need to add more here in the future
type WorkoutType string

const (
	Cycling WorkoutType = "cycling"
	Walking WorkoutType = "walking"
	Running  WorkoutType = "running"
	Yoga WorkoutType = "yoga"
	Strength WorkoutType = "strength"
)

type Workout struct {
	CustomerID string `json:"customerId"`
	Type       WorkoutType `json:"type"` 
	Date       string `json:"date"`
	Time       string `json:"time"`
	Duration   int    `json:"duration"`
	Distance   int    `json:"distance"`
}


func (w WorkoutType) Category() WorkoutCategory {
	switch w {
	case Cycling, Running, Walking:
		return CategoryDistanceBased
	case Yoga, Strength:
		return CategoryTimeBased
	default:
		return ""
	}
}
func IsValidWorkoutType(t WorkoutType) bool {
	switch t {
	case Cycling, Running, Yoga, Walking, Strength:
		return true
	default:
		return false
	}
}


func main() {
	reader := bufio.NewReader(os.Stdin)
	repo := &WorkoutRepository{}

	for {
		fmt.Println("\n=== Workout Tracker ===")
		fmt.Println("1. Record Workout")
		fmt.Println("2. List Workouts")
		fmt.Println("3. Exit")
		fmt.Print("Choose option: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "1" {
			recordWorkout(reader, repo)
		} else if choice == "2" {
			listWorkouts(reader, repo)
		} else if choice == "3" {
			fmt.Println("Goodbye!")
			break
		} else {
			fmt.Println("Invalid option")
		}
	}
}

func recordWorkout(reader *bufio.Reader, repo *WorkoutRepository) {
	var workout Workout

	fmt.Print("Enter Customer ID: ")
	customerId, _ := reader.ReadString('\n')
	workout.CustomerID = strings.TrimSpace(customerId)

	fmt.Print("Enter workout type (walking/cycling/running): ")
	typeStr, _ := reader.ReadString('\n')
	typeStr = strings.TrimSpace(typeStr)

	workoutType := WorkoutType(typeStr)
	if !IsValidWorkoutType(workoutType) {
		fmt.Println("Invalid workout type!")
		return
	}
	workout.Type = workoutType

	fmt.Print("Enter date (YYYY-MM-DD): ")
	date, _ := reader.ReadString('\n')
	workout.Date = strings.TrimSpace(date)

	_, err := time.Parse("2006-01-02", workout.Date)
	if err != nil {
		fmt.Println("Invalid date format! Use YYYY-MM-DD")
		return
	}

	fmt.Print("Enter time (HH:MM): ")
	timeStr, _ := reader.ReadString('\n')
	workout.Time = strings.TrimSpace(timeStr)

	_, err = time.Parse("15:04", workout.Time)
	if err != nil {
		fmt.Println("Invalid time format! Use HH:MM")
		return
	}

	fmt.Print("Enter duration (minutes): ")
	durationStr, _ := reader.ReadString('\n')
	duration, err := strconv.Atoi(strings.TrimSpace(durationStr))
	if err != nil {
		fmt.Println("Invalid duration!")
		return
	}
	workout.Duration = duration

	// Only ask distance if distance-based workout
	if workout.Type.Category() == CategoryDistanceBased {
		fmt.Print("Enter distance (metres): ")
		distanceStr, _ := reader.ReadString('\n')
		distance, err := strconv.Atoi(strings.TrimSpace(distanceStr))
		if err != nil || distance <= 0 {
			fmt.Println("Invalid distance!")
			return
		}
		workout.Distance = distance
	}
	

	repo.Save(workout)
	fmt.Println("Workout recorded successfully!")
}

func listWorkouts(reader *bufio.Reader, repo *WorkoutRepository) {
	fmt.Print("Enter Customer ID: ")
	customerId, _ := reader.ReadString('\n')
	customerId = strings.TrimSpace(customerId)

	workouts := repo.Fetch(customerId)

	if len(workouts) == 0 {
		fmt.Println("No workouts found for this customer!")
		return
	}

	fmt.Println("\n=== Your Workouts ===")
	for _, w := range workouts {
		fmt.Printf("\nType: %s\n", w.Type)
		fmt.Printf("Date: %s\n", w.Date)
		fmt.Printf("Time: %s\n", w.Time)
		fmt.Printf("Duration: %d minutes\n", w.Duration)

		factor := w.Type.Factor()
		//var score int doesn't work and i got cannot use float64(factor) * speed (value of type float64) as int value in assignment
		var score float64

		if w.Type.Category() == CategoryDistanceBased {
			fmt.Printf("Distance: %d metres\n", w.Distance)

			if w.Duration > 0 {
				speed := float64(w.Distance) / float64(w.Duration)
				fmt.Printf("Average Speed: %.2f metres/minute\n", speed)
				score = float64(factor) * speed
			}
			
		} else if w.Type.Category() == CategoryTimeBased {
			score = float64(factor * w.Duration)
		} 

		fmt.Printf("Score: %.2f\n", score)
	}
}

func (w WorkoutType) Factor() int {
	switch w {
	case Yoga:
		return 2
	case Strength:
		return 3
	case Walking:
		return 2
	case Cycling:
		return 4
	case Running:
		return 6
	default:
		return 1
	}
}



func (w WorkoutType) String() string {
    return strings.ToUpper(string(w[:1])) + string(w[1:])
}