package main

import (
	"sort"
)

type sorted struct {
	category    string
	probability float64
}

type Classifier struct {
	datatrain           []Condition
	datatest            []Condition
	waktu               (map[string]map[string]int)
	totalWords          int
	categoriesDocuments map[string]int
	categoriesWords     map[string]int
	threshold           float64
}

func NewClassifier(categories []string, threshold float64) Classifier {
	c := Classifier{
		waktu:               make(map[string]map[string]int),
		totalWords:          0,
		categoriesWords:     make(map[string]int),
		categoriesDocuments: make(map[string]int),
		threshold:           threshold,
	}
	for _, category := range categories {
		c.categoriesWords[category] = 0
		c.categoriesDocuments[category] = 0
		c.waktu[category] = make(map[string]int)
	}
	return c
}
func (c *Classifier) addDataTrain(cond Condition) {
	c.datatrain = append(c.datatest, cond)
}
func (c *Classifier) addDataTest(cond Condition) {
	c.datatest = append(c.datatest, cond)
}

// Train the classifier
func (c *Classifier) Train(category string, time string) {
	c.categoriesWords[category]++
	c.totalWords++
	c.categoriesDocuments[category]++
}

// Classify a document
func (c *Classifier) Classify(waktu string, dmin float64, dmax float64, tmin float64, tmax float64) (category string) {
	// get all the probabilities of each category
	prob := c.Probabilities(waktu, dmin, dmax, tmin, tmax)

	// sort the categories according to probabilities
	var sp []sorted //category, prob
	for c, p := range prob {
		sp = append(sp, sorted{c, p})

		// fmt.Println(sp[1])
	}
	sort.Slice(sp, func(i, j int) bool {
		return sp[i].probability > sp[j].probability
	})

	// if the highest probability is above threshold select that
	if sp[0].probability/sp[1].probability > c.threshold {
		category = sp[0].category
	} else {
		category = "unknown"
	}
	return
}

// Probabilities of each category
func (c *Classifier) Probabilities(waktu string, dmin float64, dmax float64, tmin float64, tmax float64) (p map[string]float64) {
	p = make(map[string]float64)
	for category := range c.waktu {
		p[category] = c.pCategoryDocument(category, waktu, dmin, dmax, tmin, tmax)
	}
	return
}

// p (category | condition1|cond2|cond3|cond4)
func (c *Classifier) pCategoryDocument(category string, waktu string, dmin float64, dmax float64, tmin float64, tmax float64) float64 {
	cuaca := CUACA_CERAH
	switch category {
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
	return c.pDocumentCategory(category, waktu) * c.pNumericalCategory(cuaca, "dmin", dmin) *
		c.pNumericalCategory(cuaca, "dmax", dmax) * c.pNumericalCategory(cuaca, "tmin", tmin) *
		c.pNumericalCategory(cuaca, "tmax", tmax) * c.pCategory(category)
}

// p (condition | category)
func (c *Classifier) pDocumentCategory(category string, condition string) float64 {
	return float64(c.waktu[category][condition]+1) / float64(c.categoriesWords[category])
}

// p (condition numerical | category)
func (c *Classifier) pNumericalCategory(category float64, condition string, cond float64) float64 {
	d := NormalDist{c.avg(condition, category), c.stdev(condition, category)}
	for i := 0; i < len(condition); i++ {
		if c.datatrain[i]["cuaca"] != category {
			continue
		}
		// cond, _ := c.datatrain[i][condition]
		// fmt.Println(cond)
	}
	// cond, _ := strconv.ParseFloat(, 64)
	return (d.PDF(cond))
}

// p (category)
func (c *Classifier) pCategory(category string) float64 {
	return float64(c.categoriesDocuments[category]) / float64(len(train))
}

// fmt.Println("Rata2 : ", nb.avg("dmin", CUACA_HUJAN))
// 	fmt.Println("Standar Deviasi : ", nb.stdev("dmin", CUACA_HUJAN))
// 	d := NormalDist{nb.avg("dmin", CUACA_HUJAN), nb.stdev("dmin", CUACA_HUJAN)}
