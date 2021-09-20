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
	CUACA_HUJAN   = 1.0
	CUACA_CERAH   = 2.0
	CUACA_MENDUNG = 3.0
)

var (
	testPercentage = 0.1 //presentasi data test
	datafile       = "data-cuaca.csv"
	threshold      = 1.1

	//exampleif `threshold` is `1.5` this means the category with the highest probability
	// needs to be 1.5 times higher than the second highest probability.
	// If the top category fails the threshold we will classify it as `unknown`.
)

// datasets
type document struct {
	time  string
	class string
}
type Condition map[string]float64

//dipisahkan untuk training dan test
var train []document
var test []document

var categories = []string{"Hujan", "Berawan", "Cerah"}

func main() {

	nb := NewClassifier(categories, threshold)
	nb.setupData(datafile)
	fmt.Println("Data file used:", datafile)
	fmt.Println("no of docs in TRAIN dataset:", len(train))
	fmt.Println("no of docs in TEST dataset:", len(test))
	// train on train dataset
	for _, doc := range train {
		nb.Train(doc.class, doc.time)
	}
	// NormalDist{
	// 	nb.avg("dmin", CUACA_HUJAN),
	// 	nb.stdev("dmin", CUACA_HUJAN),
	// }
	a := NormalDist{nb.avg("dmin", CUACA_CERAH), nb.stdev("dmin", CUACA_CERAH)}
	fmt.Println("hasil training", nb.categoriesDocuments, a)
	// validate on test dataset
	count, accurates, unknowns := 0, 0, 0
	for i, doc := range test {
		count++
		sentiment := nb.Classify(doc.time, nb.datatrain[i]["dmin"], nb.datatrain[i]["dmax"],
			nb.datatrain[i]["tmin"], nb.datatrain[i]["tmax"])
		// fmt.Println(nb.datatrain[i]["dmin"], nb.datatrain[i]["dmax"],
		// 	nb.datatrain[i]["tmin"], nb.datatrain[i]["tmax"])
		if sentiment == doc.class {
			accurates++
		}
		if sentiment == "unknown" {
			unknowns++
		}
	}
	fmt.Printf("Accuracy on TEST dataset is %2.1f%% with %2.1f%% unknowns",
		float64(accurates)*100/float64(count), float64(unknowns)*100/float64(count))
	// // validate on the first 100 docs in the train dataset
	// for i, doc := range train[0:100] {
	// 	count++
	// 	sentiment := nb.Classify(doc.time, nb.datatrain[i]["dmin"], nb.datatrain[i]["dmax"],
	// 		nb.datatrain[i]["tmin"], nb.datatrain[i]["tmax"])
	// 	if sentiment == doc.class {
	// 		accurates++
	// 	}
	// 	if sentiment == "unknown" {
	// 		unknowns++
	// 	}
	// }
	// fmt.Printf("\nAccuracy on TRAIN dataset is %2.1f%% with %2.1f%% unknowns",
	// 	float64(accurates)*100/float64(count), float64(unknowns)*100/float64(count))
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
		dmin, err := strconv.ParseFloat(csvData[i][2], 64)
		if err != nil {
			continue
		}
		dmax, err := strconv.ParseFloat(csvData[i][3], 64)
		if err != nil {
			continue
		}
		tmin, err := strconv.ParseFloat(csvData[i][4], 64)
		if err != nil {
			continue
		}
		tmax, err := strconv.ParseFloat(csvData[i][5], 64)
		if err != nil {
			continue
		}
		cuaca := CUACA_CERAH
		switch csvData[i][1] {
		case "Cerah":
			cuaca = CUACA_CERAH
			break
		case "Berawan":
			cuaca = CUACA_MENDUNG
			break
		case "Hujan":
			cuaca = CUACA_HUJAN
			break
		}
		class := csvData[i][1]
		waktu := csvData[i][0]
		// fmt.Println(waktu, dmin, dmin, tmin, tmax, class)
		//dibagi data train dan test
		if rand.Float64() > testPercentage {
			train = append(train, document{waktu, class})
			c.addDataTrain(Condition{
				"dmin":  dmin,
				"dmax":  dmax,
				"tmin":  tmin,
				"tmax":  tmax,
				"cuaca": cuaca,
			})

		} else {
			test = append(test, document{waktu, class})
			c.addDataTest(Condition{
				"dmin":  dmin,
				"dmax":  dmax,
				"tmin":  tmin,
				"tmax":  tmax,
				"cuaca": cuaca,
			})
		}
	}

}
