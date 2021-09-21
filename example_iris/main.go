package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	I_SETOSA     = 1.0
	I_VERSICOLOR = 2.0
	I_VIRGINICA  = 3.0
)

var (
	testPercentage = 0.1 //presentasi data test
	datafile       = "iris.csv"
	threshold      = 1.1

	//exampleif `threshold` is `1.5` this means the category with the highest probability
	// needs to be 1.5 times higher than the second highest probability.
	// If the top category fails the threshold we will classify it as `unknown`.
)

// datasets
type document struct {
	class string
}
type Condition map[string]float64

//dipisahkan untuk training dan test
var train []document
var test []document

var categories = []string{"Iris-setosa", "Iris-versicolor", "Iris-virginica"}

func main() {

	nb := NewClassifier(categories, threshold)
	nb.setupData(datafile)
	fmt.Println("Data file used:", datafile)
	fmt.Println("no of docs in TRAIN dataset:", len(train))
	fmt.Println("no of docs in TEST dataset:", len(test))
	// train on train dataset
	for _, doc := range train {
		nb.Train(doc.class)
	}
	// validate on test dataset
	count, accurates, unknowns := 0, 0, 0
	for i, doc := range test {
		count++
		sentiment := nb.Classify(nb.datatest[i]["l1"], nb.datatest[i]["l2"],
			nb.datatest[i]["l3"], nb.datatest[i]["l4"])
		if sentiment == doc.class {
			accurates++
		}
		if sentiment == "unknown" {
			unknowns++
		}
	}
	fmt.Printf("Accuracy on TEST dataset is %2.1f%% with %2.1f%% unknowns",
		float64(accurates)*100/float64(count), float64(unknowns)*100/float64(count))
	// validate on the first 100 docs in the train dataset
	count, accurates, unknowns = 0, 0, 0
	for i, doc := range train[0:100] {
		count++
		sentiment := nb.Classify(nb.datatrain[i]["l1"], nb.datatrain[i]["l2"],
			nb.datatrain[i]["l3"], nb.datatrain[i]["l4"])
		if sentiment == doc.class {
			accurates++
		}
		if sentiment == "unknown" {
			unknowns++
		}
	}
	fmt.Printf("\nAccuracy on TRAIN dataset is %2.1f%% with %2.1f%% unknowns",
		float64(accurates)*100/float64(count), float64(unknowns)*100/float64(count))
}

func (c *Classifier) setupData(file string) {
	rand.Seed(time.Now().UTC().UnixNano())
	f, err := os.Open(file)
	if err != nil {
		return
	}
	csvReader := csv.NewReader(f)
	csvData, err := csvReader.ReadAll()
	for i := 0; i < len(csvData); i++ {
		l1, err := strconv.ParseFloat(csvData[i][0], 64)
		if err != nil {
			continue
		}
		l2, err := strconv.ParseFloat(csvData[i][1], 64)
		if err != nil {
			continue
		}
		l3, err := strconv.ParseFloat(csvData[i][2], 64)
		if err != nil {
			continue
		}
		l4, err := strconv.ParseFloat(csvData[i][3], 64)
		if err != nil {
			continue
		}
		class := I_SETOSA
		switch csvData[i][4] {
		case "Iris-setosa":
			class = I_SETOSA
			break
		case "Iris-versicolor":
			class = I_VERSICOLOR
			break
		case "Iris-virginica":
			class = I_VIRGINICA
			break
		}
		ctg := csvData[i][4]
		// waktu := csvData[i][0]
		// fmt.Println(waktu, dmin, dmin, tmin, tmax, class)
		//dibagi data train dan test
		if rand.Float64() > testPercentage {
			train = append(train, document{ctg})
			c.addDataTrain(Condition{
				"l1":    l1,
				"l2":    l2,
				"l3":    l3,
				"l4":    l4,
				"class": class,
			})

		} else {
			test = append(test, document{ctg})
			c.addDataTest(Condition{
				"l1":    l1,
				"l2":    l2,
				"l3":    l3,
				"l4":    l4,
				"class": class,
			})
		}
	}

}
