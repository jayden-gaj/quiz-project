package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	csvName := flag.String("csv", "problems.csv", "name of csv file in question,answer format")
	timeLimit := flag.Int("time", 30, "Time limit for the test in seconds")
	shuffle := flag.Bool("Shuffle", true, "Determines if the questions should be shuffled or not")

	flag.Parse()

	problems, err := os.Open(*csvName)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(problems)
	Qs, err := reader.ReadAll()
	if err == io.EOF {
	}
	if err != nil {
		fmt.Print(err)
	}

	totalQuestions := len(Qs)

	if *shuffle {
		rand.Shuffle(totalQuestions, func(i, j int) {
			Qs[i][0], Qs[j][0] = Qs[j][0], Qs[i][0]
			Qs[i][1], Qs[j][1] = Qs[j][1], Qs[i][1]
		})
	}

	var start string
	correct := 0
	doneCh := make(chan bool)
	fmt.Printf("Type to start (Total Time = %d)\n", *timeLimit)
	fmt.Scanln(&start)

	go func() {
		for i := 0; i < totalQuestions; i++ {
			var answer string
			fmt.Print("What is " + Qs[i][0] + "\n")
			fmt.Scanln(&answer)
			if answer == Qs[i][1] {
				correct++
			}
			if totalQuestions == i+1 {
				doneCh <- true
			}
		}
	}()

	go func() {
		timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
		<-timer.C
		println(" ")
		println("Out of time")
		doneCh <- true
	}()

	<-doneCh

	fmt.Println("Correct Answers: ", correct, "\n"+"Total Questions: ", totalQuestions)
}
